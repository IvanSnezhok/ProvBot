package admin

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"provbot/internal/handlers"
	"provbot/internal/models"
	"provbot/internal/repository"
	"provbot/internal/state"
	"provbot/internal/utils"
)

type UsersHandler struct {
	billingRepo  *repository.BillingRepository
	userRepo     *repository.UserRepository
	stateManager *state.StateManager
	config       *utils.Config
}

func NewUsersHandler(
	billingRepo *repository.BillingRepository,
	userRepo *repository.UserRepository,
	stateManager *state.StateManager,
	config *utils.Config,
) *UsersHandler {
	return &UsersHandler{
		billingRepo:  billingRepo,
		userRepo:     userRepo,
		stateManager: stateManager,
		config:       config,
	}
}

// HandleAccountMenu shows account search menu
func (h *UsersHandler) HandleAccountMenu(ctx *handlers.HandlerContext) error {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("admin_search_contract"), "search_contract"),
			tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("admin_search_phone"), "search_phone"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("admin_search_name"), "search_name"),
			tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("admin_search_address"), "search_address"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("back"), "back"),
		),
	)

	text := ctx.Translator.Get("admin_search_menu")
	msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboard

	// Answer callback
	callbackConfig := tgbotapi.NewCallback(ctx.Update.CallbackQuery.ID, "")
	_, _ = ctx.Bot.Request(callbackConfig)

	_, err := ctx.Bot.Send(msg)
	return err
}

// HandleSearchContract initiates contract search
func (h *UsersHandler) HandleSearchContract(ctx *handlers.HandlerContext) error {
	h.stateManager.SetState(int64(ctx.Update.CallbackQuery.From.ID), state.StateSearchContract, nil)

	text := ctx.Translator.Get("admin_enter_contract")
	msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, text)

	// Answer callback
	callbackConfig := tgbotapi.NewCallback(ctx.Update.CallbackQuery.ID, "")
	_, _ = ctx.Bot.Request(callbackConfig)

	_, err := ctx.Bot.Send(msg)
	return err
}

// HandleSearchPhone initiates phone search
func (h *UsersHandler) HandleSearchPhone(ctx *handlers.HandlerContext) error {
	h.stateManager.SetState(int64(ctx.Update.CallbackQuery.From.ID), state.StateSearchPhone, nil)

	text := ctx.Translator.Get("admin_enter_phone")
	msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, text)

	// Answer callback
	callbackConfig := tgbotapi.NewCallback(ctx.Update.CallbackQuery.ID, "")
	_, _ = ctx.Bot.Request(callbackConfig)

	_, err := ctx.Bot.Send(msg)
	return err
}

// HandleSearchName initiates name search
func (h *UsersHandler) HandleSearchName(ctx *handlers.HandlerContext) error {
	h.stateManager.SetState(int64(ctx.Update.CallbackQuery.From.ID), state.StateSearchName, nil)

	text := ctx.Translator.Get("admin_enter_name")
	msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, text)

	// Answer callback
	callbackConfig := tgbotapi.NewCallback(ctx.Update.CallbackQuery.ID, "")
	_, _ = ctx.Bot.Request(callbackConfig)

	_, err := ctx.Bot.Send(msg)
	return err
}

// HandleSearchAddress initiates address search
func (h *UsersHandler) HandleSearchAddress(ctx *handlers.HandlerContext) error {
	h.stateManager.SetState(int64(ctx.Update.CallbackQuery.From.ID), state.StateSearchAddress, nil)

	text := ctx.Translator.Get("admin_enter_address")
	msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, text)

	// Answer callback
	callbackConfig := tgbotapi.NewCallback(ctx.Update.CallbackQuery.ID, "")
	_, _ = ctx.Bot.Request(callbackConfig)

	_, err := ctx.Bot.Send(msg)
	return err
}

// HandleContractInput handles contract number input
func (h *UsersHandler) HandleContractInput(ctx *handlers.HandlerContext, contract string) error {
	h.stateManager.ClearState(int64(ctx.Update.Message.From.ID))

	billingUser, err := h.billingRepo.SearchByContract(context.Background(), contract)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to search by contract")
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return err
	}

	if billingUser == nil {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("admin_user_not_found"))
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	return h.showUserInfo(ctx, billingUser, contract)
}

// HandlePhoneInput handles phone number input
func (h *UsersHandler) HandlePhoneInput(ctx *handlers.HandlerContext, phone string) error {
	h.stateManager.ClearState(int64(ctx.Update.Message.From.ID))

	billingUser, contract, err := h.billingRepo.SearchByPhone(context.Background(), phone)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to search by phone")
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return err
	}

	if billingUser == nil {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("admin_user_not_found"))
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	return h.showUserInfo(ctx, billingUser, contract)
}

