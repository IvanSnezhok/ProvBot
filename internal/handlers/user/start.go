package user

import (
	"context"
	"fmt"
	"strings"

	"provbot/internal/handlers"
	"provbot/internal/i18n"
	"provbot/internal/models"
	"provbot/internal/repository"
	"provbot/internal/service"
	"provbot/internal/state"
	"provbot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartHandler struct {
	userService    *service.UserService
	billingService *service.BillingService
	outageService  *service.OutageService
	billingRepo    *repository.BillingRepository
	userRepo       *repository.UserRepository
	stateManager   *state.StateManager
	config         *utils.Config
}

func NewStartHandler(
	userService *service.UserService,
	billingService *service.BillingService,
	outageService *service.OutageService,
	billingRepo *repository.BillingRepository,
	userRepo *repository.UserRepository,
	stateManager *state.StateManager,
	config *utils.Config,
) *StartHandler {
	return &StartHandler{
		userService:    userService,
		billingService: billingService,
		outageService:  outageService,
		billingRepo:    billingRepo,
		userRepo:       userRepo,
		stateManager:   stateManager,
		config:         config,
	}
}

// HandleStart handles /start command
func (h *StartHandler) HandleStart(ctx *handlers.HandlerContext) error {
	// Show language selection
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Українська", "start_lang_ua"),
			tgbotapi.NewInlineKeyboardButtonData("English", "start_lang_en"),
			tgbotapi.NewInlineKeyboardButtonData("Русский", "start_lang_ru"),
		),
	)

	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("start_message"))
	msg.ReplyMarkup = keyboard
	_, err := ctx.Bot.Send(msg)
	return err
}

// HandleLanguageCallback handles language selection callback
func (h *StartHandler) HandleLanguageCallback(ctx *handlers.HandlerContext) error {
	callback := ctx.Update.CallbackQuery
	if callback == nil {
		return nil
	}

	lang := strings.TrimPrefix(callback.Data, "start_lang_")

	// Map callback to language code
	langMap := map[string]string{
		"ua": "ua",
		"en": "en",
		"ru": "ru",
	}

	langCode, ok := langMap[lang]
	if !ok {
		langCode = i18n.DefaultLanguage
	}

	// Update user language
	if err := h.userService.UpdateLanguage(context.Background(), int64(callback.From.ID), langCode); err != nil {
		utils.Logger.WithError(err).Error("Failed to update language")
	}

	// Get translator for selected language
	translator := i18n.GetGlobalTranslator(langCode)

	// Request phone number
	contactButton := tgbotapi.NewKeyboardButtonContact(translator.Get("send_phone"))
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(contactButton),
	)
	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = true

	langNames := map[string]string{
		"ua": "Українську",
		"en": "English",
		"ru": "Русский",
	}
	langName := langNames[langCode]
	if langName == "" {
		langName = langCode
	}

	text := translator.Getf("phone_request", langName)
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboard

	// Set state
	h.stateManager.SetState(int64(callback.From.ID), state.StateWaitingPhone, nil)

	// Answer callback
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	_, _ = ctx.Bot.Request(callbackConfig)

	_, err := ctx.Bot.Send(msg)
	return err
}

// HandleContact handles phone number contact
func (h *StartHandler) HandleContact(ctx *handlers.HandlerContext) error {
	if ctx.Update.Message.Contact == nil {
		return fmt.Errorf("no contact provided")
	}

	userState, _, exists := h.stateManager.GetState(ctx.Update.Message.From.ID)
	if !exists || userState != state.StateWaitingPhone {
		return nil // Not waiting for phone
	}

	phone := ctx.Update.Message.Contact.PhoneNumber
	formattedPhone := h.formatPhoneNumber(phone)

	// Update user phone number
	user := ctx.User
	if user != nil {
		user.PhoneNumber = &formattedPhone
		if err := h.userRepo.Update(context.Background(), user); err != nil {
			utils.Logger.WithError(err).Error("Failed to update phone number")
		}
	}

	// Search in billing
	billingUser, contract, err := h.searchInBilling(formattedPhone)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to search in billing")
	}

	// Update contract if found
	if contract != "" && user != nil {
		user.Contract = &contract
		if err := h.userRepo.Update(context.Background(), user); err != nil {
			utils.Logger.WithError(err).Error("Failed to update contract")
		}
	}

	// Notify admins
	h.notifyAdmins(ctx, user, formattedPhone, contract, billingUser != nil)

	// Clear state
	h.stateManager.ClearState(ctx.Update.Message.From.ID)

	// Send welcome message with main menu
	return h.ShowMainMenu(ctx, billingUser != nil)
}

