package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration for the application
type Config struct {
	// Telegram Bot
	TelegramBotToken string

	// PostgreSQL
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSSLMode  string

	// MySQL (Billing)
	MySQLHost     string
	MySQLPort     int
	MySQLUser     string
	MySQLPassword string
	MySQLDB       string

	// Application
	AppEnv        string
	LogLevel      string
	BotWebhookURL string
	BotWebhookPort int
	ProviderToken string // Telegram Payments provider token

	// Admin
	AdminTelegramIDs []int64

	// Scheduler
	SchedulerWorkers int // Number of workers for parallel processing
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnvAsInt("POSTGRES_PORT", 5432),
		PostgresUser:     getEnv("POSTGRES_USER", "provbot"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:       getEnv("POSTGRES_DB", "provbot_db"),
		PostgresSSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		
		MySQLHost:     getEnv("MYSQL_HOST", "localhost"),
		MySQLPort:     getEnvAsInt("MYSQL_PORT", 3306),
		MySQLUser:     getEnv("MYSQL_USER", "billing_user"),
		MySQLPassword: getEnv("MYSQL_PASSWORD", ""),
		MySQLDB:       getEnv("MYSQL_DB", "billing_db"),
		
		AppEnv:        getEnv("APP_ENV", "development"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		BotWebhookURL: getEnv("BOT_WEBHOOK_URL", ""),
		BotWebhookPort: getEnvAsInt("BOT_WEBHOOK_PORT", 8443),
		ProviderToken: getEnv("PROVIDER_TOKEN", ""),
		SchedulerWorkers: getEnvAsInt("SCHEDULER_WORKERS", 10),
	}

	// Parse admin IDs
	adminIDsStr := getEnv("ADMIN_TELEGRAM_IDS", "")
	if adminIDsStr != "" {
		ids := strings.Split(adminIDsStr, ",")
		for _, idStr := range ids {
			id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
			if err == nil {
				config.AdminTelegramIDs = append(config.AdminTelegramIDs, id)
			}
		}
	}

	// Validate required fields
	if config.TelegramBotToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	return config, nil
}

// GetPostgresDSN returns PostgreSQL connection string
func (c *Config) GetPostgresDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresSSLMode)
}

// GetMySQLDSN returns MySQL connection string
func (c *Config) GetMySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
		c.MySQLUser, c.MySQLPassword, c.MySQLHost, c.MySQLPort, c.MySQLDB)
}

// IsAdmin checks if telegram ID is admin
func (c *Config) IsAdmin(telegramID int64) bool {
	for _, id := range c.AdminTelegramIDs {
		if id == telegramID {
			return true
		}
	}
	return false
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

