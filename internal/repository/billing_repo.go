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

// GetUserByContract retrieves billing user by contract number
func (r *BillingRepository) GetUserByContract(ctx context.Context, contract string) (*models.BillingUser, error) {
	query := `SELECT id, contract, fio, telefon, balance, paket, state, t_pay, start_day, srvs, address, grp 
	          FROM users WHERE contract = ? LIMIT 1`

	var user models.BillingUser
	err := database.MySQLDB.QueryRowContext(ctx, query, contract).Scan(
		&user.ID, &user.Contract, &user.Name, &user.Phone, &user.Balance,
		&user.PlanID, &user.Status, &user.TimePay, &user.StartDay, &user.Services,
		&user.Address, &user.Group,
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
func (r *BillingRepository) GetUserByID(ctx context.Context, id int64) (*models.BillingUser, error) {
	query := `SELECT id, contract, fio, telefon, balance, paket, state, t_pay, start_day, srvs, address, grp 
	          FROM users WHERE id = ? LIMIT 1`

	var user models.BillingUser
	err := database.MySQLDB.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Contract, &user.Name, &user.Phone, &user.Balance,
		&user.PlanID, &user.Status, &user.TimePay, &user.StartDay, &user.Services,
		&user.Address, &user.Group,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get billing user by id: %w", err)
	}
	return &user, nil
}

// GetPlanByID retrieves plan by ID
func (r *BillingRepository) GetPlanByID(ctx context.Context, planID int64) (*models.BillingPlan, error) {
	query := `SELECT id, price FROM plans2 WHERE id = ?`

	var plan models.BillingPlan
	err := database.MySQLDB.QueryRowContext(ctx, query, planID).Scan(&plan.ID, &plan.Price)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}
	return &plan, nil
}

// UpdateBalance updates user balance
func (r *BillingRepository) UpdateBalance(ctx context.Context, contract string, newBalance float64) error {
	query := `UPDATE users SET balance = ? WHERE contract = ?`

	_, err := database.MySQLDB.ExecContext(ctx, query, newBalance, contract)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}
	return nil
}

// AddPayment adds a payment record
func (r *BillingRepository) AddPayment(ctx context.Context, pay *models.BillingPay) error {
	query := `INSERT INTO pays (mid, cash, time, admin, reason, coment, bonus, flag) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := database.MySQLDB.ExecContext(ctx, query,
		pay.UserID, pay.Amount, pay.Time, pay.Admin, pay.Reason, pay.Comment, pay.Bonus, pay.Flag)
	if err != nil {
		return fmt.Errorf("failed to add payment: %w", err)
	}
	return nil
}

// SetUserStatus updates user status (state)
func (r *BillingRepository) SetUserStatus(ctx context.Context, contract string, status string) error {
	query := `UPDATE users SET state = ? WHERE contract = ?`
	_, err := database.MySQLDB.ExecContext(ctx, query, status, contract)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}
	return nil
}

// EnableTemporaryPayment enables temporary payment
func (r *BillingRepository) EnableTemporaryPayment(ctx context.Context, contract string, balance float64) error {
	tx, err := database.MySQLDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update balance
	_, err = tx.ExecContext(ctx, `UPDATE users SET balance = ? WHERE contract = ?`, balance, contract)
	if err != nil {
		return err
	}

	// Set state to 'on'
	_, err = tx.ExecContext(ctx, `UPDATE users SET state = 'on' WHERE contract = ?`, contract)
	if err != nil {
		return err
	}

	// Set t_pay to 1
	_, err = tx.ExecContext(ctx, `UPDATE users SET t_pay = 1 WHERE contract = ?`, contract)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// SearchByPhone searches user in billing by phone number
func (r *BillingRepository) SearchByPhone(ctx context.Context, phone string) (*models.BillingUser, error) {
	query := `SELECT id, contract, fio, telefon, balance, paket, state, t_pay, start_day, srvs, address, grp 
	          FROM users WHERE telefon LIKE ? LIMIT 1`

	// Python uses exact match or like? Usually phone search is exact or ends with.
	// Assuming exact match for now or simple LIKE
	var user models.BillingUser
	err := database.MySQLDB.QueryRowContext(ctx, query, "%"+phone).Scan(
		&user.ID, &user.Contract, &user.Name, &user.Phone, &user.Balance,
		&user.PlanID, &user.Status, &user.TimePay, &user.StartDay, &user.Services,
		&user.Address, &user.Group,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to search by phone: %w", err)
	}
	return &user, nil
}

// SearchByName searches users in billing by name
func (r *BillingRepository) SearchByName(ctx context.Context, name string) ([]models.BillingUser, error) {
	query := `SELECT id, contract, fio, telefon, balance, paket, state, t_pay, start_day, srvs, address, grp 
	          FROM users WHERE fio LIKE ? LIMIT 10`

	rows, err := database.MySQLDB.QueryContext(ctx, query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search by name: %w", err)
	}
	defer rows.Close()

	var users []models.BillingUser
	for rows.Next() {
		var user models.BillingUser
		err := rows.Scan(
			&user.ID, &user.Contract, &user.Name, &user.Phone, &user.Balance,
			&user.PlanID, &user.Status, &user.TimePay, &user.StartDay, &user.Services,
			&user.Address, &user.Group,
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
	query := `SELECT id, contract, fio, telefon, balance, paket, state, t_pay, start_day, srvs, address, grp 
	          FROM users WHERE address LIKE ? LIMIT 10`

	rows, err := database.MySQLDB.QueryContext(ctx, query, "%"+address+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search by address: %w", err)
	}
	defer rows.Close()

	var users []models.BillingUser
	for rows.Next() {
		var user models.BillingUser
		err := rows.Scan(
			&user.ID, &user.Contract, &user.Name, &user.Phone, &user.Balance,
			&user.PlanID, &user.Status, &user.TimePay, &user.StartDay, &user.Services,
			&user.Address, &user.Group,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

// GetUsersByGroup retrieves users by group ID
func (r *BillingRepository) GetUsersByGroup(ctx context.Context, groupID int) ([]models.BillingUser, error) {
	query := `SELECT id, contract, fio, telefon, balance, paket, state, t_pay, start_day, srvs, address, grp 
	          FROM users WHERE grp = ?`

	rows, err := database.MySQLDB.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by group: %w", err)
	}
	defer rows.Close()

	var users []models.BillingUser
	for rows.Next() {
		var user models.BillingUser
		err := rows.Scan(
			&user.ID, &user.Contract, &user.Name, &user.Phone, &user.Balance,
			&user.PlanID, &user.Status, &user.TimePay, &user.StartDay, &user.Services,
			&user.Address, &user.Group,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
