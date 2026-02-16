package dynamodb

import "fmt"

type Table struct {
	Name string
}

type DynamoTables struct {
	Table []Table
}

func NewDynamoTables(tableNames []string) (*DynamoTables, error) {
	if len(tableNames) == 0 {
		return nil, fmt.Errorf("table names list cannot be empty")
	}

	tables := make([]Table, len(tableNames))
	for i, name := range tableNames {
		if name == "" {
			return nil, fmt.Errorf("table name at index %d cannot be empty", i)
		}
		tables[i] = Table{Name: name}
	}

	return &DynamoTables{Table: tables}, nil
}

func (t *DynamoTables) List() []string {
	result := make([]string, len(t.Table))
	for i, table := range t.Table {
		result[i] = table.Name
	}
	return result
}

func (t *DynamoTables) Validate(tableName string) (Table, error) {
	if len(t.Table) == 0 {
		return Table{}, fmt.Errorf("DynamoTables must have at least one table")
	}
	for _, table := range t.Table {
		if table.Name == tableName {
			return table, nil
		}
	}
	return Table{}, fmt.Errorf("Table '%s' is not in DynamoTables", tableName)
}
