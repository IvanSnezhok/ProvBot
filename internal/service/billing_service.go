package service

import (
	"context"
	"fmt"

	"provbot/internal/models"
	"provbot/internal/repository"
)

type BillingService struct {
	billingRepo *repository.BillingRepository
}

func NewBillingService(billingRepo *repository.BillingRepository) *BillingService {
	return &BillingService{
		billingRepo: billingRepo,
	}
}

// GetBalance retrieves user balance from billing system
func (s *BillingService) GetBalance(ctx context.Context, userID int64) (float64, error) {
	user, err := s.billingRepo.GetUserByID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get billing user: %w", err)
	}
	if user == nil {
		return 0, fmt.Errorf("user not found in billing system")
	}
	return user.Balance, nil
}

// TopUpBalance adds amount to user balance
func (s *BillingService) TopUpBalance(ctx context.Context, userID int64, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	return s.billingRepo.UpdateBalance(ctx, userID, amount)
}

// GetServices retrieves user services
func (s *BillingService) GetServices(ctx context.Context, userID int64) ([]models.BillingService, error) {
	services, err := s.billingRepo.GetServicesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}
	return services, nil
}

// UpdateServiceStatus updates service status in billing
func (s *BillingService) UpdateServiceStatus(ctx context.Context, serviceID int64, status string) error {
	return s.billingRepo.UpdateServiceStatus(ctx, serviceID, status)
}

// UpdateUserStatus updates user status in billing
func (s *BillingService) UpdateUserStatus(ctx context.Context, userID int64, status string) error {
	return s.billingRepo.UpdateUserStatus(ctx, userID, status)
}

// GetBillingUser retrieves billing user by ID
func (s *BillingService) GetBillingUser(ctx context.Context, userID int64) (*models.BillingUser, error) {
	user, err := s.billingRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing user: %w", err)
	}
	return user, nil
}

