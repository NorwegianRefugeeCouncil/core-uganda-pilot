package sqlschema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSQLSchema_DDL(t *testing.T) {
	type fields struct {
		Name string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "simple",
			fields: fields{Name: "schema"},
			want:   "CREATE SCHEMA schema;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SQLSchema{
				Name: tt.fields.Name,
			}
			assert.Equal(t, tt.want, s.DDL())
		})
	}
}
