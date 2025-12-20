package service

import (
	"context"
	"fmt"

	"provbot/internal/models"
	"provbot/internal/repository"
)

type OutageService struct {
	outageRepo *repository.OutageRepository
}

func NewOutageService(outageRepo *repository.OutageRepository) *OutageService {
	return &OutageService{
		outageRepo: outageRepo,
	}
}

// CreateOutage creates a new outage
func (s *OutageService) CreateOutage(ctx context.Context, location, description string, groupID, street *string, createdBy *int) (*models.Outage, error) {
	outage := &models.Outage{
		Location:    location,
		Description: description,
		Status:      models.OutageStatusActive,
		GroupID:     groupID,
		Street:      street,
		CreatedBy:   createdBy,
	}

	if err := s.outageRepo.Create(ctx, outage); err != nil {
		return nil, fmt.Errorf("failed to create outage: %w", err)
	}

	return outage, nil
}

// GetOutageByID retrieves an outage by ID
func (s *OutageService) GetOutageByID(ctx context.Context, id int) (*models.Outage, error) {
	return s.outageRepo.GetByID(ctx, id)
}

// GetActiveOutages retrieves all active outages
func (s *OutageService) GetActiveOutages(ctx context.Context) ([]models.Outage, error) {
	return s.outageRepo.GetAllActive(ctx)
}

// GetOutagesByGroupID retrieves active outages for a specific billing group
func (s *OutageService) GetOutagesByGroupID(ctx context.Context, groupID string) ([]models.Outage, error) {
	return s.outageRepo.GetByGroupID(ctx, groupID)
}

// GetOutagesByContract retrieves active outages for a specific contract
func (s *OutageService) GetOutagesByContract(ctx context.Context, contract string) ([]models.Outage, error) {
	return s.outageRepo.GetByContract(ctx, contract)
}

// HasActiveOutageForUser checks if there's an active outage affecting the user
func (s *OutageService) HasActiveOutageForUser(ctx context.Context, groupID, contract string) (bool, error) {
	return s.outageRepo.HasActiveOutageForUser(ctx, groupID, contract)
}

// GetOutageMessageForUser gets the outage message for a user if there's an active outage
func (s *OutageService) GetOutageMessageForUser(ctx context.Context, groupID, contract string) (string, error) {
	return s.outageRepo.GetOutageMessageForUser(ctx, groupID, contract)
}

// ResolveOutage resolves an outage by ID
func (s *OutageService) ResolveOutage(ctx context.Context, id int) error {
	return s.outageRepo.Resolve(ctx, id)
}

// UpdateOutage updates an outage
func (s *OutageService) UpdateOutage(ctx context.Context, outage *models.Outage) error {
	return s.outageRepo.Update(ctx, outage)
}
