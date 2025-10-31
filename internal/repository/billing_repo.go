package repository

import (
	"context"
	"database/sql"
	"fmt"

	"provbot/internal/database"
	"provbot/internal/models"
)

type BillingRepository struct{}

func NewBillingRepository() *BillingRepository {
	return &BillingRepository{}
}

// GetUserByTelegramID retrieves billing user by Telegram ID
// Note: This assumes there's a mapping table or field in billing DB
func (r *BillingRepository) GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.BillingUser, error) {
	// This is a placeholder - actual query depends on billing DB structure
	query := `SELECT id, username, balance, status, service_id, created_at, updated_at 
	          FROM users WHERE telegram_id = ? LIMIT 1`
	
	var user models.BillingUser
	err := database.MySQLDB.QueryRowContext(ctx, query, telegramID).Scan(
		&user.ID, &user.Username, &user.Balance, &user.Status,
		&user.ServiceID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get billing user: %w", err)
	}
	return &user, nil
}

// GetUserByID retrieves billing user by ID
func (r *BillingRepository) GetUserByID(ctx context.Context, userID int64) (*models.BillingUser, error) {
	query := `SELECT id, username, balance, status, service_id, created_at, updated_at 
	          FROM users WHERE id = ? LIMIT 1`
	
	var user models.BillingUser
	err := database.MySQLDB.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Username, &user.Balance, &user.Status,
		&user.ServiceID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get billing user: %w", err)
	}
	return &user, nil
}

