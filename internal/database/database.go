package database

import (
	"embed"
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/compareflow/compareflow/internal/models"
)

func Initialize(databaseURL string, migrationsFS embed.FS) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if strings.HasPrefix(databaseURL, "sqlite://") {
		dbPath := strings.TrimPrefix(databaseURL, "sqlite://")
		db, err = gorm.Open(sqlite.Open(dbPath), config)
	} else if strings.HasPrefix(databaseURL, "postgres://") || strings.HasPrefix(databaseURL, "postgresql://") {
		db, err = gorm.Open(postgres.Open(databaseURL), config)
	} else {
		return nil, fmt.Errorf("unsupported database URL: %s", databaseURL)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run auto-migrations
	if err := db.AutoMigrate(&models.User{}, &models.Connection{}, &models.Validation{}); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Create default admin user if not exists
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		adminUser := &models.User{
			Username: "admin",
			Email:    "admin@compareflow.com",
		}
		if err := adminUser.SetPassword("admin123"); err != nil {
			return nil, fmt.Errorf("failed to set admin password: %w", err)
		}
		if err := db.Create(adminUser).Error; err != nil {
			return nil, fmt.Errorf("failed to create admin user: %w", err)
		}
	}

	return db, nil
}