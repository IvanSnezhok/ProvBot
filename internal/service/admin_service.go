package service

import (
	"context"
	"fmt"

	"provbot/internal/models"
	"provbot/internal/repository"
)

type AdminService struct {
	userRepo    *repository.UserRepository
	outageRepo  *repository.OutageRepository
	billingRepo *repository.BillingRepository
}

func NewAdminService(userRepo *repository.UserRepository, outageRepo *repository.OutageRepository, billingRepo *repository.BillingRepository) *AdminService {
	return &AdminService{
		userRepo:    userRepo,
		outageRepo:  outageRepo,
		billingRepo: billingRepo,
	}
}

// CreateOutage creates a new outage
func (s *AdminService) CreateOutage(ctx context.Context, location, description string, createdBy int) (*models.Outage, error) {
	outage := &models.Outage{
		Location:    location,
		Description: description,
		Status:      models.OutageStatusActive,
		CreatedBy:   &createdBy,
	}

	if err := s.outageRepo.Create(ctx, outage); err != nil {
		return nil, fmt.Errorf("failed to create outage: %w", err)
	}
	return outage, nil
}

// ResolveOutage resolves an outage
func (s *AdminService) ResolveOutage(ctx context.Context, outageID int) error {
	return s.outageRepo.Resolve(ctx, outageID)
}

// GetOutages retrieves all active outages
func (s *AdminService) GetOutages(ctx context.Context) ([]models.Outage, error) {
	return s.outageRepo.GetAllActive(ctx)
}

// GetOutagesByLocation retrieves outages for a specific location
func (s *AdminService) GetOutagesByLocation(ctx context.Context, location string) ([]models.Outage, error) {
	return s.outageRepo.GetByLocation(ctx, location)
}

// UpdateUserInBilling updates user status in billing system
func (s *AdminService) UpdateUserInBilling(ctx context.Context, userID int64, status string) error {
	// We need contract to update status
	user, err := s.billingRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}
	return s.billingRepo.SetUserStatus(ctx, user.Contract, status)
}

// UpdateServiceInBilling updates service status in billing system
func (s *AdminService) UpdateServiceInBilling(ctx context.Context, serviceID int64, status string) error {
	// return s.billingRepo.UpdateServiceStatus(ctx, serviceID, status)
	return fmt.Errorf("not implemented")
}

// GetAllUsers retrieves all active users for broadcast
func (s *AdminService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.GetAllActiveUsers(ctx)
}
