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
	CreatedBy   *int         `json:"created_by" db:"created_by"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	ResolvedAt  *time.Time   `json:"resolved_at" db:"resolved_at"`
}

