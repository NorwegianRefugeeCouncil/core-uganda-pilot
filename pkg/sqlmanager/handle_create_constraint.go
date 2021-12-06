package sqlmanager

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/nrc-no/core/pkg/sql/schema"
)

func (s sqlManager) handleCreateConstraint(constraint sqlActionCreateConstraint) (sqlManager, error) {

	state := s.State

	// getting the table on which to append the new constraint
	table, err := state.Tables.GetTable(constraint.tableName)
	if err != nil {
		return sqlManager{}, err
	}

	// checking if the constraint already exists
	_, err = table.Constraints.GetConstraint(constraint.sqlConstraint.Name)
	if err != nil {
		if errors.Is(err, schema.ErrTableConstraintNotFound) {
			// the constraint does not already exist. Create a new constraint

			table.Constraints = append(table.Constraints, constraint.sqlConstraint)
			state.Tables, err = state.Tables.ReplaceTable(table)
			if err != nil {
				return sqlManager{}, err
			}

			// create the DDL Statement for creating the constraint
			addConstraintDDL := schema.NewDDL(
				fmt.Sprintf("alter table %s.%s add ",
					pq.QuoteIdentifier(constraint.schemaName),
					pq.QuoteIdentifier(constraint.tableName))).
				Merge(constraint.sqlConstraint.DDL()).
				Write(";")

			s.Statements = append(s.Statements, addConstraintDDL)

		} else {
			return sqlManager{}, err
		}

	} else {
		// todo modify the constraint to match the desired state
	}

	s.State = state
	return s, nil
}
