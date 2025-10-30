package database

import (
	"effective_mobile/internal/entity"
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create db connection: %w", err)
	}

	err = db.AutoMigrate(&entity.User{}, &entity.Subscription{})
	if err != nil {
		return nil, fmt.Errorf("unable to migrate database: %w", err)
	}

	slog.Info("Successfully connected to PostgreSQL")
	return db, nil
}
