package sqlmanager

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/nrc-no/core/pkg/sql/schema"
	"strings"
)

func (s writer) handleInsertRow(insertRow sqlActionInsertRow) (writer, error) {

	var placeholders = make([]string, len(insertRow.columns))
	var columnNames = make([]string, len(insertRow.columns))
	for i, column := range insertRow.columns {
		columnNames[i] = pq.QuoteIdentifier(column)
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	addRowQuery := fmt.Sprintf(`insert into %s.%s (%s) values (%s);`,
		pq.QuoteIdentifier(insertRow.schemaName),
		pq.QuoteIdentifier(insertRow.tableName),
		strings.Join(columnNames, ","),
		strings.Join(placeholders, ","),
	)

	s.Statements = append(s.Statements, schema.DDL{
		Query: addRowQuery,
		Args:  insertRow.values,
	})

	return s, nil
}
