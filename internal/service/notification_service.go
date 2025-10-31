package service

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"provbot/internal/models"
	"provbot/internal/repository"
)

type NotificationService struct {
	bot       *tgbotapi.BotAPI
	logRepo   *repository.LogRepository
	userRepo  *repository.UserRepository
}

func NewNotificationService(bot *tgbotapi.BotAPI, logRepo *repository.LogRepository, userRepo *repository.UserRepository) *NotificationService {
	return &NotificationService{
		bot:      bot,
		logRepo:  logRepo,
		userRepo: userRepo,
	}
}

// SendMessage sends a message to a user and logs it
func (s *NotificationService) SendMessage(ctx context.Context, chatID int64, text string) (int, error) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	
	sentMsg, err := s.bot.Send(msg)
	if err != nil {
		return 0, fmt.Errorf("failed to send message: %w", err)
	}

	// Log outgoing message
	user, _ := s.userRepo.GetByTelegramID(ctx, chatID)
	var userID *int
	if user != nil {
		userID = &user.ID
	}
	
	messageID := int64(sentMsg.MessageID)
	log := &models.MessageLog{
		UserID:     userID,
		TelegramID: chatID,
		Direction:  models.DirectionOutgoing,
		MessageText: &text,
		MessageID:  &messageID,
	}
	
	if err := s.logRepo.LogMessage(ctx, log); err != nil {
		// Log error but don't fail the send
		fmt.Printf("Failed to log message: %v\n", err)
	}

	return sentMsg.MessageID, nil
}

// BroadcastMessage sends a message to all active users
func (s *NotificationService) BroadcastMessage(ctx context.Context, message string) (int, int, error) {
	users, err := s.userRepo.GetAllActiveUsers(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get users: %w", err)
	}

	successCount := 0
	totalCount := len(users)

	for _, user := range users {
		if _, err := s.SendMessage(ctx, user.TelegramID, message); err == nil {
			successCount++
		}
	}

	return successCount, totalCount, nil
}

// SendOutageNotification sends outage notification to users in affected location
func (s *NotificationService) SendOutageNotification(ctx context.Context, outage *models.Outage, affectedUserIDs []int64) error {
	message := fmt.Sprintf("⚠️ Аварія\n\nЛокація: %s\nОпис: %s", outage.Location, outage.Description)
	
	for _, userID := range affectedUserIDs {
		if _, err := s.SendMessage(ctx, userID, message); err != nil {
			// Continue with other users even if one fails
			fmt.Printf("Failed to send outage notification to %d: %v\n", userID, err)
		}
	}
	
	return nil
}

