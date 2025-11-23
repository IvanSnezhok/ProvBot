package user

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"provbot/internal/handlers"
	"provbot/internal/i18n"
	"provbot/internal/repository"
	"provbot/internal/service"
	"provbot/internal/state"
	"provbot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PayBillHandler struct {
	billingService *service.BillingService
	userRepo       *repository.UserRepository
	stateManager   *state.StateManager
	config         *utils.Config
}

func NewPayBillHandler(
	billingService *service.BillingService,
	userRepo *repository.UserRepository,
	stateManager *state.StateManager,
	config *utils.Config,
) *PayBillHandler {
	return &PayBillHandler{
		billingService: billingService,
		userRepo:       userRepo,
		stateManager:   stateManager,
		config:         config,
	}
}

// HandleTopUp handles "Поповнити рахунок" button/text
func (h *PayBillHandler) HandleTopUp(ctx *handlers.HandlerContext) error {
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

	// Get tariff for user
	tariffName, _, err := h.billingService.GetTariff(context.Background(), contract)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to get tariff")
		return err
	}

	// Handle predefined tariffs
	// Note: Python code checks for specific tariff names. This logic might be fragile if names change.
	// Consider moving this config to DB or config file.
	if tariffName == "СТАНДАРТ(180грн)." {
		return h.handlePredefinedTariff(ctx, 180, 1080, contract)
	} else if tariffName == "PON-100(200грн)" || tariffName == "VIP WIFI-200" {
		return h.handlePredefinedTariff(ctx, 200, 1200, contract)
	} else if tariffName == "PON-300(350грн)" {
		return h.handlePredefinedTariff(ctx, 350, 2100, contract)
	}

	// For other tariffs, ask for amount
	return h.requestCustomAmount(ctx, contract)
}

// handlePredefinedTariff sends predefined invoice options
func (h *PayBillHandler) handlePredefinedTariff(ctx *handlers.HandlerContext, monthlyAmount, sixMonthAmount int, contract string) error {
	// Send info message
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("topup_notice"))
	_, _ = ctx.Bot.Send(msg)

	// Send monthly invoice
	invoice1 := h.createInvoice(monthlyAmount*100, contract, ctx.Translator) // Amount in kopecks
	_, err := ctx.Bot.Send(invoice1)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to send invoice")
	}

	// Send promotion message
	promoMsg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("topup_promotion"))
	_, _ = ctx.Bot.Send(promoMsg)

	// Send 6-month invoice
	invoice2 := h.createInvoice(sixMonthAmount*100, contract, ctx.Translator) // Amount in kopecks
	_, err = ctx.Bot.Send(invoice2)
	return err
}

// requestCustomAmount requests custom amount from user
func (h *PayBillHandler) requestCustomAmount(ctx *handlers.HandlerContext, contract string) error {
	text := ctx.Translator.Get("topup_custom_amount")
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)

	// Add back button
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(ctx.Translator.Get("back")),
		),
	)
	keyboard.ResizeKeyboard = true
	msg.ReplyMarkup = keyboard

	h.stateManager.SetState(ctx.Update.Message.From.ID, state.StateWaitingInvoicePayload, state.StateData{
		"contract": contract,
	})

	_, err := ctx.Bot.Send(msg)
	return err
}

// HandleAmountInput handles user input for payment amount
func (h *PayBillHandler) HandleAmountInput(ctx *handlers.HandlerContext) error {
	userState, data, exists := h.stateManager.GetState(ctx.Update.Message.From.ID)
	if !exists || userState != state.StateWaitingInvoicePayload {
		return nil
	}

	amountStr := strings.TrimSpace(ctx.Update.Message.Text)
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("invalid_amount"))
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	// Store amount in state
	h.stateManager.UpdateData(ctx.Update.Message.From.ID, state.StateData{
		"amount": amount,
	})

	// Request contract number
	text := ctx.Translator.Get("topup_enter_contract")
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
	h.stateManager.SetState(ctx.Update.Message.From.ID, state.StateWaitingInvoiceContract, data)

	_, err = ctx.Bot.Send(msg)
	return err
}

// HandleContractInput handles contract number input
func (h *PayBillHandler) HandleContractInput(ctx *handlers.HandlerContext) error {
	userState, data, exists := h.stateManager.GetState(ctx.Update.Message.From.ID)
	if !exists || userState != state.StateWaitingInvoiceContract {
		return nil
	}

	contract := strings.TrimSpace(ctx.Update.Message.Text)

	// Validate contract format (8 digits)
	matched, _ := regexp.MatchString(`^\d{8}$`, contract)
	if !matched {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("invalid_contract_format"))
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	// Check if contract exists
	existsContract, err := h.billingService.CheckContract(context.Background(), contract)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to check contract")
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return err
	}

	if !existsContract {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("contract_not_found"))
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	// Get amount from state
	amount, ok := data["amount"].(float64)
	if !ok {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	// Create invoice
	invoice := h.createCustomInvoice(int(amount*100), contract, ctx.Translator)

	// Notify admins
	h.notifyAdminsAboutInvoice(ctx, contract, amount)

	_, err = ctx.Bot.Send(invoice)
	if err != nil {
		// Handle currency amount invalid error
		if strings.Contains(err.Error(), "CurrencyTotalAmountInvalid") {
			msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("invalid_amount_minimum"))
			_, _ = ctx.Bot.Send(msg)
			h.stateManager.SetState(ctx.Update.Message.From.ID, state.StateInvalidPayload, data)
			return nil
		}
		return err
	}

	h.stateManager.ClearState(ctx.Update.Message.From.ID)
	return nil
}

// createInvoice creates Telegram invoice for payment
func (h *PayBillHandler) createInvoice(amountKopecks int, contract string, translator i18n.Translator) tgbotapi.InvoiceConfig {
	amountGrn := float64(amountKopecks) / 100.0

	return tgbotapi.InvoiceConfig{
		Title:         translator.Getf("invoice_title", amountGrn),
		Description:   translator.Getf("invoice_description", contract, amountGrn),
		Payload:       fmt.Sprintf("%d", amountKopecks),
		ProviderToken: h.config.ProviderToken,
		Currency:      "UAH",
		Prices: []tgbotapi.LabeledPrice{
			{
				Label:  translator.Getf("invoice_label", contract),
				Amount: amountKopecks,
			},
		},
	}
}

// createCustomInvoice creates custom invoice
func (h *PayBillHandler) createCustomInvoice(amountKopecks int, contract string, translator i18n.Translator) tgbotapi.InvoiceConfig {
	return h.createInvoice(amountKopecks, contract, translator)
}

// notifyAdminsAboutInvoice notifies admins about created invoice
func (h *PayBillHandler) notifyAdminsAboutInvoice(ctx *handlers.HandlerContext, contract string, amount float64) {
	message := fmt.Sprintf("Створено інвойс для користувача %d:\nДоговір: %s\nСума: %.2f грн\n",
		ctx.Update.Message.From.ID, contract, amount)

	for _, adminID := range h.config.AdminTelegramIDs {
		msg := tgbotapi.NewMessage(adminID, message)
		_, _ = ctx.Bot.Send(msg)
	}
}
