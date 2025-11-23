package admin

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"provbot/internal/handlers"
	"provbot/internal/repository"
	"provbot/internal/service"
	"provbot/internal/state"
	"provbot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BillingHandler struct {
	billingService *service.BillingService
	userRepo       *repository.UserRepository
	stateManager   *state.StateManager
	config         *utils.Config
}

func NewBillingHandler(
	billingService *service.BillingService,
	userRepo *repository.UserRepository,
	stateManager *state.StateManager,
	config *utils.Config,
) *BillingHandler {
	return &BillingHandler{
		billingService: billingService,
		userRepo:       userRepo,
		stateManager:   stateManager,
		config:         config,
	}
}

// HandleBalanceChange initiates balance change
func (h *BillingHandler) HandleBalanceChange(ctx *handlers.HandlerContext, callbackData string) error {
	// Extract user ID from callback data
	parts := strings.Split(callbackData, "_")
	if len(parts) < 4 {
		return fmt.Errorf("invalid callback data format")
	}

	userID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return err
	}

	// Store user ID in state
	h.stateManager.SetState(int64(ctx.Update.CallbackQuery.From.ID), state.StateAdminChangeBalance, state.StateData{
		"user_id": userID,
	})

	text := ctx.Translator.Get("admin_enter_balance_amount")
	msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, text)

	// Answer callback
	callbackConfig := tgbotapi.NewCallback(ctx.Update.CallbackQuery.ID, "")
	_, _ = ctx.Bot.Request(callbackConfig)

	_, err = ctx.Bot.Send(msg)
	return err
}

// HandleBalanceChangeInput handles balance amount input
func (h *BillingHandler) HandleBalanceChangeInput(ctx *handlers.HandlerContext, amountText string) error {
	userState, stateData, exists := h.stateManager.GetState(int64(ctx.Update.Message.From.ID))
	if !exists || userState != state.StateAdminChangeBalance {
		return nil
	}

	userID, ok := stateData["user_id"].(int64)
	if !ok {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return fmt.Errorf("user ID not found in state")
	}

	amount, err := strconv.ParseFloat(amountText, 64)
	if err != nil {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("admin_invalid_amount"))
		_, _ = ctx.Bot.Send(msg)
		return err
	}

	// Update balance
	err = h.billingService.UpdateBalance(context.Background(), userID, amount)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to update balance")
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return err
	}

	// Get updated user info
	billingUser, err := h.billingService.GetUserByID(context.Background(), userID)
	if err == nil && billingUser != nil {
		text := ctx.Translator.Getf("admin_balance_updated", amount, billingUser.Balance)
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
		_, _ = ctx.Bot.Send(msg)
	}

	h.stateManager.ClearState(int64(ctx.Update.Message.From.ID))
	return nil
}

// HandleTemporaryPayment handles temporary payment activation
func (h *BillingHandler) HandleTemporaryPayment(ctx *handlers.HandlerContext, callbackData string) error {
	// Extract user ID from callback data
	parts := strings.Split(callbackData, "_")
	if len(parts) < 4 {
		return fmt.Errorf("invalid callback data format")
	}

	userID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return err
	}

	// Get user contract from state or billing
	_, stateData, exists := h.stateManager.GetState(int64(ctx.Update.CallbackQuery.From.ID))
	var contract string
	if exists {
		if c, ok := stateData["contract"].(string); ok {
			contract = c
		}
	}

	// If no contract in state, try to get from billing user
	if contract == "" {
		billingUser, err := h.billingService.GetUserByID(context.Background(), userID)
		if err == nil && billingUser != nil {
			contract = billingUser.Contract
		}
	}

	if contract == "" {
		msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, ctx.Translator.Get("admin_contract_not_found"))
		_, _ = ctx.Bot.Send(msg)
		return fmt.Errorf("contract not found")
	}

	// Enable temporary payment
	success, err := h.billingService.TemporaryPay(context.Background(), contract)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to enable temporary payment")
		msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return err
	}

	// Answer callback
	callbackConfig := tgbotapi.NewCallback(ctx.Update.CallbackQuery.ID, "")
	_, _ = ctx.Bot.Request(callbackConfig)

	if success {
		// We don't have amount returned from TemporaryPay anymore, just success
		text := ctx.Translator.Get("admin_temporary_payment_success") // Removed format arg
		msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, text)
		_, _ = ctx.Bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, ctx.Translator.Get("admin_temporary_payment_failed"))
		_, _ = ctx.Bot.Send(msg)
	}

	return nil
}
