package sqlmanager

import (
	"errors"
	"github.com/nrc-no/core/pkg/sql/schema"
)

func (s writer) handleCreateTable(createTable sqlActionCreateTable) (writer, error) {
	state := s.State

	// checking if there is already a table existing
	_, err := state.Tables.GetTable(createTable.sqlTable.Name)
	if err != nil {
		if errors.Is(err, schema.ErrTableNotFound) {
			// the table does not already exist. Create a new one
			state.Tables = append(state.Tables, createTable.sqlTable)
			s.Statements = append(s.Statements, createTable.sqlTable.DDL())
		} else {
			return writer{}, err
		}
	} else {
		// todo update an existing table to match the desired state
	}
	s.State = state
	return s, nil
}
