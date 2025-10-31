package admin

import (
	"context"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"provbot/internal/handlers"
	"provbot/internal/repository"
	"provbot/internal/service"
	"provbot/internal/state"
	"provbot/internal/utils"
)

type PanelHandler struct {
	adminService   *service.AdminService
	billingRepo    *repository.BillingRepository
	userRepo       *repository.UserRepository
	stateManager   *state.StateManager
	config         *utils.Config
	UsersHandler   *UsersHandler
	BillingHandler *BillingHandler
	BroadcastHandler *BroadcastHandler
	AccountHandler *AccountHandler
	LogsHandler    *LogsHandler
}

func NewPanelHandler(
	adminService *service.AdminService,
	billingRepo *repository.BillingRepository,
	userRepo *repository.UserRepository,
	stateManager *state.StateManager,
	config *utils.Config,
) *PanelHandler {
	ph := &PanelHandler{
		adminService: adminService,
		billingRepo: billingRepo,
		userRepo:    userRepo,
		stateManager: stateManager,
		config:      config,
	}

	// Initialize sub-handlers
	ph.UsersHandler = NewUsersHandler(billingRepo, userRepo, stateManager, config)
	ph.BillingHandler = NewBillingHandler(billingRepo, userRepo, stateManager, config)
	ph.BroadcastHandler = NewBroadcastHandler(userRepo, stateManager, config)
	ph.AccountHandler = NewAccountHandler(billingRepo, userRepo, stateManager, config)
	ph.LogsHandler = NewLogsHandler(stateManager, config)

	return ph
}

// HandleAdmin handles /admin command or "Назад" text for admins
func (h *PanelHandler) HandleAdmin(ctx *handlers.HandlerContext) error {
	if !h.config.IsAdmin(ctx.Update.Message.From.ID) {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("admin_unauthorized"))
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	// Check if text is a contract number (8 digits)
	text := strings.TrimSpace(ctx.Update.Message.Text)
	if isContractNumber(text) {
		// Redirect to user search
		return h.UsersHandler.HandleContractSearch(ctx, text)
	}

	// Show admin panel menu
	return h.ShowAdminPanel(ctx)
}

// ShowAdminPanel shows admin panel keyboard
func (h *PanelHandler) ShowAdminPanel(ctx *handlers.HandlerContext) error {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("admin_users"), "account_menu"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("admin_broadcast"), "panel_send_message"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("admin_message_history"), "message_history"),
		),
	)

	text := ctx.Translator.Get("admin_panel_title")
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboard
	_, err := ctx.Bot.Send(msg)
	return err
}

// HandleStats handles /stats command
func (h *PanelHandler) HandleStats(ctx *handlers.HandlerContext) error {
	if !h.config.IsAdmin(ctx.Update.Message.From.ID) {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("admin_unauthorized"))
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	// Get user count
	users, err := h.userRepo.GetAllActiveUsers(context.Background())
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to get user count")
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return err
	}

	// Count users with contracts
	contractCount := 0
	for _, u := range users {
		if u.Contract != nil && *u.Contract != "" {
			contractCount++
		}
	}

	text := ctx.Translator.Getf("admin_stats", len(users), contractCount)
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
	_, err = ctx.Bot.Send(msg)
	return err
}

// HandleTextMessage handles text messages for admin panel
func (h *PanelHandler) HandleTextMessage(ctx *handlers.HandlerContext) error {
	if !h.config.IsAdmin(ctx.Update.Message.From.ID) {
		return nil
	}

	text := ctx.Update.Message.Text
	userState, stateData, exists := h.stateManager.GetState(int64(ctx.Update.Message.From.ID))

	// Handle "Назад" button
	if text == ctx.Translator.Get("menu_back") {
		h.stateManager.ClearState(int64(ctx.Update.Message.From.ID))
		return h.ShowAdminPanel(ctx)
	}

	// Handle states
	if exists {
		switch userState {
		case state.StateSearchContract:
			return h.UsersHandler.HandleContractInput(ctx, text)
		case state.StateSearchPhone:
			return h.UsersHandler.HandlePhoneInput(ctx, text)
		case state.StateSearchName:
			return h.UsersHandler.HandleNameInput(ctx, text)
		case state.StateSearchAddress:
			return h.UsersHandler.HandleAddressInput(ctx, text)
		case state.StateAdminChangeBalance:
			return h.BillingHandler.HandleBalanceChangeInput(ctx, text)
		case state.StateSendMessagePhone:
			return h.BroadcastHandler.HandlePhoneInput(ctx, text)
		case state.StateSendMessageText:
			return h.BroadcastHandler.HandleMessageInput(ctx, text)
		case state.StateAnswer:
			return h.BroadcastHandler.HandleAnswerInput(ctx, text, stateData)
		case state.StateMessageHistory:
			return h.LogsHandler.HandleHistoryCount(ctx, text)
		}
	}

	// Default: check if it's a contract number
	if isContractNumber(text) {
		return h.UsersHandler.HandleContractSearch(ctx, text)
	}

	// Show admin panel
	return h.ShowAdminPanel(ctx)
}

// isContractNumber checks if text is a valid contract number (8 digits)
func isContractNumber(text string) bool {
	if len(text) != 8 {
		return false
	}
	_, err := strconv.Atoi(text)
	return err == nil
}

