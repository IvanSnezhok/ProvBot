package service

import (
	"context"
	"fmt"

	"provbot/internal/i18n"
	"provbot/internal/models"
	"provbot/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// RegisterOrUpdateUser registers a new user or updates existing one
func (s *UserService) RegisterOrUpdateUser(ctx context.Context, telegramID int64, username, firstName, lastName *string) (*models.User, error) {
	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		// Create new user
		user = &models.User{
			TelegramID: telegramID,
			Username:   username,
			FirstName:  firstName,
			LastName:   lastName,
			Language:   i18n.DefaultLanguage,
			IsActive:   true,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		// Update existing user
		user.Username = username
		user.FirstName = firstName
		user.LastName = lastName
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	return user, nil
}

// GetUser retrieves user by Telegram ID
func (s *UserService) GetUser(ctx context.Context, telegramID int64) (*models.User, error) {
	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// UpdateLanguage updates user language
func (s *UserService) UpdateLanguage(ctx context.Context, telegramID int64, language string) error {
	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Validate language
	if language != i18n.LanguageUA && language != i18n.LanguageEN {
		language = i18n.DefaultLanguage
	}

	user.Language = language
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update language: %w", err)
	}
	return nil
}

// IsAdmin checks if user is admin
func (s *UserService) IsAdmin(ctx context.Context, telegramID int64) (bool, error) {
	return s.userRepo.IsAdmin(ctx, telegramID)
}

// IsBanned checks if user is banned
func (s *UserService) IsBanned(ctx context.Context, telegramID int64) (bool, error) {
	return s.userRepo.IsBanned(ctx, telegramID)
}

// SetBan sets ban status for a user
func (s *UserService) SetBan(ctx context.Context, telegramID int64, banned bool) error {
	return s.userRepo.SetBan(ctx, telegramID, banned)
}

// GetBannedUsers retrieves all banned users
func (s *UserService) GetBannedUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.GetBannedUsers(ctx)
}

