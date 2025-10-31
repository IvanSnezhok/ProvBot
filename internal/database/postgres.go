package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"provbot/internal/utils"
)

var PostgresDB *pgxpool.Pool

// InitPostgres initializes PostgreSQL connection pool
func InitPostgres(config *utils.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dsn := config.GetPostgresDSN()
	
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse PostgreSQL config: %w", err)
	}

	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 5 * time.Minute
	poolConfig.MaxConnIdleTime = 1 * time.Minute

	PostgresDB, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("failed to create PostgreSQL connection pool: %w", err)
	}

	// Test connection
	if err := PostgresDB.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	utils.Logger.Info("PostgreSQL connection established successfully")
	return nil
}

// ClosePostgres closes PostgreSQL connection pool
func ClosePostgres() {
	if PostgresDB != nil {
		PostgresDB.Close()
		utils.Logger.Info("PostgreSQL connection closed")
	}
}