// formatPhoneNumber formats phone number (removes + and spaces)
func (h *StartHandler) formatPhoneNumber(phone string) string {
	phone = strings.ReplaceAll(phone, "+", "")
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	return phone
}

// searchInBilling searches user in billing by phone number
func (h *StartHandler) searchInBilling(phone string) (*models.BillingUser, string, error) {
	user, err := h.billingService.SearchUser(context.Background(), phone)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", nil
	}
	return user, user.Contract, nil
}

// notifyAdmins notifies admins about new user
func (h *StartHandler) notifyAdmins(ctx *handlers.HandlerContext, user *models.User, phone, contract string, found bool) {
	fullName := ctx.Update.Message.From.FirstName
	if ctx.Update.Message.From.LastName != "" {
		fullName += " " + ctx.Update.Message.From.LastName
	}
	message := fmt.Sprintf("Новий клієнт: %s, ID: %d\nЗ номером телефону: %s\n",
		fullName, ctx.Update.Message.From.ID, phone)

	if found && contract != "" {
		message += fmt.Sprintf("Клієнт знайдений в білінгу, номер договору: %s\n", contract)
	} else {
		message += "Клієнт не знайдений в білінгу\n"
	}

	for _, adminID := range h.config.AdminTelegramIDs {
		msg := tgbotapi.NewMessage(adminID, message)
		_, _ = ctx.Bot.Send(msg)
	}
}

// ShowMainMenu shows main menu keyboard
func (h *StartHandler) ShowMainMenu(ctx *handlers.HandlerContext, hasContract bool) error {
	var rows [][]tgbotapi.KeyboardButton

	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(ctx.Translator.Get("menu_topup")),
		tgbotapi.NewKeyboardButton(ctx.Translator.Get("menu_shop")),
	))
	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(ctx.Translator.Get("menu_support")),
		tgbotapi.NewKeyboardButton(ctx.Translator.Get("menu_time_pay")),
	))

	// Add "Connect Friend" button only for users with contracts
	// Add "Connection Request" button only for users without contracts
	if hasContract {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(ctx.Translator.Get("connect_friend")),
		))
	} else {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(ctx.Translator.Get("connection_request")),
		))
	}

	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(ctx.Translator.Get("menu_back")),
	))

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = false

	var text string

	// Check for active outages and show network status if user has contract
	if hasContract && ctx.User != nil && ctx.User.Contract != nil {
		text = ctx.Translator.Get("welcome_registered")

		// Check for outages
		outageMsg, err := h.checkOutageForUser(ctx)
		if err != nil {
			utils.Logger.WithError(err).Error("Failed to check outage")
		} else if outageMsg != "" {
			text = ctx.Translator.Get("outage_warning") + "\n\n" + outageMsg + "\n\n" + text
		}

		// Show network status and balance
		statusInfo := h.getNetworkStatus(ctx)
		if statusInfo != "" {
			text += "\n\n" + statusInfo
		}
	} else {
		// User not found in billing
		text = ctx.Translator.Get("not_found_billing")
	}

	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboard
	_, err := ctx.Bot.Send(msg)
	return err
}

