package sqlmanager

import (
	"errors"
	"github.com/nrc-no/core/pkg/sql/schema"
)

func (s sqlManager) handleCreateTable(createTable sqlActionCreateTable) (sqlManager, error) {
	state := s.State

	// checking if there is already a table existing
	_, err := state.Tables.GetTable(createTable.sqlTable.Name)
	if err != nil {
		if errors.Is(err, schema.ErrTableNotFound) {
			// the table does not already exist. Create a new one
			state.Tables = append(state.Tables, createTable.sqlTable)
			s.Statements = append(s.Statements, createTable.sqlTable.DDL())
		} else {
			return sqlManager{}, err
		}
	} else {
		// todo update an existing table to match the desired state
	}
	s.State = state
	return s, nil
}
