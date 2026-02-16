package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

func (i *InitDynamoDb) Init() {
	i.logger.Info("Initializing DynamoDB connection...")
	// Here you would add the actual initialization code for DynamoDB, such as creating a client, setting up tables, etc.
}

// EnsureTablesExist checks the provided list of tables and creates them if missing.
func (i *InitDynamoDb) EnsureTablesExist(ctx context.Context) error {
	for _, table := range i.dbClient.Tables.Table {
		_, err := i.dbClient.Client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(table.Name),
		})

		if err != nil {
			var notFound *types.ResourceNotFoundException
			if errors.As(err, &notFound) {
				i.logger.Info("Table not found creating it...", "name", table.Name)
				if err := i.createCalculationsTable(ctx, table.Name); err != nil {
					return err
				}
				continue
			}
			return fmt.Errorf("failed to describe table %s: %w", table.Name, err)
		}
		i.logger.Info("Table already exists", "name", table.Name)
	}
	return nil
}

func (i *InitDynamoDb) createCalculationsTable(ctx context.Context, name string) error {
	_, err := i.dbClient.Client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(name),
		// 1. Define the Attributes (PK and SK)
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("pk"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("sk"), AttributeType: types.ScalarAttributeTypeS},
		},
		// 2. Define the Key Schema (Hash and Range)
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("pk"), KeyType: types.KeyTypeHash},  // Partition Key
			{AttributeName: aws.String("sk"), KeyType: types.KeyTypeRange}, // Sort Key
		},
		// 3. Set Billing Mode (Pay Per Request is easiest for local/dev)
		BillingMode: types.BillingModePayPerRequest,
	})
	return err
}