// getNetworkStatus returns the current network status and balance info for the user
func (h *StartHandler) getNetworkStatus(ctx *handlers.HandlerContext) string {
	if ctx.User == nil {
		return ""
	}

	billingUser, err := h.billingService.GetBillingUser(context.Background(), int64(ctx.User.ID))
	if err != nil || billingUser == nil {
		return ""
	}

	// Determine service status
	var statusText string
	if billingUser.Status == "on" {
		statusText = ctx.Translator.Get("service_on")
	} else {
		statusText = ctx.Translator.Get("service_off")
	}

	// Format: Balance: X grn | Status: On/Off
	return fmt.Sprintf("%s: %.2f грн | %s: %s",
		ctx.Translator.Get("admin_balance"),
		billingUser.Balance,
		ctx.Translator.Get("admin_status"),
		statusText)
}

// checkOutageForUser checks if there's an active outage for the user
func (h *StartHandler) checkOutageForUser(ctx *handlers.HandlerContext) (string, error) {
	if h.outageService == nil || ctx.User == nil || ctx.User.Contract == nil {
		return "", nil
	}

	// Get billing group from user contract
	groupID := ""
	contract := *ctx.User.Contract

	// Try to get billing user to get the group
	billingUser, err := h.billingService.GetBillingUser(context.Background(), int64(ctx.User.ID))
	if err == nil && billingUser != nil && billingUser.Group != 0 {
		groupID = fmt.Sprintf("%d", billingUser.Group)
	}

	// Check for active outage
	outageMsg, err := h.outageService.GetOutageMessageForUser(context.Background(), groupID, contract)
	if err != nil {
		return "", err
	}

	return outageMsg, nil
}

// HandleConnectFriend handles the "Connect Friend" promotion message
func (h *StartHandler) HandleConnectFriend(ctx *handlers.HandlerContext) error {
	if ctx.User == nil || ctx.User.Contract == nil {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("no_contract"))
		_, err := ctx.Bot.Send(msg)
		return err
	}

	contract := *ctx.User.Contract
	address := ""

	// Get user's address from billing
	billingUser, err := h.billingService.GetBillingUser(context.Background(), int64(ctx.User.ID))
	if err == nil && billingUser != nil {
		address = billingUser.Address
	}

	// Format promo message with contract and address
	text := ctx.Translator.Getf("connect_friend_promo", contract, address)

	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	_, err = ctx.Bot.Send(msg)
	return err
}

// HandleConnectionRequest starts the connection request flow
func (h *StartHandler) HandleConnectionRequest(ctx *handlers.HandlerContext) error {
	// Set state to wait for connection request details
	h.stateManager.SetState(ctx.Update.Message.From.ID, state.StateWaitingConnectionRequest, nil)

	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("connection_request_prompt"))
	_, err := ctx.Bot.Send(msg)
	return err
}

// HandleConnectionRequestInput handles the connection request input
func (h *StartHandler) HandleConnectionRequestInput(ctx *handlers.HandlerContext) error {
	text := ctx.Update.Message.Text

	// Clear state
	h.stateManager.ClearState(ctx.Update.Message.From.ID)

	// Notify admins about the connection request
	h.notifyAdminsConnectionRequest(ctx, text)

	// Send confirmation to user
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("connection_request_sent"))
	_, err := ctx.Bot.Send(msg)
	return err
}

// notifyAdminsConnectionRequest notifies admins about a new connection request
func (h *StartHandler) notifyAdminsConnectionRequest(ctx *handlers.HandlerContext, requestText string) {
	fullName := ctx.Update.Message.From.FirstName
	if ctx.Update.Message.From.LastName != "" {
		fullName += " " + ctx.Update.Message.From.LastName
	}

	username := ""
	if ctx.Update.Message.From.UserName != "" {
		username = "@" + ctx.Update.Message.From.UserName
	}

	message := fmt.Sprintf("📝 Нова заявка на підключення!\n\nВід: %s %s\nTelegram ID: %d\n\nТекст заявки:\n%s",
		fullName, username, ctx.Update.Message.From.ID, requestText)

	for _, adminID := range h.config.AdminTelegramIDs {
		msg := tgbotapi.NewMessage(adminID, message)
		_, _ = ctx.Bot.Send(msg)
	}
}
