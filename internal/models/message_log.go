package models

import "time"

// MessageDirection represents message direction
type MessageDirection string

const (
	DirectionIncoming MessageDirection = "incoming"
	DirectionOutgoing MessageDirection = "outgoing"
)

// MessageLog represents a logged message
type MessageLog struct {
	ID         int64           `json:"id"`
	UserID     *int            `json:"user_id" db:"user_id"`
	TelegramID int64           `json:"telegram_id" db:"telegram_id"`
	Direction  MessageDirection `json:"direction"`
	MessageText *string         `json:"message_text" db:"message_text"`
	MessageID  *int64          `json:"message_id" db:"message_id"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}

