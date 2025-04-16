package executors

import (
	stdlib_helpers "backend-testing-module-checker/stdlib/helpers"
	stdlib_types "backend-testing-module-checker/stdlib/types"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	DB_RESULTS = "db-results"
)

type DBQueryArgs struct {
	Query            string `yaml:"query"`
	ConnectionString string `yaml:"connection-string"`
}

func tryInjectQuery(query string, storage stdlib_types.Storage) (string, error) {
	injectedQuery, err := stdlib_helpers.TryInject(query, storage)
	if err != nil {
		return "", fmt.Errorf("failed to inject query parameters: %v", err)
	}
	return injectedQuery, nil
}

func DBQueryHandler(id string, args any, storage stdlib_types.Storage) stdlib_types.ExecutorResult {
	arguments, ok := args.(*DBQueryArgs)
	if !ok {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.EF,
			Comment: "failed to retrieve arguments. Checker bug",
		}
	}

	injectedQuery, err := tryInjectQuery(arguments.Query, storage)
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.PE,
			Comment: err.Error(),
		}
	}

	db, err := sql.Open("postgres", arguments.ConnectionString)
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.PE,
			Comment: fmt.Sprintf("failed to connect to database: %v", err),
		}
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.PE,
			Comment: fmt.Sprintf("failed to ping database: %v", err),
		}
	}

	rows, err := db.Query(injectedQuery)
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.PE,
			Comment: fmt.Sprintf("failed to execute query: %v", err),
		}
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.PE,
			Comment: fmt.Sprintf("failed to get columns: %v", err),
		}
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return stdlib_types.ExecutorResult{
				Verdict: stdlib_types.PE,
				Comment: fmt.Sprintf("failed to scan row: %v", err),
			}
		}

		rowData := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			switch v := val.(type) {
			case []byte:
				rowData[col] = string(v)
			default:
				rowData[col] = v
			}
		}
		results = append(results, rowData)
	}

	if err := rows.Err(); err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.PE,
			Comment: fmt.Sprintf("error after scanning rows: %v", err),
		}
	}

	if len(results) == 0 {
		results = []map[string]interface{}{}
	}

	jsonResults, err := json.Marshal(results)
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.PE,
			Comment: fmt.Sprintf("failed to marshal results to JSON: %v", err),
		}
	}

	storage[id] = map[string]interface{}{
		DB_RESULTS: string(jsonResults),
	}

	return stdlib_types.ExecutorResult{
		Verdict: stdlib_types.OK,
	}
}
