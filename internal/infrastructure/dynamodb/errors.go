package dynamodb

import "errors"

type DynamoDbError error

var (
	ErrDynamoDB       DynamoDbError = errors.New("general DynamoDB error")
	ErrTablesCreation DynamoDbError = errors.New("failed to create DynamoDB tables")
)

type TableNotFoundError struct {
	TableName string
}

func (e TableNotFoundError) Error() string {
	return "Table '" + e.TableName + "' not found in DynamoDB"
}

type SavingRecordError struct {
	Err    error
	Record any
}

func (e SavingRecordError) Error() string {
	return "Failed to save record " + string(e.Record.(string)) + " to DynamoDB: " + e.Err.Error()
}
