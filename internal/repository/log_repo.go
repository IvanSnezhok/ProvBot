package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"provbot/internal/database"
	"provbot/internal/models"
)

type LogRepository struct{}

func NewLogRepository() *LogRepository {
	return &LogRepository{}
}

// LogMessage logs a message to the database
func (r *LogRepository) LogMessage(ctx context.Context, log *models.MessageLog) error {
	query := `INSERT INTO message_logs (user_id, telegram_id, direction, message_text, message_id)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	
	err := database.PostgresDB.QueryRow(ctx, query,
		log.UserID, log.TelegramID, log.Direction, log.MessageText, log.MessageID,
	).Scan(&log.ID, &log.CreatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	return nil
}

// LogBotEvent logs a bot event to the database
func (r *LogRepository) LogBotEvent(ctx context.Context, log *models.BotLog) error {
	var fieldsJSON []byte
	if log.Fields != nil {
		var err error
		fieldsJSON, err = json.Marshal(log.Fields)
		if err != nil {
			return fmt.Errorf("failed to marshal fields: %w", err)
		}
	}
	
	query := `INSERT INTO bot_logs (level, message, fields) VALUES ($1, $2, $3) RETURNING id, created_at`
	
	err := database.PostgresDB.QueryRow(ctx, query, log.Level, log.Message, fieldsJSON).Scan(&log.ID, &log.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to log bot event: %w", err)
	}
	return nil
}

// GetMessageLogs retrieves message logs for a user
func (r *LogRepository) GetMessageLogs(ctx context.Context, telegramID int64, limit int) ([]models.MessageLog, error) {
	query := `SELECT id, user_id, telegram_id, direction, message_text, message_id, created_at
	          FROM message_logs WHERE telegram_id = $1 ORDER BY created_at DESC LIMIT $2`
	
	rows, err := database.PostgresDB.Query(ctx, query, telegramID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get message logs: %w", err)
	}
	defer rows.Close()

	var logs []models.MessageLog
	for rows.Next() {
		var log models.MessageLog
		err := rows.Scan(
			&log.ID, &log.UserID, &log.TelegramID, &log.Direction,
			&log.MessageText, &log.MessageID, &log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message log: %w", err)
		}
		logs = append(logs, log)
	}
	return logs, nil
}

// GetBotLogs retrieves bot logs by level
func (r *LogRepository) GetBotLogs(ctx context.Context, level models.LogLevel, limit int) ([]models.BotLog, error) {
	query := `SELECT id, level, message, fields, created_at FROM bot_logs 
	          WHERE level = $1 ORDER BY created_at DESC LIMIT $2`
	
	rows, err := database.PostgresDB.Query(ctx, query, level, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot logs: %w", err)
	}
	defer rows.Close()

	var logs []models.BotLog
	for rows.Next() {
		var log models.BotLog
		var fieldsJSON []byte
		err := rows.Scan(&log.ID, &log.Level, &log.Message, &fieldsJSON, &log.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bot log: %w", err)
		}
		
		if len(fieldsJSON) > 0 {
			if err := json.Unmarshal(fieldsJSON, &log.Fields); err != nil {
				log.Fields = make(map[string]interface{})
			}
		}
		logs = append(logs, log)
	}
	return logs, nil
}

