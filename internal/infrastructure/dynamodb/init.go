package dynamodb

import (
	"log/slog"
)

type InitDynamoDb struct {
	logger   *slog.Logger
	dbClient *DynamoDbClient
}

func NewInitDynamoDb(logger *slog.Logger, dbClient *DynamoDbClient) *InitDynamoDb {
	return &InitDynamoDb{
		logger:   logger,
		dbClient: dbClient,
	}
}
