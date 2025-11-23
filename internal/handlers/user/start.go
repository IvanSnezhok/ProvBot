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
	billingRepo    *repository.BillingRepository
	userRepo       *repository.UserRepository
	stateManager   *state.StateManager
	config         *utils.Config
}

func NewStartHandler(
	userService *service.UserService,
	billingService *service.BillingService,
	billingRepo *repository.BillingRepository,
	userRepo *repository.UserRepository,
	stateManager *state.StateManager,
	config *utils.Config,
) *StartHandler {
	return &StartHandler{
		userService:    userService,
		billingService: billingService,
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
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(ctx.Translator.Get("menu_topup")),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(ctx.Translator.Get("menu_support")),
			tgbotapi.NewKeyboardButton(ctx.Translator.Get("menu_time_pay")),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(ctx.Translator.Get("menu_back")),
		),
	)
	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = false

	text := ctx.Translator.Get("welcome_registered")
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboard
	_, err := ctx.Bot.Send(msg)
	return err
}
