package models

import "time"

// OutageStatus represents outage status
type OutageStatus string

const (
	OutageStatusActive    OutageStatus = "active"
	OutageStatusResolved  OutageStatus = "resolved"
	OutageStatusCancelled OutageStatus = "cancelled"
)

// Outage represents a service outage
type Outage struct {
	ID          int          `json:"id"`
	Location    string       `json:"location"`
	Description string       `json:"description"`
	Status      OutageStatus `json:"status"`
	GroupID     *string      `json:"group_id" db:"group_id"`     // Billing group ID for targeted outage
	Street      *string      `json:"street" db:"street"`         // Street for targeted outage
	CreatedBy   *int         `json:"created_by" db:"created_by"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	ResolvedAt  *time.Time   `json:"resolved_at" db:"resolved_at"`
}

// UserAlarmStatus represents user's alarm notification status
type UserAlarmStatus struct {
	TelegramID int64  `json:"telegram_id" db:"telegram_id"`
	Contract   string `json:"contract" db:"contract"`
	HasAlarm   bool   `json:"has_alarm" db:"has_alarm"`
}

