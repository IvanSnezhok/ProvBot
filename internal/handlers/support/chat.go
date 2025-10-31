package support

import (
	"context"
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"provbot/internal/i18n"
	"provbot/internal/models"
	"provbot/internal/repository"
	"provbot/internal/state"
	"provbot/internal/utils"
)

// ActiveChats stores active support chats: userID -> adminID
var ActiveChats = struct {
	sync.RWMutex
	chats map[int64]int64
}{
	chats: make(map[int64]int64),
}

type SupportChatHandler struct {
	logRepo      *repository.LogRepository
	stateManager *state.StateManager
	config       *utils.Config
	bot          *tgbotapi.BotAPI
}

func NewSupportChatHandler(
	logRepo *repository.LogRepository,
	stateManager *state.StateManager,
	config *utils.Config,
	bot *tgbotapi.BotAPI,
) *SupportChatHandler {
	return &SupportChatHandler{
		logRepo:      logRepo,
		stateManager: stateManager,
		config:       config,
		bot:          bot,
	}
}

// SupportContext holds context for support handlers
type SupportContext struct {
	Update       *tgbotapi.Update
	User         *models.User
	Translator   i18n.Translator
	Config       *utils.Config
	StateManager *state.StateManager
}

// HandleStartSupport handles "Чат з тех. підтримкою" button/text
func (h *SupportChatHandler) HandleStartSupport(ctx *SupportContext) error {
	userID := ctx.Update.Message.From.ID

	// Check if user already in active chat
	ActiveChats.RLock()
	_, exists := ActiveChats.chats[int64(userID)]
	ActiveChats.RUnlock()

	if exists {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("support_already_active"))
		_, _ = h.bot.Send(msg)
		return nil
	}

	// Set waiting state
	h.stateManager.SetState(int64(userID), state.StateWaitingForSupport, nil)

	// Send message to user
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("support_waiting"))
	
	// Add back button
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(ctx.Translator.Get("menu_back")),
		),
	)
	keyboard.ResizeKeyboard = true
	msg.ReplyMarkup = keyboard
	
	_, err := h.bot.Send(msg)

	// Notify admins
	h.notifyAdmins(ctx, userID)

	return err
}

// HandleSupportMessage handles messages in support chat
func (h *SupportChatHandler) HandleSupportMessage(ctx *SupportContext) error {
	userID := ctx.Update.Message.From.ID
	userState, _, exists := h.stateManager.GetState(int64(userID))

	// Check if user is waiting or chatting
	if !exists || (userState != state.StateWaitingForSupport && userState != state.StateChatting) {
		return nil
	}

	ActiveChats.RLock()
	adminID, hasAdmin := ActiveChats.chats[int64(userID)]
	ActiveChats.RUnlock()

	if !hasAdmin {
		// Still waiting for admin
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("support_waiting_manager"))
		_, _ = h.bot.Send(msg)
		return nil
	}

	// Forward message to admin
	text := fmt.Sprintf("Повідомлення від користувача: %s", ctx.Update.Message.Text)
	msg := tgbotapi.NewMessage(adminID, text)
	_, err := h.bot.Send(msg)

	// Log message
	h.logSupportMessage(int64(userID), adminID, ctx.Update.Message.Text, true)

	return err
}

// HandleAdminConnect handles /connect command for admins
func (h *SupportChatHandler) HandleAdminConnect(ctx *SupportContext, userID int64) error {
	adminID := ctx.Update.Message.From.ID

	// Check if user already has an admin
	ActiveChats.Lock()
	if _, exists := ActiveChats.chats[userID]; exists {
		ActiveChats.Unlock()
		msg := tgbotapi.NewMessage(adminID, ctx.Translator.Get("support_already_connected"))
		_, _ = h.bot.Send(msg)
		return nil
	}

	// Connect admin to user
	ActiveChats.chats[userID] = adminID
	ActiveChats.Unlock()

	// Update user state
	h.stateManager.SetState(userID, state.StateChatting, nil)

	// Notify user
	msg := tgbotapi.NewMessage(userID, ctx.Translator.Get("support_manager_connected"))
	_, _ = h.bot.Send(msg)

	// Notify admin
	msg2 := tgbotapi.NewMessage(adminID, fmt.Sprintf("Ви підключилися до чату з користувачем %d", userID))
	_, _ = h.bot.Send(msg2)

	// Notify other admins
	h.notifyOtherAdmins(ctx, adminID, userID)

	return nil
}

