package schema

import (
	"github.com/nrc-no/core/pkg/utils/pointers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSQLDataType_DDL(t *testing.T) {
	type fields struct {
		Array           *SQLDataTypeArray
		VarChar         *SQLDataTypeVarChar
		Int             *SQLDataTypeInt
		SmallInt        *SQLDataTypeSmallInt
		BigInt          *SQLDataTypeBigInt
		Numeric         *SQLDataTypeNumeric
		DoublePrecision *SQLDataTypeDoublePrecision
		Serial          *SQLDataTypeSerial
		SmallSerial     *SQLDataTypeSmallSerial
		BigSerial       *SQLDataTypeBigSerial
		Boolean         *SQLDataTypeBoolean
		Float           *SQLDataTypeFloat
		Timestamp       *SQLDataTypeTimestamp
		Date            *SQLDataTypeDate
		Time            *SQLDataTypeTime
		Interval        *SQLDataTypeInterval
	}
	tests := []struct {
		name   string
		fields fields
		want   DDL
	}{
		{
			name:   "varchar",
			fields: fields{VarChar: &SQLDataTypeVarChar{Length: 10}},
			want:   NewDDL("varchar(10)"),
		}, {
			name:   "int",
			fields: fields{Int: &SQLDataTypeInt{}},
			want:   NewDDL("int"),
		}, {
			name:   "bigint",
			fields: fields{BigInt: &SQLDataTypeBigInt{}},
			want:   NewDDL("bigint"),
		}, {
			name:   "smallint",
			fields: fields{SmallInt: &SQLDataTypeSmallInt{}},
			want:   NewDDL("smallint"),
		}, {
			name:   "numeric",
			fields: fields{Numeric: &SQLDataTypeNumeric{}},
			want:   NewDDL("numeric"),
		}, {
			name:   "numeric with precision",
			fields: fields{Numeric: &SQLDataTypeNumeric{Precision: pointers.Int(10)}},
			want:   NewDDL("numeric(?)", 10),
		}, {
			name:   "numeric with precision and scale",
			fields: fields{Numeric: &SQLDataTypeNumeric{Precision: pointers.Int(10), Scale: pointers.Int(5)}},
			want:   NewDDL("numeric(?, ?)", 10, 5),
		}, {
			name:   "double precision",
			fields: fields{DoublePrecision: &SQLDataTypeDoublePrecision{}},
			want:   NewDDL("double precision"),
		}, {
			name:   "serial",
			fields: fields{Serial: &SQLDataTypeSerial{}},
			want:   NewDDL("serial"),
		}, {
			name:   "small serial",
			fields: fields{SmallSerial: &SQLDataTypeSmallSerial{}},
			want:   NewDDL("smallserial"),
		}, {
			name:   "big serial",
			fields: fields{BigSerial: &SQLDataTypeBigSerial{}},
			want:   NewDDL("bigserial"),
		}, {
			name:   "boolean",
			fields: fields{Boolean: &SQLDataTypeBoolean{}},
			want:   NewDDL("boolean"),
		}, {
			name:   "float",
			fields: fields{Float: &SQLDataTypeFloat{}},
			want:   NewDDL("float"),
		}, {
			name:   "float with precision",
			fields: fields{Float: &SQLDataTypeFloat{Precision: 10}},
			want:   NewDDL("float(?)", 10),
		}, {
			name:   "timestamp",
			fields: fields{Timestamp: &SQLDataTypeTimestamp{}},
			want:   NewDDL("timestamp"),
		}, {
			name:   "timestamp with tz",
			fields: fields{Timestamp: &SQLDataTypeTimestamp{Timezone: &TimestampWithTimeZone}},
			want:   NewDDL("timestamp with time zone"),
		}, {
			name:   "timestamp without tz",
			fields: fields{Timestamp: &SQLDataTypeTimestamp{Timezone: &TimestampWithoutTimeZone}},
			want:   NewDDL("timestamp without time zone"),
		}, {
			name:   "timestamp with p",
			fields: fields{Timestamp: &SQLDataTypeTimestamp{Digits: pointers.Int(5)}},
			want:   NewDDL("timestamp ?", 5),
		}, {
			name:   "timestamp with p with tz",
			fields: fields{Timestamp: &SQLDataTypeTimestamp{Digits: pointers.Int(5), Timezone: &TimestampWithTimeZone}},
			want:   NewDDL("timestamp ? with time zone", 5),
		}, {
			name:   "timestamp with p without tz",
			fields: fields{Timestamp: &SQLDataTypeTimestamp{Digits: pointers.Int(5), Timezone: &TimestampWithoutTimeZone}},
			want:   NewDDL("timestamp ? without time zone", 5),
		}, {
			name:   "time with tz",
			fields: fields{Time: &SQLDataTypeTime{Timezone: &TimestampWithTimeZone}},
			want:   NewDDL("time with time zone"),
		}, {
			name:   "time without tz",
			fields: fields{Time: &SQLDataTypeTime{Timezone: &TimestampWithoutTimeZone}},
			want:   NewDDL("time without time zone"),
		}, {
			name:   "time with p",
			fields: fields{Time: &SQLDataTypeTime{Digits: pointers.Int(5)}},
			want:   NewDDL("time ?", 5),
		}, {
			name:   "time with p with tz",
			fields: fields{Time: &SQLDataTypeTime{Digits: pointers.Int(5), Timezone: &TimestampWithTimeZone}},
			want:   NewDDL("time ? with time zone", 5),
		}, {
			name:   "time with p without tz",
			fields: fields{Time: &SQLDataTypeTime{Digits: pointers.Int(5), Timezone: &TimestampWithoutTimeZone}},
			want:   NewDDL("time ? without time zone", 5),
		}, {
			name:   "interval",
			fields: fields{Interval: &SQLDataTypeInterval{}},
			want:   NewDDL("interval"),
		}, {
			name:   "interval with p",
			fields: fields{Interval: &SQLDataTypeInterval{Digits: pointers.Int(10)}},
			want:   NewDDL("interval ?", 10),
		}, {
			name:   "interval with fields",
			fields: fields{Interval: &SQLDataTypeInterval{Fields: []string{"A"}}},
			want:   NewDDL("interval ?", "A"),
		}, {
			name:   "interval with fields with p",
			fields: fields{Interval: &SQLDataTypeInterval{Fields: []string{"A"}, Digits: pointers.Int(10)}},
			want:   NewDDL("interval ? ?", "A", 10),
		}, {
			name: "array",
			fields: fields{
				Array: &SQLDataTypeArray{
					DataType: SQLDataType{VarChar: &SQLDataTypeVarChar{Length: 10}},
				},
			},
			want: NewDDL("varchar(10)[]"),
		}, {
			name: "array with length",
			fields: fields{
				Array: &SQLDataTypeArray{
					Length:   5,
					DataType: SQLDataType{VarChar: &SQLDataTypeVarChar{Length: 10}},
				},
			},
			want: NewDDL("varchar(10)[5]"),
		}, {
			name: "multi dimensional array",
			fields: fields{
				Array: &SQLDataTypeArray{
					DataType: SQLDataType{
						Array: &SQLDataTypeArray{
							DataType: SQLDataType{
								VarChar: &SQLDataTypeVarChar{Length: 10},
							},
						},
					},
				},
			},
			want: NewDDL("varchar(10)[][]"),
		}, {
			name: "multi dimensional array with length",
			fields: fields{
				Array: &SQLDataTypeArray{
					Length: 1,
					DataType: SQLDataType{
						Array: &SQLDataTypeArray{
							Length: 3,
							DataType: SQLDataType{
								VarChar: &SQLDataTypeVarChar{Length: 10},
							},
						},
					},
				},
			},
			want: NewDDL("varchar(10)[3][1]"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SQLDataType{
				Array:           tt.fields.Array,
				VarChar:         tt.fields.VarChar,
				Int:             tt.fields.Int,
				SmallInt:        tt.fields.SmallInt,
				BigInt:          tt.fields.BigInt,
				Numeric:         tt.fields.Numeric,
				DoublePrecision: tt.fields.DoublePrecision,
				Serial:          tt.fields.Serial,
				SmallSerial:     tt.fields.SmallSerial,
				BigSerial:       tt.fields.BigSerial,
				Boolean:         tt.fields.Boolean,
				Float:           tt.fields.Float,
				Timestamp:       tt.fields.Timestamp,
				Date:            tt.fields.Date,
				Time:            tt.fields.Time,
				Interval:        tt.fields.Interval,
			}
			assert.Equal(t, tt.want, s.DDL())

		})
	}
}
