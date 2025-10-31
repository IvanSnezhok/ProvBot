package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"provbot/internal/database"
	"provbot/internal/models"
)

type OutageRepository struct{}

func NewOutageRepository() *OutageRepository {
	return &OutageRepository{}
}

// Create creates a new outage
func (r *OutageRepository) Create(ctx context.Context, outage *models.Outage) error {
	query := `INSERT INTO outages (location, description, status, created_by)
	          VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	
	err := database.PostgresDB.QueryRow(ctx, query,
		outage.Location, outage.Description, outage.Status, outage.CreatedBy,
	).Scan(&outage.ID, &outage.CreatedAt, &outage.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create outage: %w", err)
	}
	return nil
}

// GetByID retrieves an outage by ID
func (r *OutageRepository) GetByID(ctx context.Context, id int) (*models.Outage, error) {
	query := `SELECT id, location, description, status, created_by, created_at, updated_at, resolved_at
	          FROM outages WHERE id = $1`
	
	var outage models.Outage
	err := database.PostgresDB.QueryRow(ctx, query, id).Scan(
		&outage.ID, &outage.Location, &outage.Description, &outage.Status,
		&outage.CreatedBy, &outage.CreatedAt, &outage.UpdatedAt, &outage.ResolvedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get outage: %w", err)
	}
	return &outage, nil
}

// GetByLocation retrieves active outages for a location
func (r *OutageRepository) GetByLocation(ctx context.Context, location string) ([]models.Outage, error) {
	query := `SELECT id, location, description, status, created_by, created_at, updated_at, resolved_at
	          FROM outages WHERE location = $1 AND status = 'active' ORDER BY created_at DESC`
	
	rows, err := database.PostgresDB.Query(ctx, query, location)
	if err != nil {
		return nil, fmt.Errorf("failed to get outages by location: %w", err)
	}
	defer rows.Close()

	var outages []models.Outage
	for rows.Next() {
		var outage models.Outage
		err := rows.Scan(
			&outage.ID, &outage.Location, &outage.Description, &outage.Status,
			&outage.CreatedBy, &outage.CreatedAt, &outage.UpdatedAt, &outage.ResolvedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan outage: %w", err)
		}
		outages = append(outages, outage)
	}
	return outages, nil
}

// GetAllActive retrieves all active outages
func (r *OutageRepository) GetAllActive(ctx context.Context) ([]models.Outage, error) {
	query := `SELECT id, location, description, status, created_by, created_at, updated_at, resolved_at
	          FROM outages WHERE status = 'active' ORDER BY created_at DESC`
	
	rows, err := database.PostgresDB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active outages: %w", err)
	}
	defer rows.Close()

	var outages []models.Outage
	for rows.Next() {
		var outage models.Outage
		err := rows.Scan(
			&outage.ID, &outage.Location, &outage.Description, &outage.Status,
			&outage.CreatedBy, &outage.CreatedAt, &outage.UpdatedAt, &outage.ResolvedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan outage: %w", err)
		}
		outages = append(outages, outage)
	}
	return outages, nil
}

// Update updates an outage
func (r *OutageRepository) Update(ctx context.Context, outage *models.Outage) error {
	var resolvedAt interface{}
	if outage.ResolvedAt != nil {
		resolvedAt = *outage.ResolvedAt
	}
	
	query := `UPDATE outages SET location = $1, description = $2, status = $3, 
	          updated_at = $4, resolved_at = $5 WHERE id = $6`
	
	outage.UpdatedAt = time.Now()
	_, err := database.PostgresDB.Exec(ctx, query,
		outage.Location, outage.Description, outage.Status,
		outage.UpdatedAt, resolvedAt, outage.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update outage: %w", err)
	}
	return nil
}

// Resolve resolves an outage
func (r *OutageRepository) Resolve(ctx context.Context, id int) error {
	now := time.Now()
	query := `UPDATE outages SET status = 'resolved', updated_at = $1, resolved_at = $1 WHERE id = $2`
	
	_, err := database.PostgresDB.Exec(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("failed to resolve outage: %w", err)
	}
	return nil
}

