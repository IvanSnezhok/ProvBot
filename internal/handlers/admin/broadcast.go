package admin

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"provbot/internal/handlers"
	"provbot/internal/repository"
	"provbot/internal/state"
	"provbot/internal/utils"
)

type BroadcastHandler struct {
	userRepo     *repository.UserRepository
	stateManager *state.StateManager
	config       *utils.Config
}

func NewBroadcastHandler(
	userRepo *repository.UserRepository,
	stateManager *state.StateManager,
	config *utils.Config,
) *BroadcastHandler {
	return &BroadcastHandler{
		userRepo:     userRepo,
		stateManager: stateManager,
		config:       config,
	}
}

// HandleSendMessage initiates message sending
func (h *BroadcastHandler) HandleSendMessage(ctx *handlers.HandlerContext) error {
	h.stateManager.SetState(int64(ctx.Update.CallbackQuery.From.ID), state.StateSendMessagePhone, nil)

	text := ctx.Translator.Get("admin_enter_phone_or_userid")
	msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, text)

	// Answer callback
	callbackConfig := tgbotapi.NewCallback(ctx.Update.CallbackQuery.ID, "")
	_, _ = ctx.Bot.Request(callbackConfig)

	_, err := ctx.Bot.Send(msg)
	return err
}

// HandlePhoneInput handles phone number or user ID input
func (h *BroadcastHandler) HandlePhoneInput(ctx *handlers.HandlerContext, input string) error {
	var userID int64
	var err error

	// Try to parse as user ID first
	userID, err = strconv.ParseInt(input, 10, 64)
	if err != nil {
		// If not a number, treat as phone number
		// Search for user by phone
		users, err := h.userRepo.GetAllActiveUsers(context.Background())
		if err != nil {
			msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
			_, _ = ctx.Bot.Send(msg)
			return err
		}

		// Find user by phone
		found := false
		for _, u := range users {
			if u.PhoneNumber != nil && strings.Contains(*u.PhoneNumber, input) {
				userID = int64(u.TelegramID)
				found = true
				break
			}
		}

		if !found {
			msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("admin_user_not_found"))
			_, _ = ctx.Bot.Send(msg)
			return fmt.Errorf("user not found")
		}
	}

	// Store user ID in state and switch to message text state
	h.stateManager.SetState(int64(ctx.Update.Message.From.ID), state.StateSendMessageText, state.StateData{
		"target_user_id": userID,
	})

	text := ctx.Translator.Get("admin_enter_message_text")
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
	_, err = ctx.Bot.Send(msg)
	return err
}

// HandleMessageInput handles message text/media input
func (h *BroadcastHandler) HandleMessageInput(ctx *handlers.HandlerContext, text string) error {
	userState, stateData, exists := h.stateManager.GetState(int64(ctx.Update.Message.From.ID))
	if !exists || userState != state.StateSendMessageText {
		return nil
	}

	targetUserID, ok := stateData["target_user_id"].(int64)
	if !ok {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return fmt.Errorf("target user ID not found in state")
	}

	// Send message to user
	if ctx.Update.Message.Photo != nil && len(ctx.Update.Message.Photo) > 0 {
		// Send photo
		photo := ctx.Update.Message.Photo[len(ctx.Update.Message.Photo)-1]
		photoMsg := tgbotapi.NewPhoto(targetUserID, tgbotapi.FileID(photo.FileID))
		photoMsg.Caption = text
		_, err := ctx.Bot.Send(photoMsg)
		if err != nil {
			utils.Logger.WithError(err).Error("Failed to send photo")
		}
	} else if ctx.Update.Message.Document != nil {
		// Send document
		docMsg := tgbotapi.NewDocument(targetUserID, tgbotapi.FileID(ctx.Update.Message.Document.FileID))
		docMsg.Caption = text
		_, err := ctx.Bot.Send(docMsg)
		if err != nil {
			utils.Logger.WithError(err).Error("Failed to send document")
		}
	} else if ctx.Update.Message.Video != nil {
		// Send video
		videoMsg := tgbotapi.NewVideo(targetUserID, tgbotapi.FileID(ctx.Update.Message.Video.FileID))
		videoMsg.Caption = text
		_, err := ctx.Bot.Send(videoMsg)
		if err != nil {
			utils.Logger.WithError(err).Error("Failed to send video")
		}
	} else {
		// Send text message
		msg := tgbotapi.NewMessage(targetUserID, text)
		_, err := ctx.Bot.Send(msg)
		if err != nil {
			utils.Logger.WithError(err).Error("Failed to send message")
			msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
			_, _ = ctx.Bot.Send(msg)
			return err
		}
	}

	// Notify admins
	h.notifyAdmins(ctx, targetUserID, text)

	h.stateManager.ClearState(int64(ctx.Update.Message.From.ID))

	successMsg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("admin_message_sent"))
	_, _ = ctx.Bot.Send(successMsg)
	return nil
}

// HandleAnswer initiates answer to user
func (h *BroadcastHandler) HandleAnswer(ctx *handlers.HandlerContext, callbackData string) error {
	// Extract user ID from callback data (format: "answer_<user_id>")
	parts := strings.Split(callbackData, "_")
	if len(parts) < 2 {
		return fmt.Errorf("invalid callback data format")
	}

	userID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return err
	}

	// Store user ID in state
	h.stateManager.SetState(int64(ctx.Update.CallbackQuery.From.ID), state.StateAnswer, state.StateData{
		"target_user_id": userID,
	})

	text := ctx.Translator.Get("admin_enter_answer_text")
	msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, text)

	// Answer callback
	callbackConfig := tgbotapi.NewCallback(ctx.Update.CallbackQuery.ID, "")
	_, _ = ctx.Bot.Request(callbackConfig)

	_, err = ctx.Bot.Send(msg)
	return err
}

// HandleAnswerInput handles answer text input
func (h *BroadcastHandler) HandleAnswerInput(ctx *handlers.HandlerContext, text string, stateData state.StateData) error {
	targetUserID, ok := stateData["target_user_id"].(int64)
	if !ok {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return fmt.Errorf("target user ID not found in state")
	}

	// Send answer to user
	msg := tgbotapi.NewMessage(targetUserID, text)
	_, err := ctx.Bot.Send(msg)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to send answer")
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return err
	}

	// Notify admins
	h.notifyAdmins(ctx, targetUserID, text)

	h.stateManager.ClearState(int64(ctx.Update.Message.From.ID))

	successMsg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("admin_answer_sent"))
	_, _ = ctx.Bot.Send(successMsg)
	return nil
}

// notifyAdmins notifies all admins about sent message
func (h *BroadcastHandler) notifyAdmins(ctx *handlers.HandlerContext, targetUserID int64, message string) {
	adminMessage := fmt.Sprintf("Повідомлення відправлено користувачу ID: %d\nТекст: %s", targetUserID, message)
	for _, adminID := range h.config.AdminTelegramIDs {
		msg := tgbotapi.NewMessage(adminID, adminMessage)
		_, _ = ctx.Bot.Send(msg)
	}
}

