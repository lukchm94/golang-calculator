package dynamodbModels

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

func TestCalculationDynamoRecordMarshalMap_UsesTerraformKeyNames(t *testing.T) {
	record := CalculationDynamoRecord{
		PartitionKey: "2026-04-13",
		SortKey:      "OP#DIVIDE#session",
		Operation:    "DIVIDE",
		Number1:      1423,
		Number2:      25,
		Result:       56.92,
	}

	item, err := attributevalue.MarshalMap(record)
	if err != nil {
		t.Fatalf("expected record to marshal successfully, got error: %v", err)
	}

	if _, ok := item["PartitionKey"]; !ok {
		t.Fatalf("expected marshalled item to include PartitionKey, got keys: %v", item)
	}

	if _, ok := item["SortKey"]; !ok {
		t.Fatalf("expected marshalled item to include SortKey, got keys: %v", item)
	}

	if _, ok := item["pk"]; ok {
		t.Fatalf("expected marshalled item not to include legacy pk key, got keys: %v", item)
	}

	if _, ok := item["sk"]; ok {
		t.Fatalf("expected marshalled item not to include legacy sk key, got keys: %v", item)
	}
}
