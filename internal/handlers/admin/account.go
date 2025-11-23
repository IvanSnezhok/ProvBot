package admin

import (
	"context"
	"fmt"

	"provbot/internal/handlers"
	"provbot/internal/repository"
	"provbot/internal/service"
	"provbot/internal/state"
	"provbot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AccountHandler struct {
	billingService *service.BillingService
	userRepo       *repository.UserRepository
	stateManager   *state.StateManager
	config         *utils.Config
}

func NewAccountHandler(
	billingService *service.BillingService,
	userRepo *repository.UserRepository,
	stateManager *state.StateManager,
	config *utils.Config,
) *AccountHandler {
	return &AccountHandler{
		billingService: billingService,
		userRepo:       userRepo,
		stateManager:   stateManager,
		config:         config,
	}
}

// HandleAccountMenu shows account menu for selected user
func (h *AccountHandler) HandleAccountMenu(ctx *handlers.HandlerContext) error {
	userState, stateData, exists := h.stateManager.GetState(int64(ctx.Update.CallbackQuery.From.ID))
	if !exists || userState != state.StateAccountMenuList {
		return nil
	}

	userID, ok := stateData["selected_user_id"].(int64)
	if !ok {
		msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return fmt.Errorf("user ID not found in state")
	}

	billingUser, err := h.billingService.GetUserByID(context.Background(), userID)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to get billing user")
		msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return err
	}

	if billingUser == nil {
		msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, ctx.Translator.Get("admin_user_not_found"))
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	contract := ""
	if c, ok := stateData["contract"].(string); ok {
		contract = c
	}

	// Show user info (reuse UsersHandler method)
	usersHandler := NewUsersHandler(h.billingService, h.userRepo, h.stateManager, h.config)
	return usersHandler.showUserInfo(ctx, billingUser, contract)
}
