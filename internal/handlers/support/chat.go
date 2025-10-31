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

// HandleAdminConnect handles /connect command or callback for admins
func (h *SupportChatHandler) HandleAdminConnect(ctx *SupportContext, userID int64) error {
	var adminID int64
	var chatID int64
	var messageID int
	
	// Get admin ID from Message or CallbackQuery
	if ctx.Update.Message != nil {
		adminID = ctx.Update.Message.From.ID
		chatID = ctx.Update.Message.Chat.ID
	} else if ctx.Update.CallbackQuery != nil {
		adminID = ctx.Update.CallbackQuery.From.ID
		chatID = ctx.Update.CallbackQuery.Message.Chat.ID
		messageID = ctx.Update.CallbackQuery.Message.MessageID
		
		// Answer callback query
		callbackConfig := tgbotapi.NewCallback(ctx.Update.CallbackQuery.ID, "")
		_, _ = h.bot.Request(callbackConfig)
	} else {
		return fmt.Errorf("cannot determine admin ID")
	}

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

	// Notify user with end chat button
	userKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(ctx.Translator.Get("support_end_chat_button")),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(ctx.Translator.Get("menu_back")),
		),
	)
	userKeyboard.ResizeKeyboard = true
	
	msg := tgbotapi.NewMessage(userID, ctx.Translator.Get("support_manager_connected"))
	msg.ReplyMarkup = userKeyboard
	_, _ = h.bot.Send(msg)

	// Notify admin with end chat button (inline)
	adminKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				ctx.Translator.Get("support_end_chat_button"),
				"support_end_chat",
			),
		),
	)
	
	adminMessage := ctx.Translator.Get("support_admin_connected")
	
	// If it's a callback, edit the message; otherwise send new one
	if ctx.Update.CallbackQuery != nil && messageID > 0 {
		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, adminMessage)
		editMsg.ReplyMarkup = &adminKeyboard
		_, _ = h.bot.Send(editMsg)
	} else {
		msg2 := tgbotapi.NewMessage(adminID, adminMessage)
		msg2.ReplyMarkup = adminKeyboard
		_, _ = h.bot.Send(msg2)
	}

	// Notify other admins
	h.notifyOtherAdmins(ctx, adminID, userID)

	return nil
}

// IsAdminInChat checks if admin is connected to an active support chat
func (h *SupportChatHandler) IsAdminInChat(adminID int64) bool {
	ActiveChats.RLock()
	defer ActiveChats.RUnlock()
	
	for _, aid := range ActiveChats.chats {
		if aid == adminID {
			return true
		}
	}
	return false
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

// HandleEndChat handles /end_chat command or callback for admins and users
func (h *SupportChatHandler) HandleEndChat(ctx *SupportContext) error {
	var chatID int64
	var messageID int
	var isAdmin bool
	var requesterID int64
	
	// Get requester ID from Message or CallbackQuery
	if ctx.Update.Message != nil {
		requesterID = ctx.Update.Message.From.ID
		chatID = ctx.Update.Message.Chat.ID
		isAdmin = ctx.Config.IsAdmin(requesterID)
	} else if ctx.Update.CallbackQuery != nil {
		requesterID = ctx.Update.CallbackQuery.From.ID
		chatID = ctx.Update.CallbackQuery.Message.Chat.ID
		messageID = ctx.Update.CallbackQuery.Message.MessageID
		isAdmin = ctx.Config.IsAdmin(requesterID)
		
		// Answer callback query
		callbackConfig := tgbotapi.NewCallback(ctx.Update.CallbackQuery.ID, "")
		_, _ = h.bot.Request(callbackConfig)
	} else {
		return fmt.Errorf("cannot determine requester ID")
	}

	// Find and remove chat
	ActiveChats.Lock()
	var userID int64 = -1
	var adminID int64 = -1
	
	if isAdmin {
		// Admin ending chat - find user by admin ID
		for uid, aid := range ActiveChats.chats {
			if aid == requesterID {
				userID = uid
				adminID = requesterID
				delete(ActiveChats.chats, uid)
				break
			}
		}
	} else {
		// User ending chat - find admin by user ID
		if aid, exists := ActiveChats.chats[requesterID]; exists {
			userID = requesterID
			adminID = aid
			delete(ActiveChats.chats, requesterID)
		}
	}
	ActiveChats.Unlock()

	if userID == -1 {
		msg := tgbotapi.NewMessage(requesterID, ctx.Translator.Get("support_no_active_chat"))
		_, _ = h.bot.Send(msg)
		return nil
	}

	// Clear user state
	h.stateManager.ClearState(userID)

	// Remove keyboard for user
	removeKeyboard := tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
	
	// Notify user
	msg := tgbotapi.NewMessage(userID, ctx.Translator.Get("support_chat_ended"))
	msg.ReplyMarkup = removeKeyboard
	_, _ = h.bot.Send(msg)

	// Notify admin
	endMessage := fmt.Sprintf("Чат з користувачем %d завершено.", userID)
	
	// If it's a callback, edit the message; otherwise send new one
	if ctx.Update.CallbackQuery != nil && messageID > 0 {
		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, endMessage)
		_, _ = h.bot.Send(editMsg)
	} else {
		msg2 := tgbotapi.NewMessage(adminID, endMessage)
		_, _ = h.bot.Send(msg2)
	}

	return nil
}

// notifyAdmins notifies admins about new support request
func (h *SupportChatHandler) notifyAdmins(ctx *SupportContext, userID int64) {
	fullName := ctx.Update.Message.From.FirstName
	if ctx.Update.Message.From.LastName != "" {
		fullName += " " + ctx.Update.Message.From.LastName
	}

	message := fmt.Sprintf("Новий запит на чат підтримки від користувача %s (ID: %d)", fullName, userID)

	// Create inline keyboard with connect button
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				ctx.Translator.Get("support_connect_button"),
				fmt.Sprintf("support_connect_%d", userID),
			),
		),
	)

	for _, adminID := range h.config.AdminTelegramIDs {
		msg := tgbotapi.NewMessage(adminID, message)
		msg.ReplyMarkup = keyboard
		_, _ = h.bot.Send(msg)
	}
}

// notifyOtherAdmins notifies other admins about connection
func (h *SupportChatHandler) notifyOtherAdmins(ctx *SupportContext, connectedAdminID, userID int64) {
	var adminFirstName string
	
	// Get admin first name from Message or CallbackQuery
	if ctx.Update.Message != nil {
		adminFirstName = ctx.Update.Message.From.FirstName
	} else if ctx.Update.CallbackQuery != nil {
		adminFirstName = ctx.Update.CallbackQuery.From.FirstName
	} else {
		adminFirstName = "Адміністратор" // Fallback if cannot determine
	}
	
	message := fmt.Sprintf("Адміністратор %s підключився до чату з користувачем %d",
		adminFirstName, userID)

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

