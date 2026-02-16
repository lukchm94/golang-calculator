package dynamodbModels

import (
	"fmt"
	"time"
)

const TABLE_CALCULATIONS = "Calculations"

type CalculationDynamoRecord struct {
	// PK: Date (YYYY-MM-DD)
	PK string `dynamodbav:"pk"`
	// SK: OP#Addition#SessionID
	SK        string    `dynamodbav:"sk"`
	Operation string    `dynamodbav:"operation"`
	Number1   float64   `dynamodbav:"number1"`
	Number2   float64   `dynamodbav:"number2"`
	Result    float64   `dynamodbav:"result"`
	CreatedAt time.Time `dynamodbav:"created_at"`
}

// NewCalculationRecord is a helper to format your keys correctly
func NewCalculationRecord(op string, n1, n2, res float64, sessionId string) CalculationDynamoRecord {
	now := time.Now()

	return CalculationDynamoRecord{
		PK:        now.Format("2006-01-02"), // YYYY-MM-DD
		SK:        fmt.Sprintf("OP#%s#%s", op, sessionId),
		Operation: op,
		Number1:   n1,
		Number2:   n2,
		Result:    res,
		CreatedAt: now,
	}
}
