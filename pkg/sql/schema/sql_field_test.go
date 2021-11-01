package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSQLField_DDL(t *testing.T) {
	type fields struct {
		Name     string
		DataType SQLDataType
		Collate  string
	}
	tests := []struct {
		name   string
		fields fields
		want   DDL
	}{
		{
			name:   "simple",
			fields: fields{Name: "field", DataType: SQLDataType{Int: &SQLDataTypeInt{}}},
			want:   NewDDL("? int", "field"),
		}, {
			name:   "with collate",
			fields: fields{Name: "field", Collate: "collate", DataType: SQLDataType{Int: &SQLDataTypeInt{}}},
			want:   NewDDL("? int collate ?", "field", "collate"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SQLField{
				Name:     tt.fields.Name,
				DataType: tt.fields.DataType,
				Collate:  tt.fields.Collate,
			}
			assert.Equal(t, tt.want, s.DDL())
		})
	}
}
