package sqlmanager

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/nrc-no/core/pkg/sql/schema"
)

func (s sqlManager) handleCreateColumn(createColumn sqlActionCreateColumn) (sqlManager, error) {

	state := s.State

	// finding the table on which to append the column
	table, err := state.Tables.GetTable(createColumn.tableName)
	if err != nil {
		return sqlManager{}, err
	}

	// checking if the column already exists
	_, err = table.Columns.GetColumn(createColumn.sqlColumn.Name)
	if err != nil {

		if errors.Is(err, schema.ErrColumnNotFound) {
			// the column does not exist. create it

			table.Columns = append(table.Columns, createColumn.sqlColumn)
			state.Tables, err = state.Tables.ReplaceTable(table)
			if err != nil {
				return sqlManager{}, err
			}

			// add the DDL statement to create the column
			addColumnDDL :=
				schema.NewDDL(fmt.Sprintf("alter table %s.%s add ",
					pq.QuoteIdentifier(createColumn.schemaName),
					pq.QuoteIdentifier(createColumn.tableName))).
					Merge(createColumn.sqlColumn.DDL()).
					Write(";")

			s.Statements = append(s.Statements, addColumnDDL)
		} else {
			return sqlManager{}, err
		}

	} else {
		// todo update an existing column to match the desired column
	}
	s.State = state
	return s, nil
}
