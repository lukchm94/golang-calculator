package postgres

import (
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
}

func NewGormClient(config PostgresConfig, logger *slog.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.host, config.user, config.password, config.dbName, config.port)

	logger.Debug("Setting up Postgres DB with the following", "DSN", dsn)

	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		return nil, ErrPostgresInit
	}

	return db, nil
}
