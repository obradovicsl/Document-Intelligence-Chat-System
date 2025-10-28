package repository

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() error {
    dsn := os.Getenv("DB_URL")
    
    slog.Info("connecting to database")
    
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }
    
    // Connection pool settings
    sqlDB, err := DB.DB()
    if err != nil {
        return err
    }
    
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(5)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)
    
    // Auto migrate
    slog.Info("running auto migrations")

    if err := DB.AutoMigrate(
        &models.Document{},
        &models.Chat{},
    ); err != nil {
        return fmt.Errorf("failed to run migrations: %w", err)
    }

    
    slog.Info("database initialized successfully")
    return nil
}