package dynamoRepo

import (
	"app/cmd/config"
	calculatorApplication "app/internal/application/calculator"
	calculatorDomain "app/internal/domain/calculator"
	dynamodb "app/internal/infrastructure/dynamodb"
	dynamodbModels "app/internal/infrastructure/dynamodb/models"
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	dynamodbSDK "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type CalculationsDynamoRepository struct {
	logger    *slog.Logger
	db        *dynamodb.DynamoDbClient
	tableName string
}

func NewCalculationsDynamoRepository(
	db *dynamodb.DynamoDbClient,
	logger *slog.Logger,
	awsConfig config.AwsConfig,
) (*CalculationsDynamoRepository, error) {

	tableName := string(awsConfig.Prefix) + dynamodb.CalculationsTable

	logger.Info("Initializing CalculationsDynamoRepository", "tableName", tableName)

	return &CalculationsDynamoRepository{
		logger:    logger,
		db:        db,
		tableName: tableName,
	}, nil
}

func (r *CalculationsDynamoRepository) Save(ctx context.Context, input calculatorApplication.SavedCalculationInput) error {

	stringOperation, err := r.convertMathSignToString(input)

	if err != nil {
		return err
	}
	record := dynamodbModels.NewCalculationRecord(
		stringOperation,
		input.CalculationInput.Number1,
		input.CalculationInput.Number2,
		input.Result.Result,
		input.SessionId,
	)

	r.logger.Debug("Saving calculation record", "record", record)

	err = r.saveRecord(ctx, r.tableName, record)

	if err != nil {
		r.logger.Error("Failed to save calculation record to DynamoDB", "error", err)
		return err
	}

	return nil
}

func (r *CalculationsDynamoRepository) saveRecord(ctx context.Context, tableName string, record dynamodbModels.CalculationDynamoRecord) error {
	item, err := attributevalue.MarshalMap(record)

	r.logger.Debug("Marshalling record for DynamoDB", "record", record)

	if err != nil {
		return dynamodb.SavingRecordError{Err: err, Record: record}
	}

	var putItemInput = &dynamodbSDK.PutItemInput{
		TableName: &tableName,
		Item:      item,
	}

	r.logger.Info("Saving record to DynamoDB", "table", tableName, "record", record)
	result, err := r.db.Client.PutItem(ctx, putItemInput)

	if err != nil {
		return dynamodb.SavingRecordError{Err: err, Record: record}
	}

	r.logger.Info("Record saved to DynamoDB", "result", result)

	return nil
}

func (r *CalculationsDynamoRepository) convertMathSignToString(input calculatorApplication.SavedCalculationInput) (string, error) {
	r.logger.Info("Converting operation type to string", "operation", input.Operation)

	switch input.Operation {
	case calculatorDomain.Add:
		return "ADD", nil

	case calculatorDomain.Substract:
		return "SUBSTRACT", nil

	case calculatorDomain.Multiply:
		return "MULTIPLY", nil

	case calculatorDomain.Divide:
		return "DIVIDE", nil

	default:
		r.logger.Error("Invalid operation type", "operation", input.Operation)
		return "", calculatorDomain.ErrInvalidOperation
	}
}
