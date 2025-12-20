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
	query := `INSERT INTO outages (location, description, status, group_id, street, created_by)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`

	err := database.PostgresDB.QueryRow(ctx, query,
		outage.Location, outage.Description, outage.Status, outage.GroupID, outage.Street, outage.CreatedBy,
	).Scan(&outage.ID, &outage.CreatedAt, &outage.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create outage: %w", err)
	}
	return nil
}

// GetByID retrieves an outage by ID
func (r *OutageRepository) GetByID(ctx context.Context, id int) (*models.Outage, error) {
	query := `SELECT id, location, description, status, group_id, street, created_by, created_at, updated_at, resolved_at
	          FROM outages WHERE id = $1`

	var outage models.Outage
	err := database.PostgresDB.QueryRow(ctx, query, id).Scan(
		&outage.ID, &outage.Location, &outage.Description, &outage.Status,
		&outage.GroupID, &outage.Street, &outage.CreatedBy, &outage.CreatedAt, &outage.UpdatedAt, &outage.ResolvedAt,
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
	query := `SELECT id, location, description, status, group_id, street, created_by, created_at, updated_at, resolved_at
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
			&outage.GroupID, &outage.Street, &outage.CreatedBy, &outage.CreatedAt, &outage.UpdatedAt, &outage.ResolvedAt,
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
	query := `SELECT id, location, description, status, group_id, street, created_by, created_at, updated_at, resolved_at
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
			&outage.GroupID, &outage.Street, &outage.CreatedBy, &outage.CreatedAt, &outage.UpdatedAt, &outage.ResolvedAt,
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
	          group_id = $4, street = $5, updated_at = $6, resolved_at = $7 WHERE id = $8`

	outage.UpdatedAt = time.Now()
	_, err := database.PostgresDB.Exec(ctx, query,
		outage.Location, outage.Description, outage.Status,
		outage.GroupID, outage.Street, outage.UpdatedAt, resolvedAt, outage.ID,
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

// GetByGroupID retrieves active outages for a specific billing group
func (r *OutageRepository) GetByGroupID(ctx context.Context, groupID string) ([]models.Outage, error) {
	query := `SELECT id, location, description, status, group_id, street, created_by, created_at, updated_at, resolved_at
	          FROM outages WHERE (group_id = $1 OR group_id IS NULL) AND status = 'active' ORDER BY created_at DESC`

	rows, err := database.PostgresDB.Query(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get outages by group: %w", err)
	}
	defer rows.Close()

	var outages []models.Outage
	for rows.Next() {
		var outage models.Outage
		err := rows.Scan(
			&outage.ID, &outage.Location, &outage.Description, &outage.Status,
			&outage.GroupID, &outage.Street, &outage.CreatedBy, &outage.CreatedAt, &outage.UpdatedAt, &outage.ResolvedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan outage: %w", err)
		}
		outages = append(outages, outage)
	}
	return outages, nil
}

// GetByContract retrieves active outages for a specific contract (uses contract as group identifier)
func (r *OutageRepository) GetByContract(ctx context.Context, contract string) ([]models.Outage, error) {
	// First try to get by contract directly, then fall back to general outages
	query := `SELECT id, location, description, status, group_id, street, created_by, created_at, updated_at, resolved_at
	          FROM outages WHERE (group_id = $1 OR group_id IS NULL) AND status = 'active' ORDER BY created_at DESC`

	rows, err := database.PostgresDB.Query(ctx, query, contract)
	if err != nil {
		return nil, fmt.Errorf("failed to get outages by contract: %w", err)
	}
	defer rows.Close()

	var outages []models.Outage
	for rows.Next() {
		var outage models.Outage
		err := rows.Scan(
			&outage.ID, &outage.Location, &outage.Description, &outage.Status,
			&outage.GroupID, &outage.Street, &outage.CreatedBy, &outage.CreatedAt, &outage.UpdatedAt, &outage.ResolvedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan outage: %w", err)
		}
		outages = append(outages, outage)
	}
	return outages, nil
}

// HasActiveOutageForUser checks if there's an active outage affecting the user
func (r *OutageRepository) HasActiveOutageForUser(ctx context.Context, groupID, contract string) (bool, error) {
	query := `SELECT COUNT(*) FROM outages
	          WHERE status = 'active' AND (group_id = $1 OR group_id = $2 OR group_id IS NULL)`

	var count int
	err := database.PostgresDB.QueryRow(ctx, query, groupID, contract).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check outage status: %w", err)
	}
	return count > 0, nil
}

// GetOutageMessageForUser gets the first active outage message for a user
func (r *OutageRepository) GetOutageMessageForUser(ctx context.Context, groupID, contract string) (string, error) {
	query := `SELECT description FROM outages
	          WHERE status = 'active' AND (group_id = $1 OR group_id = $2)
	          ORDER BY created_at DESC LIMIT 1`

	var message string
	err := database.PostgresDB.QueryRow(ctx, query, groupID, contract).Scan(&message)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get outage message: %w", err)
	}
	return message, nil
}

