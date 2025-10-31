package user

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"provbot/internal/handlers"
	"provbot/internal/repository"
	"provbot/internal/utils"
)

type TimePayHandler struct {
	billingRepo *repository.BillingRepository
	userRepo    *repository.UserRepository
	config      *utils.Config
}

func NewTimePayHandler(
	billingRepo *repository.BillingRepository,
	userRepo *repository.UserRepository,
	config *utils.Config,
) *TimePayHandler {
	return &TimePayHandler{
		billingRepo: billingRepo,
		userRepo:    userRepo,
		config:      config,
	}
}

// HandleTimePay handles "Тимчасовий платіж" button/text
func (h *TimePayHandler) HandleTimePay(ctx *handlers.HandlerContext) error {
	user := ctx.User
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Check if user is banned
	// ban := await db.get_ban() - need to implement ban check
	
	// Get user contract
	if user.Contract == nil || *user.Contract == "" {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("no_contract"))
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	contract := *user.Contract
	
	// Enable temporary payment
	success, amount, err := h.billingRepo.EnableTemporaryPayment(context.Background(), contract)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to enable temporary payment")
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return err
	}

	if success {
		// Success message
		text := ctx.Translator.Getf("time_pay_success", amount)
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
		
		// Add back button
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(ctx.Translator.Get("menu_back")),
			),
		)
		keyboard.ResizeKeyboard = true
		msg.ReplyMarkup = keyboard
		
		_, err = ctx.Bot.Send(msg)
		
		// Notify admins
		h.notifyAdmins(ctx, contract)
		
		return err
	} else {
		// Failed - already used this month
		text := ctx.Translator.Get("time_pay_failed")
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
		
		// Add back button
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(ctx.Translator.Get("menu_back")),
			),
		)
		keyboard.ResizeKeyboard = true
		msg.ReplyMarkup = keyboard
		
		_, err = ctx.Bot.Send(msg)
		return err
	}
}

// notifyAdmins notifies admins about temporary payment usage
func (h *TimePayHandler) notifyAdmins(ctx *handlers.HandlerContext, contract string) {
	message := fmt.Sprintf("Користувач %s використав тимчасовий платіж!", contract)

	for _, adminID := range h.config.AdminTelegramIDs {
		msg := tgbotapi.NewMessage(adminID, message)
		_, _ = ctx.Bot.Send(msg)
	}
}

