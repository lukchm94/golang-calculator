package postgres

import (
	"fmt"
	"log/slog"

	postgresModels "app/internal/infrastructure/postgres/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

func NewGormClient(config PostgresConfig, logger *slog.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s dbname=%s port=%s sslmode=disable",
		config.Host, config.DbName, config.Port)

	logger.Debug("Setting up Postgres DB with the following", "DSN", dsn)

	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		return nil, ErrPostgresInit
	}

	return db, nil
}

func InitPostgresTables(db *gorm.DB, logger *slog.Logger) error {
	err := db.AutoMigrate(&postgresModels.UserPostgres{})

	if err != nil {
		logger.Error("Failed to initialize Postgres tables", "error", err)
		return ErrPostgresTablesInit
	}

	return nil
}
