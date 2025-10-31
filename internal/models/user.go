package models

import "time"

// User represents a bot user
type User struct {
	ID         int       `json:"id"`
	TelegramID int64     `json:"telegram_id" db:"telegram_id"`
	Username   *string   `json:"username"`
	FirstName  *string   `json:"first_name" db:"first_name"`
	LastName   *string   `json:"last_name" db:"last_name"`
	PhoneNumber *string  `json:"phone_number" db:"phone_number"`
	Contract   *string   `json:"contract" db:"contract"`
	Language   string    `json:"language"`
	IsActive   bool      `json:"is_active" db:"is_active"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// AdminUser represents an admin user
type AdminUser struct {
	ID          int                    `json:"id"`
	UserID      int                    `json:"user_id" db:"user_id"`
	TelegramID  int64                  `json:"telegram_id" db:"telegram_id"`
	Permissions map[string]interface{} `json:"permissions"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
}

