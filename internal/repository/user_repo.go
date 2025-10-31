package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"provbot/internal/database"
	"provbot/internal/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// GetByTelegramID retrieves a user by Telegram ID
func (r *UserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	query := `SELECT id, telegram_id, username, first_name, last_name, phone_number, contract, language, is_active, created_at, updated_at 
	          FROM users WHERE telegram_id = $1`
	
	var user models.User
	err := database.PostgresDB.QueryRow(ctx, query, telegramID).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName,
		&user.PhoneNumber, &user.Contract, &user.Language, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by telegram ID: %w", err)
	}
	return &user, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (telegram_id, username, first_name, last_name, phone_number, contract, language, is_active)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`
	
	err := database.PostgresDB.QueryRow(ctx, query,
		user.TelegramID, user.Username, user.FirstName, user.LastName,
		user.PhoneNumber, user.Contract, user.Language, user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET username = $1, first_name = $2, last_name = $3, 
	          phone_number = $4, contract = $5, language = $6, is_active = $7, updated_at = $8 
	          WHERE id = $9`
	
	user.UpdatedAt = time.Now()
	_, err := database.PostgresDB.Exec(ctx, query,
		user.Username, user.FirstName, user.LastName, user.PhoneNumber, user.Contract,
		user.Language, user.IsActive, user.UpdatedAt, user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// IsAdmin checks if user is admin
func (r *UserRepository) IsAdmin(ctx context.Context, telegramID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM admin_users WHERE telegram_id = $1`
	var count int
	err := database.PostgresDB.QueryRow(ctx, query, telegramID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check admin status: %w", err)
	}
	return count > 0, nil
}

// GetAllActiveUsers retrieves all active users
func (r *UserRepository) GetAllActiveUsers(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, telegram_id, username, first_name, last_name, phone_number, contract, language, is_active, created_at, updated_at 
	          FROM users WHERE is_active = true`
	
	rows, err := database.PostgresDB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName,
			&user.PhoneNumber, &user.Contract, &user.Language, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

// GetUsersWithContracts retrieves all active users with contracts
func (r *UserRepository) GetUsersWithContracts(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, telegram_id, username, first_name, last_name, phone_number, contract, language, is_active, created_at, updated_at 
	          FROM users WHERE is_active = true AND contract IS NOT NULL AND contract != ''`
	
	rows, err := database.PostgresDB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users with contracts: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName,
			&user.PhoneNumber, &user.Contract, &user.Language, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

