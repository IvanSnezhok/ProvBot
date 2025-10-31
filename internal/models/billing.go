package models

import "time"

// BillingUser represents a user in the billing system
type BillingUser struct {
	ID          int64     `json:"id" db:"id"`
	Username    string    `json:"username" db:"username"`
	Balance     float64   `json:"balance" db:"balance"`
	Status      string    `json:"status" db:"status"`
	ServiceID   *int64    `json:"service_id" db:"service_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// BillingService represents a service in the billing system
type BillingService struct {
	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Status      string    `json:"status" db:"status"`
	Location    string    `json:"location" db:"location"`
	IPAddress   *string   `json:"ip_address" db:"ip_address"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Payment represents a payment transaction
type Payment struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	TelegramPayID string    `json:"telegram_pay_id" db:"telegram_pay_id"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