// UpdateBalance updates user balance
func (r *BillingRepository) UpdateBalance(ctx context.Context, userID int64, amount float64) error {
	query := `UPDATE users SET balance = balance + ?, updated_at = NOW() WHERE id = ?`
	
	result, err := database.MySQLDB.ExecContext(ctx, query, amount, userID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// GetServicesByUserID retrieves services for a user
func (r *BillingRepository) GetServicesByUserID(ctx context.Context, userID int64) ([]models.BillingService, error) {
	query := `SELECT id, user_id, name, status, location, ip_address, created_at, updated_at 
	          FROM services WHERE user_id = ?`
	
	rows, err := database.MySQLDB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}
	defer rows.Close()

	var services []models.BillingService
	for rows.Next() {
		var service models.BillingService
		err := rows.Scan(
			&service.ID, &service.UserID, &service.Name, &service.Status,
			&service.Location, &service.IPAddress, &service.CreatedAt, &service.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service: %w", err)
		}
		services = append(services, service)
	}
	return services, nil
}

// UpdateServiceStatus updates service status
func (r *BillingRepository) UpdateServiceStatus(ctx context.Context, serviceID int64, status string) error {
	query := `UPDATE services SET status = ?, updated_at = NOW() WHERE id = ?`
	
	result, err := database.MySQLDB.ExecContext(ctx, query, status, serviceID)
	if err != nil {
		return fmt.Errorf("failed to update service status: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("service not found")
	}
	return nil
}

// UpdateUserStatus updates user status in billing
func (r *BillingRepository) UpdateUserStatus(ctx context.Context, userID int64, status string) error {
	query := `UPDATE users SET status = ?, updated_at = NOW() WHERE id = ?`
	
	result, err := database.MySQLDB.ExecContext(ctx, query, status, userID)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// SearchByPhone searches user in billing by phone number
// Returns user data, contract number, and error
func (r *BillingRepository) SearchByPhone(ctx context.Context, phone string) (*models.BillingUser, string, error) {
	// This is a placeholder - actual query depends on billing DB structure
	// Assuming there's a phone field in users table
	query := `SELECT id, username, balance, status, service_id, created_at, updated_at, contract 
	          FROM users WHERE phone = ? LIMIT 1`
	
	var user models.BillingUser
	var contract string
	err := database.MySQLDB.QueryRowContext(ctx, query, phone).Scan(
		&user.ID, &user.Username, &user.Balance, &user.Status,
		&user.ServiceID, &user.CreatedAt, &user.UpdatedAt, &contract,
	)
	if err == sql.ErrNoRows {
		return nil, "", nil
	}
	if err != nil {
		return nil, "", fmt.Errorf("failed to search by phone: %w", err)
	}
	return &user, contract, nil
}

// SearchByContract searches user in billing by contract number
func (r *BillingRepository) SearchByContract(ctx context.Context, contract string) (*models.BillingUser, error) {
	query := `SELECT id, username, balance, status, service_id, created_at, updated_at 
	          FROM users WHERE contract = ? LIMIT 1`
	
	var user models.BillingUser
	err := database.MySQLDB.QueryRowContext(ctx, query, contract).Scan(
		&user.ID, &user.Username, &user.Balance, &user.Status,
		&user.ServiceID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to search by contract: %w", err)
	}
	return &user, nil
}

// SearchByName searches users in billing by name
func (r *BillingRepository) SearchByName(ctx context.Context, name string) ([]models.BillingUser, error) {
	query := `SELECT id, username, balance, status, service_id, created_at, updated_at 
	          FROM users WHERE fio LIKE ? LIMIT 10`
	
	searchPattern := "%" + name + "%"
	rows, err := database.MySQLDB.QueryContext(ctx, query, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search by name: %w", err)
	}
	defer rows.Close()

	var users []models.BillingUser
	for rows.Next() {
		var user models.BillingUser
		err := rows.Scan(
			&user.ID, &user.Username, &user.Balance, &user.Status,
			&user.ServiceID, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

// SearchByAddress searches users in billing by address
func (r *BillingRepository) SearchByAddress(ctx context.Context, address string) ([]models.BillingUser, error) {
	query := `SELECT id, username, balance, status, service_id, created_at, updated_at 
	          FROM users WHERE address LIKE ? LIMIT 10`
	
	searchPattern := "%" + address + "%"
	rows, err := database.MySQLDB.QueryContext(ctx, query, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search by address: %w", err)
	}
	defer rows.Close()

	var users []models.BillingUser
	for rows.Next() {
		var user models.BillingUser
		err := rows.Scan(
			&user.ID, &user.Username, &user.Balance, &user.Status,
			&user.ServiceID, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

// CheckContractExists checks if contract exists in billing
func (r *BillingRepository) CheckContractExists(ctx context.Context, contract string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE contract = ?`
	var count int
	err := database.MySQLDB.QueryRowContext(ctx, query, contract).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check contract: %w", err)
	}
	return count > 0, nil
}

// GetTariffByContract gets tariff name by contract
func (r *BillingRepository) GetTariffByContract(ctx context.Context, contract string) (string, error) {
	query := `SELECT tariff FROM users WHERE contract = ? LIMIT 1`
	var tariff string
	err := database.MySQLDB.QueryRowContext(ctx, query, contract).Scan(&tariff)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get tariff: %w", err)
	}
	return tariff, nil
}

// EnableTemporaryPayment enables temporary payment for user
func (r *BillingRepository) EnableTemporaryPayment(ctx context.Context, contract string) (bool, float64, error) {
	// This is a placeholder - actual implementation depends on billing system
	// Should enable temporary payment and return success status and amount
	query := `UPDATE users SET temp_payment = 1, temp_payment_date = NOW() WHERE contract = ?`
	result, err := database.MySQLDB.ExecContext(ctx, query, contract)
	if err != nil {
		return false, 0, fmt.Errorf("failed to enable temporary payment: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, 0, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return false, 0, nil
	}
	
	// Get temporary payment amount (placeholder - should come from billing logic)
	amount := 0.0
	return true, amount, nil
}

// GetBalanceByContract retrieves user balance by contract number
func (r *BillingRepository) GetBalanceByContract(ctx context.Context, contract string) (float64, error) {
	query := `SELECT balance FROM users WHERE contract = ? LIMIT 1`
	var balance float64
	err := database.MySQLDB.QueryRowContext(ctx, query, contract).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get balance by contract: %w", err)
	}
	return balance, nil
}

