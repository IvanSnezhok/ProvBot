package admin

import (
	"strings"

	"provbot/internal/handlers"
)

// HandleCallbackQuery handles admin callback queries
func HandleCallbackQuery(ctx *handlers.HandlerContext, panelHandler *PanelHandler) error {
	callback := ctx.Update.CallbackQuery
	if callback == nil {
		return nil
	}

	data := callback.Data

	// Handle back button
	if data == "back" {
		// Clear state and show admin panel
		ctx.StateManager.ClearState(int64(callback.From.ID))
		return panelHandler.ShowAdminPanel(ctx)
	}

	// Handle account menu
	if data == "account_menu" {
		return panelHandler.UsersHandler.HandleAccountMenu(ctx)
	}

	// Handle search methods
	if data == "search_contract" {
		return panelHandler.UsersHandler.HandleSearchContract(ctx)
	}
	if data == "search_phone" {
		return panelHandler.UsersHandler.HandleSearchPhone(ctx)
	}
	if data == "search_name" {
		return panelHandler.UsersHandler.HandleSearchName(ctx)
	}
	if data == "search_address" {
		return panelHandler.UsersHandler.HandleSearchAddress(ctx)
	}

	// Handle account selection (format: "account_<contract>")
	if strings.HasPrefix(data, "account_") {
		contract := strings.TrimPrefix(data, "account_")
		return panelHandler.UsersHandler.HandleAccountSelection(ctx, contract)
	}

	// Handle balance change (format: "admin_change_balance_<user_id>")
	if strings.HasPrefix(data, "admin_change_balance_") {
		return panelHandler.BillingHandler.HandleBalanceChange(ctx, data)
	}

	// Handle temporary payment (format: "admin_temporary_payment_<user_id>")
	if strings.HasPrefix(data, "admin_temporary_payment_") {
		return panelHandler.BillingHandler.HandleTemporaryPayment(ctx, data)
	}

	// Handle send message
	if data == "panel_send_message" {
		return panelHandler.BroadcastHandler.HandleSendMessage(ctx)
	}

	// Handle answer (format: "answer_<user_id>")
	if strings.HasPrefix(data, "answer_") {
		return panelHandler.BroadcastHandler.HandleAnswer(ctx, data)
	}

	// Handle message history
	if data == "message_history" {
		return panelHandler.LogsHandler.HandleMessageHistory(ctx)
	}

	// Handle account menu (from user selection)
	if data == "admin_account_menu" {
		return panelHandler.AccountHandler.HandleAccountMenu(ctx)
	}

	return nil
}

