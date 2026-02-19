package postgres

import "errors"

var (
	ErrPostgresInit = errors.New("failed to start Postgres DB")
	ErrPostgres     = errors.New("general Postgres error")
)
