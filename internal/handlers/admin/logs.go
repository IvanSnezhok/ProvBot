package admin

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"provbot/internal/handlers"
	"provbot/internal/repository"
	"provbot/internal/state"
	"provbot/internal/utils"
)

type LogsHandler struct {
	stateManager *state.StateManager
	config       *utils.Config
}

func NewLogsHandler(stateManager *state.StateManager, config *utils.Config) *LogsHandler {
	return &LogsHandler{
		stateManager: stateManager,
		config:       config,
	}
}

// HandleMessageHistory initiates message history view
func (h *LogsHandler) HandleMessageHistory(ctx *handlers.HandlerContext) error {
	h.stateManager.SetState(int64(ctx.Update.CallbackQuery.From.ID), state.StateMessageHistory, nil)

	// Request number of messages to show
	text := ctx.Translator.Get("admin_enter_message_count")
	msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, text)

	// Answer callback
	callbackConfig := tgbotapi.NewCallback(ctx.Update.CallbackQuery.ID, "")
	_, _ = ctx.Bot.Request(callbackConfig)

	_, err := ctx.Bot.Send(msg)
	return err
}

// HandleHistoryCount handles message count input
func (h *LogsHandler) HandleHistoryCount(ctx *handlers.HandlerContext, countText string) error {
	count := 10 // Default
	var err error

	if countText != "" {
		count, err = strconv.Atoi(countText)
		if err != nil || count < 1 || count > 100 {
			msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("admin_invalid_count"))
			_, _ = ctx.Bot.Send(msg)
			return err
		}
	}

	// Get message logs
	logRepo := repository.NewLogRepository()
	logs, err := logRepo.GetMessageLogs(context.Background(), int64(ctx.Update.Message.From.ID), count)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to get message logs")
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return err
	}

	if len(logs) == 0 {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("admin_no_messages"))
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	// Format and send logs
	var text string
	for i, log := range logs {
		if i >= count {
			break
		}
		direction := "→"
		if log.Direction == "incoming" {
			direction = "←"
		}
		text += fmt.Sprintf("%s [%s] %s\n", direction, log.CreatedAt.Format("2006-01-02 15:04:05"), log.MessageText)
		if i < len(logs)-1 {
			text += "\n"
		}
	}

	h.stateManager.ClearState(int64(ctx.Update.Message.From.ID))

	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
	_, err = ctx.Bot.Send(msg)
	return err
}

