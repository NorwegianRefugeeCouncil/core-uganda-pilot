package utils

import (
	"database/sql"
	"fmt"

	"github.com/nrc-no/core/pkg/server/data/api"
)

// SQLResultReader is a wrapper around sql.Rows that implements the ResultReader api.
// It takes care of deserializing the rows into a map of column name to value.
type SQLResultReader struct {
	rows     *sql.Rows
	columns  []string
	values   []api.Value
	pointers []interface{}
}

// Err returns the error, if any, that was encountered during iteration.
func (r SQLResultReader) Err() error {
	return r.rows.Err()
}

// Close closes the Rows, preventing further enumeration. If Next is called
func (r SQLResultReader) Close() error {
	return r.rows.Close()
}

// Next prepares the next result row for reading. It returns true if there is
func (r SQLResultReader) Next() bool {
	return r.rows.Next()
}

// Read reads the next result row
func (r SQLResultReader) Read(columnKinds []api.ValueKind) (map[string]api.Value, error) {
	var values []api.Value
	var pointers []interface{}
	for _, kind := range columnKinds {
		switch kind {
		case api.ValueKindString:
			values = append(values, api.Value{Kind: kind, String: &api.String{}})
			pointers = append(pointers, values[len(values)-1].String)
		case api.ValueKindInt:
			values = append(values, api.Value{Kind: kind, Int: &api.Int{}})
			pointers = append(pointers, values[len(values)-1].Int)
		case api.ValueKindFloat:
			values = append(values, api.Value{Kind: kind, Float: &api.Float{}})
			pointers = append(pointers, values[len(values)-1].Float)
		case api.ValueKindBool:
			values = append(values, api.Value{Kind: kind, Bool: &api.Bool{}})
			pointers = append(pointers, values[len(values)-1].Bool)
		default:
			return nil, fmt.Errorf("unsupported value kind: %v", kind)
		}
	}
	if err := r.rows.Scan(pointers...); err != nil {
		return nil, err
	}
	result := make(map[string]api.Value)
	for i, column := range r.columns {
		result[column] = values[i]
	}
	return result, nil
}

// NewSQLResultReader creates a new SQLResultReader from the given sql.Rows.
func NewSQLResultReader(rows *sql.Rows) (*SQLResultReader, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]api.Value, len(columns))
	pointers := make([]interface{}, len(columns))
	for i := range values {
		pointers[i] = &values[i]
	}
	return &SQLResultReader{
		rows:     rows,
		columns:  columns,
		values:   values,
		pointers: pointers,
	}, nil
}
