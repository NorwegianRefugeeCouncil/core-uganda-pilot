package validation

import (
	"github.com/nrc-no/core/pkg/validation"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestValidateDatabaseName(t *testing.T) {
	tests := []struct {
		name         string
		databaseName string
		field        *validation.Path
		want         validation.ErrorList
	}{
		{
			name:         "valid name",
			databaseName: "some-database",
			field:        validation.NewPath("name"),
			want:         validation.ErrorList{},
		}, {
			name:         "too short",
			databaseName: "a",
			field:        validation.NewPath("name"),
			want: validation.ErrorList{
				validation.TooShort(validation.NewPath("name"), "a", databaseNameMinLength),
			},
		}, {
			name:         "too long",
			databaseName: strings.Repeat("a", databaseNameMaxLength+1),
			field:        validation.NewPath("name"),
			want: validation.ErrorList{
				validation.TooLong(validation.NewPath("name"), strings.Repeat("a", databaseNameMaxLength+1), databaseNameMaxLength),
			},
		}, {
			name:         "invalid",
			databaseName: "  abc  ",
			field:        validation.NewPath("name"),
			want: validation.ErrorList{
				validation.Invalid(validation.NewPath("name"), "  abc  ", invalidDatabaseNameMsg),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateDatabaseName(tt.databaseName, tt.field)
			assert.Equal(t, tt.want, got)
		})
	}
}