// HandleNameInput handles name input
func (h *UsersHandler) HandleNameInput(ctx *handlers.HandlerContext, name string) error {
	billingUsers, err := h.billingRepo.SearchByName(context.Background(), name)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to search by name")
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return err
	}

	if len(billingUsers) == 0 {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("admin_user_not_found"))
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	if len(billingUsers) == 1 {
		h.stateManager.ClearState(int64(ctx.Update.Message.From.ID))
		return h.showUserInfo(ctx, &billingUsers[0], "")
	}

	// Multiple users found - show list
	return h.showUserList(ctx, billingUsers)
}

// HandleAddressInput handles address input
func (h *UsersHandler) HandleAddressInput(ctx *handlers.HandlerContext, address string) error {
	billingUsers, err := h.billingRepo.SearchByAddress(context.Background(), address)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to search by address")
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return err
	}

	if len(billingUsers) == 0 {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("admin_user_not_found"))
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	if len(billingUsers) == 1 {
		h.stateManager.ClearState(int64(ctx.Update.Message.From.ID))
		return h.showUserInfo(ctx, &billingUsers[0], "")
	}

	// Multiple users found - show list
	return h.showUserList(ctx, billingUsers)
}

// HandleContractSearch handles direct contract search (from text message)
func (h *UsersHandler) HandleContractSearch(ctx *handlers.HandlerContext, contract string) error {
	return h.HandleContractInput(ctx, contract)
}

// HandleAccountSelection handles account selection from list
func (h *UsersHandler) HandleAccountSelection(ctx *handlers.HandlerContext, contract string) error {
	billingUser, err := h.billingRepo.SearchByContract(context.Background(), contract)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to search by contract")
		msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return err
	}

	if billingUser == nil {
		msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, ctx.Translator.Get("admin_user_not_found"))
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	// Answer callback
	callbackConfig := tgbotapi.NewCallback(ctx.Update.CallbackQuery.ID, "")
	_, _ = ctx.Bot.Request(callbackConfig)

	return h.showUserInfo(ctx, billingUser, contract)
}

// showUserInfo shows detailed user information
func (h *UsersHandler) showUserInfo(ctx *handlers.HandlerContext, billingUser *models.BillingUser, contract string) error {
	text := h.formatUserInfo(ctx, billingUser, contract)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("admin_change_balance"), fmt.Sprintf("admin_change_balance_%d", billingUser.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("admin_temporary_payment"), fmt.Sprintf("admin_temporary_payment_%d", billingUser.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("admin_answer_user"), fmt.Sprintf("answer_%d", billingUser.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("back"), "back"),
		),
	)

	var chatID int64
	if ctx.Update.Message != nil {
		chatID = ctx.Update.Message.Chat.ID
	} else if ctx.Update.CallbackQuery != nil {
		chatID = ctx.Update.CallbackQuery.Message.Chat.ID
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	msg.ParseMode = tgbotapi.ModeHTML

	// Store selected user ID in state data
	h.stateManager.SetState(int64(ctx.Update.Message.From.ID), state.StateAccountMenuList, state.StateData{
		"selected_user_id": billingUser.ID,
		"contract":        contract,
	})

	_, err := ctx.Bot.Send(msg)
	return err
}

// showUserList shows list of found users
func (h *UsersHandler) showUserList(ctx *handlers.HandlerContext, users []models.BillingUser) error {
	var rows [][]tgbotapi.InlineKeyboardButton
	for i, user := range users {
		if i >= 10 { // Limit to 10 users
			break
		}
		contract := fmt.Sprintf("contract_%d", user.ID) // Placeholder - should get actual contract
		buttonText := fmt.Sprintf("%s (ID: %d)", user.Username, user.ID)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonText, fmt.Sprintf("account_%s", contract)),
		))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("back"), "back"),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	text := ctx.Translator.Getf("admin_users_found", len(users))
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboard

	h.stateManager.SetState(int64(ctx.Update.Message.From.ID), state.StateAccountMenuList, nil)

	_, err := ctx.Bot.Send(msg)
	return err
}

// formatUserInfo formats user information for display
func (h *UsersHandler) formatUserInfo(ctx *handlers.HandlerContext, user *models.BillingUser, contract string) string {
	var text strings.Builder
	text.WriteString(ctx.Translator.Get("admin_user_info_title"))
	text.WriteString("\n\n")
	text.WriteString(fmt.Sprintf("<b>%s:</b> %s\n", ctx.Translator.Get("admin_user_id"), strconv.FormatInt(user.ID, 10)))
	text.WriteString(fmt.Sprintf("<b>%s:</b> %s\n", ctx.Translator.Get("admin_username"), user.Username))
	text.WriteString(fmt.Sprintf("<b>%s:</b> %.2f грн\n", ctx.Translator.Get("admin_balance"), user.Balance))
	text.WriteString(fmt.Sprintf("<b>%s:</b> %s\n", ctx.Translator.Get("admin_status"), user.Status))
	if contract != "" {
		text.WriteString(fmt.Sprintf("<b>%s:</b> %s\n", ctx.Translator.Get("admin_contract"), contract))
	}
	return text.String()
}