// HandleAdminMessage handles admin messages in support chat
func (h *SupportChatHandler) HandleAdminMessage(ctx *SupportContext) error {
	adminID := ctx.Update.Message.From.ID

	// Find user connected to this admin
	ActiveChats.RLock()
	var userID int64 = -1
	for uid, aid := range ActiveChats.chats {
		if aid == adminID {
			userID = uid
			break
		}
	}
	ActiveChats.RUnlock()

	if userID == -1 {
		msg := tgbotapi.NewMessage(adminID, ctx.Translator.Get("support_no_active_chat"))
		_, _ = h.bot.Send(msg)
		return nil
	}

	// Send message to user
	text := fmt.Sprintf("Відповідь від менеджера: %s", ctx.Update.Message.Text)
	msg := tgbotapi.NewMessage(userID, text)
	_, err := h.bot.Send(msg)

	// Log message
	h.logSupportMessage(userID, adminID, ctx.Update.Message.Text, false)

	return err
}

// HandleEndChat handles /end_chat command for admins
func (h *SupportChatHandler) HandleEndChat(ctx *SupportContext) error {
	adminID := ctx.Update.Message.From.ID

	// Find and remove chat
	ActiveChats.Lock()
	var userID int64 = -1
	for uid, aid := range ActiveChats.chats {
		if aid == adminID {
			userID = uid
			delete(ActiveChats.chats, uid)
			break
		}
	}
	ActiveChats.Unlock()

	if userID == -1 {
		msg := tgbotapi.NewMessage(adminID, ctx.Translator.Get("support_no_active_chat"))
		_, _ = h.bot.Send(msg)
		return nil
	}

	// Clear user state
	h.stateManager.ClearState(userID)

	// Notify user
	msg := tgbotapi.NewMessage(userID, ctx.Translator.Get("support_chat_ended"))
	_, _ = h.bot.Send(msg)

	// Notify admin
	msg2 := tgbotapi.NewMessage(adminID, fmt.Sprintf("Чат з користувачем %d завершено.", userID))
	_, _ = h.bot.Send(msg2)

	return nil
}

// notifyAdmins notifies admins about new support request
func (h *SupportChatHandler) notifyAdmins(ctx *SupportContext, userID int64) {
	fullName := ctx.Update.Message.From.FirstName
	if ctx.Update.Message.From.LastName != "" {
		fullName += " " + ctx.Update.Message.From.LastName
	}

	message := fmt.Sprintf("Новий запит на чат підтримки від користувача %s (ID: %d)", fullName, userID)

	for _, adminID := range h.config.AdminTelegramIDs {
		msg := tgbotapi.NewMessage(adminID, message)
		_, _ = h.bot.Send(msg)
	}
}

// notifyOtherAdmins notifies other admins about connection
func (h *SupportChatHandler) notifyOtherAdmins(ctx *SupportContext, connectedAdminID, userID int64) {
	message := fmt.Sprintf("Адміністратор %s підключився до чату з користувачем %d",
		ctx.Update.Message.From.FirstName, userID)

	for _, adminID := range h.config.AdminTelegramIDs {
		if adminID != connectedAdminID {
			msg := tgbotapi.NewMessage(adminID, message)
			_, _ = h.bot.Send(msg)
		}
	}
}

// logSupportMessage logs support chat message
func (h *SupportChatHandler) logSupportMessage(userID, adminID int64, text string, fromUser bool) {
	if fromUser {
		// Log as incoming from user
		log := &models.MessageLog{
			TelegramID:  userID,
			Direction:   models.DirectionIncoming,
			MessageText: &text,
		}
		_ = h.logRepo.LogMessage(context.Background(), log)
	} else {
		// Log as outgoing to user
		log := &models.MessageLog{
			TelegramID:  userID,
			Direction:   models.DirectionOutgoing,
			MessageText: &text,
		}
		_ = h.logRepo.LogMessage(context.Background(), log)
	}
}

