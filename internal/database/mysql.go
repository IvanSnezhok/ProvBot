package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"provbot/internal/utils"
)

var MySQLDB *sql.DB

// InitMySQL initializes MySQL connection for billing system
func InitMySQL(config *utils.Config) error {
	dsn := config.GetMySQLDSN() + "&charset=cp1251&collation=cp1251_general_ci"
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open MySQL connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping MySQL: %w", err)
	}

	MySQLDB = db
	utils.Logger.Info("MySQL connection established successfully")
	return nil
}

// CloseMySQL closes MySQL connection
func CloseMySQL() {
	if MySQLDB != nil {
		MySQLDB.Close()
		utils.Logger.Info("MySQL connection closed")
	}
}

