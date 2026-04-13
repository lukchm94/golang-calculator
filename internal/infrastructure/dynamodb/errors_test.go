package dynamodb

import (
	"errors"
	"strings"
	"testing"

	dynamodbModels "app/internal/infrastructure/dynamodb/models"
)

func TestSavingRecordError_WithStructuredRecord(t *testing.T) {
	innerErr := errors.New("put item failed")
	record := dynamodbModels.CalculationDynamoRecord{
		PartitionKey: "2026-04-13",
		SortKey:      "OP#DIVIDE#session",
		Operation:    "DIVIDE",
		Number1:      1423,
		Number2:      25,
		Result:       56.92,
	}

	err := SavingRecordError{Err: innerErr, Record: record}

	msg := err.Error()

	if !strings.Contains(msg, "failed to save record") {
		t.Fatalf("expected formatted error prefix, got %q", msg)
	}

	if !strings.Contains(msg, "OP#DIVIDE#session") {
		t.Fatalf("expected structured record details in error, got %q", msg)
	}

	if !errors.Is(err, innerErr) {
		t.Fatal("expected SavingRecordError to unwrap the inner error")
	}
}
