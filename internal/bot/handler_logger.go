package bot

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"provbot/internal/utils"
)

// HandlerLoggingMiddleware logs handler execution
func HandlerLoggingMiddleware() MiddlewareFunc {
	return func(ctx *BotContext, next HandlerFunc) error {
		startTime := time.Now()
		handlerName := getHandlerName(ctx.Update)

		// Log handler start
		logFields := map[string]interface{}{
			"handler": handlerName,
			"type":    getUpdateType(ctx.Update),
		}

		if ctx.Update.Message != nil {
			logFields["chat_id"] = ctx.Update.Message.Chat.ID
			logFields["user_id"] = ctx.Update.Message.From.ID
			logFields["username"] = ctx.Update.Message.From.UserName
			if ctx.Update.Message.Text != "" {
				logFields["text"] = truncateString(ctx.Update.Message.Text, 100)
			}
		}

		if ctx.Update.CallbackQuery != nil {
			logFields["chat_id"] = ctx.Update.CallbackQuery.Message.Chat.ID
			logFields["user_id"] = ctx.Update.CallbackQuery.From.ID
			logFields["callback_data"] = ctx.Update.CallbackQuery.Data
		}

		utils.Logger.WithFields(logFields).Info("Handler started")

		// Execute handler
		err := next(ctx)

		// Log handler completion
		duration := time.Since(startTime)
		logFields["duration_ms"] = duration.Milliseconds()
		logFields["success"] = err == nil

		if err != nil {
			logFields["error"] = err.Error()
			utils.Logger.WithFields(logFields).Error("Handler completed with error")
		} else {
			utils.Logger.WithFields(logFields).Info("Handler completed successfully")
		}

		return err
	}
}

// getHandlerName extracts handler name from update
func getHandlerName(update *tgbotapi.Update) string {
	if update.Message != nil {
		if update.Message.IsCommand() {
			return update.Message.Command()
		}
		if update.Message.Text != "" {
			// Check for common text commands
			text := update.Message.Text
			switch text {
			case "Поповнити рахунок", "Top up account":
				return "topup"
			case "Чат з тех. підтримкою", "Support chat":
				return "support"
			case "Тимчасовий платіж", "Temporary payment":
				return "time_pay"
			case "Назад", "Back":
				return "back"
			default:
				return "text_message"
			}
		}
		if update.Message.Contact != nil {
			return "contact"
		}
	}
	if update.CallbackQuery != nil {
		return "callback_" + truncateString(update.CallbackQuery.Data, 50)
	}
	if update.PreCheckoutQuery != nil {
		return "pre_checkout"
	}
	if update.Message != nil && update.Message.SuccessfulPayment != nil {
		return "successful_payment"
	}
	return "unknown"
}

// getUpdateType returns the type of update
func getUpdateType(update *tgbotapi.Update) string {
	if update.Message != nil {
		return "message"
	}
	if update.CallbackQuery != nil {
		return "callback_query"
	}
	if update.PreCheckoutQuery != nil {
		return "pre_checkout_query"
	}
	if update.Message != nil && update.Message.SuccessfulPayment != nil {
		return "successful_payment"
	}
	return "unknown"
}

// truncateString truncates a string to max length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

