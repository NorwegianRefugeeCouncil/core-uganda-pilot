package schema

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorIs(t *testing.T) {

	tests := []struct {
		name  string
		err   error
		isErr error
	}{
		{
			name:  "tableNotFound",
			err:   newTableNotFoundErr("tableName"),
			isErr: ErrTableNotFound,
		},
		{
			name:  "columnNotFound",
			err:   newColumnNotFoundErr("columnName"),
			isErr: ErrColumnNotFound,
		},
		{
			name:  "tableConstraintNotFound",
			err:   newTableConstraintNotFoundErr("constraintName"),
			isErr: ErrTableConstraintNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.True(t, errors.Is(test.err, test.isErr))
		})
	}

}
