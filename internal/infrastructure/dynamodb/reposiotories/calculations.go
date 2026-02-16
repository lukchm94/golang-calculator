package dynamoRepo

import (
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
	tableName dynamodb.Table
}

func NewCalculationsDynamoRepository(db *dynamodb.DynamoDbClient, logger *slog.Logger, tableName string) (*CalculationsDynamoRepository, error) {
	validTableName, err := db.Tables.Validate(tableName)

	if err != nil {
		logger.Error("Failed to validate table name", "error", err)

		return nil, dynamodb.TableNotFoundError{TableName: tableName}
	}

	return &CalculationsDynamoRepository{
		logger:    logger,
		db:        db,
		tableName: validTableName,
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

	err = r.saveRecord(ctx, r.tableName.Name, record)

	if err != nil {
		r.logger.Error("Failed to save calculation record to DynamoDB", "error", err)
		return err
	}

	return nil
}

func (r *CalculationsDynamoRepository) saveRecord(ctx context.Context, tableName string, record dynamodbModels.CalculationDynamoRecord) error {
	item, err := attributevalue.MarshalMap(record)
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
