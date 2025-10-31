package service

import (
	"context"
	"fmt"
	"time"

	"provbot/internal/models"
	"provbot/internal/repository"
)

type SupportService struct {
	logRepo *repository.LogRepository
}

func NewSupportService(logRepo *repository.LogRepository) *SupportService {
	return &SupportService{
		logRepo: logRepo,
	}
}

// CreateTicket creates a support ticket (logs as message)
func (s *SupportService) CreateTicket(ctx context.Context, userID *int, telegramID int64, message string) error {
	// Log support message as incoming
	log := &models.MessageLog{
		UserID:      userID,
		TelegramID:  telegramID,
		Direction:   models.DirectionIncoming,
		MessageText: &message,
	}
	
	if err := s.logRepo.LogMessage(ctx, log); err != nil {
		return fmt.Errorf("failed to log support ticket: %w", err)
	}
	
	// Here you could add additional logic like:
	// - Create ticket in ticketing system
	// - Notify admins
	// - Store in separate tickets table
	
	return nil
}

// GetTicketHistory retrieves user's support message history
func (s *SupportService) GetTicketHistory(ctx context.Context, telegramID int64, limit int) ([]models.MessageLog, error) {
	logs, err := s.logRepo.GetMessageLogs(ctx, telegramID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket history: %w", err)
	}
	return logs, nil
}

// SupportTicket represents a support ticket
type SupportTicket struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

