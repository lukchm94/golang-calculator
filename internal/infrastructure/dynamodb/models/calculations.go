package dynamodbModels

import (
	"fmt"
	"time"
)

type CalculationDynamoRecord struct {
	PartitionKey string    `dynamodbav:"PartitionKey"` // PK: Date (YYYY-MM-DD)
	SortKey      string    `dynamodbav:"SortKey"`      // SK: OP#Addition#SessionID
	Operation    string    `dynamodbav:"operation"`
	Number1      float64   `dynamodbav:"number1"`
	Number2      float64   `dynamodbav:"number2"`
	Result       float64   `dynamodbav:"result"`
	CreatedAt    time.Time `dynamodbav:"created_at"`
}

// NewCalculationRecord is a helper to format your keys correctly
func NewCalculationRecord(op string, n1, n2, res float64, sessionId string) CalculationDynamoRecord {
	now := time.Now()

	return CalculationDynamoRecord{
		PartitionKey: now.Format("2006-01-02"), // YYYY-MM-DD
		SortKey:      fmt.Sprintf("OP#%s#%s", op, sessionId),
		Operation:    op,
		Number1:      n1,
		Number2:      n2,
		Result:       res,
		CreatedAt:    now,
	}
}
