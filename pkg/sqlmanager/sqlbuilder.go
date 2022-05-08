package sqlmanager

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/nrc-no/core/pkg/sql/schema"
)

type SQLBuilder interface {
	InsertRow(
		schemaName string,
		tableName string,
		columns []string,
		values []interface{},
	) schema.DDL
}

type sqlBuilder struct {
}

func NewSQLBuilder() SQLBuilder {
	return sqlBuilder{}
}

func (s sqlBuilder) InsertRow(
	schemaName string,
	tableName string,
	columns []string,
	values []interface{},
) schema.DDL {
	var placeholders = make([]string, len(columns))
	var columnNames = make([]string, len(columns))

	for i, column := range columns {
		columnNames[i] = pq.QuoteIdentifier(column)
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	addRowQuery := fmt.Sprintf(`INSERT INTO %s.%s (%s) VALUES (%s);`,
		pq.QuoteIdentifier(schemaName),
		pq.QuoteIdentifier(tableName),
		strings.Join(columnNames, ","),
		strings.Join(placeholders, ","),
	)

	query := schema.DDL{
		Query: addRowQuery,
		Args:  values,
	}

	return query
}
