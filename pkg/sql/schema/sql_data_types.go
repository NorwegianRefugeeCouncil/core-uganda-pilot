package schema

import (
	"reflect"
	"strconv"
	"strings"
)

type SQLDataType struct {
	Array           *SQLDataTypeArray           `json:"array,omitempty" yaml:"array,omitempty"`
	VarChar         *SQLDataTypeVarChar         `json:"varChar,omitempty" yaml:"varChar,omitempty"`
	Int             *SQLDataTypeInt             `json:"int,omitempty" yaml:"int,omitempty"`
	SmallInt        *SQLDataTypeSmallInt        `json:"smallInt,omitempty" yaml:"smallInt,omitempty"`
	BigInt          *SQLDataTypeBigInt          `json:"bigInt,omitempty" yaml:"bigInt,omitempty"`
	Numeric         *SQLDataTypeNumeric         `json:"numeric,omitempty" yaml:"numeric,omitempty"`
	DoublePrecision *SQLDataTypeDoublePrecision `json:"doublePrecision,omitempty" yaml:"doublePrecision,omitempty"`
	Serial          *SQLDataTypeSerial          `json:"serial,omitempty" yaml:"serial,omitempty"`
	SmallSerial     *SQLDataTypeSmallSerial     `json:"smallSerial,omitempty" yaml:"smallSerial,omitempty"`
	BigSerial       *SQLDataTypeBigSerial       `json:"bigSerial,omitempty" yaml:"bigSerial,omitempty"`
	Boolean         *SQLDataTypeBoolean         `json:"boolean,omitempty" yaml:"boolean,omitempty"`
	Float           *SQLDataTypeFloat           `json:"float,omitempty" yaml:"float,omitempty"`
	Timestamp       *SQLDataTypeTimestamp       `json:"timestamp,omitempty" yaml:"timestamp,omitempty"`
	Date            *SQLDataTypeDate            `json:"date,omitempty" yaml:"date,omitempty"`
	Time            *SQLDataTypeTime            `json:"time,omitempty" yaml:"time,omitempty"`
	Interval        *SQLDataTypeInterval        `json:"interval,omitempty" yaml:"interval,omitempty"`
	Text            *SQLDataTypeText            `json:"text,omitempty" yaml:"text,omitempty"`
}

func (s SQLDataType) DDL() DDL {

	for _, ddLer := range []DDLGenerator{
		s.Array,
		s.VarChar,
		s.SmallInt,
		s.BigInt,
		s.Int,
		s.Numeric,
		s.DoublePrecision,
		s.Serial,
		s.SmallSerial,
		s.BigSerial,
		s.Boolean,
		s.Float,
		s.Timestamp,
		s.Date,
		s.Time,
		s.Interval,
		s.Text,
	} {
		if ddLer != nil {
			if reflect.ValueOf(ddLer).IsNil() {
				continue
			}
			return ddLer.DDL()
		}
	}
	return DDL{}
}

type SQLDataTypeVarChar struct {
	Length int
}

type SQLDataTypeText struct{}

func (S SQLDataTypeText) DDL() DDL {
	return NewDDL("text")
}

func (c SQLDataTypeVarChar) DDL() DDL {
	return DDL{}.WriteF("varchar(%d)", c.Length)
}

type SQLDataTypeInt struct{}

func (c SQLDataTypeInt) DDL() DDL {
	return NewDDL("int")
}

type SQLDataTypeSmallInt struct{}

func (c SQLDataTypeSmallInt) DDL() DDL {
	return DDL{Query: "smallint"}
}

type SQLDataTypeBigInt struct{}

func (c SQLDataTypeBigInt) DDL() DDL {
	return NewDDL("bigint")
}

type SQLDataTypeNumeric struct {
	Precision *int
	Scale     *int
}

func (c SQLDataTypeNumeric) DDL() DDL {
	if c.Precision == nil && c.Scale == nil {
		return NewDDL("numeric")
	}

	sb := &strings.Builder{}
	sb.WriteString("numeric")
	sb.WriteString("(")
	var args []interface{}
	if c.Precision != nil {
		args = append(args, *c.Precision)
	}
	if c.Scale != nil {
		args = append(args, *c.Scale)
	}
	sb.WriteString(writeParamPlaceholders(len(args)))
	sb.WriteString(")")
	return NewDDL(sb.String(), args...)
}

type SQLDataTypeDoublePrecision struct{}

func (c SQLDataTypeDoublePrecision) DDL() DDL {
	return NewDDL("double precision")
}

type SQLDataTypeSerial struct{}

func (c SQLDataTypeSerial) DDL() DDL {
	return NewDDL("serial")
}

type SQLDataTypeSmallSerial struct{}

func (c SQLDataTypeSmallSerial) DDL() DDL {
	return NewDDL("smallserial")
}

type SQLDataTypeBigSerial struct{}

func (c SQLDataTypeBigSerial) DDL() DDL {
	return NewDDL("bigserial")
}

type SQLDataTypeBoolean struct{}

func (c SQLDataTypeBoolean) DDL() DDL {
	return NewDDL("boolean")
}

type SQLDataTypeFloat struct {
	Precision int
}

func (c SQLDataTypeFloat) DDL() DDL {
	if c.Precision == 0 {
		return NewDDL("float")
	}
	return NewDDL("float(?)", c.Precision)
}

type SQLDataTypeMoney struct{}

func (c SQLDataTypeMoney) DDL() DDL {
	return NewDDL("money")
}

type SQLDataTypeTimestampTZMode string

var (
	TimestampWithTimeZone    SQLDataTypeTimestampTZMode = "WithTimezone"
	TimestampWithoutTimeZone SQLDataTypeTimestampTZMode = "WithoutTimezone"
)

type SQLDataTypeTimestamp struct {
	Timezone *SQLDataTypeTimestampTZMode `json:"timezone,omitempty" yaml:"timezone,omitempty"`
	Digits   *int                        `json:"digits,omitempty" yaml:"digits,omitempty"`
}

func (c SQLDataTypeTimestamp) DDL() DDL {
	sb := strings.Builder{}
	var args []interface{}
	sb.WriteString("timestamp")

	if c.Digits != nil {
		args = append(args, *c.Digits)
		sb.WriteString(" ?")
	}

	if c.Timezone != nil {
		switch *c.Timezone {
		case TimestampWithTimeZone:
			sb.WriteString(" with time zone")
		case TimestampWithoutTimeZone:
			sb.WriteString(" without time zone")
		}
	}
	return NewDDL(sb.String(), args...)
}

type SQLDataTypeDate struct{}

func (c SQLDataTypeDate) DDL() DDL {
	return NewDDL("date")
}

type SQLDataTypeTime struct {
	Timezone *SQLDataTypeTimestampTZMode
	Digits   *int
}

func (c SQLDataTypeTime) DDL() DDL {
	sb := strings.Builder{}
	sb.WriteString("time")
	var args []interface{}

	if c.Digits != nil {
		args = append(args, *c.Digits)
		sb.WriteString(" ?")
	}

	if c.Timezone != nil {
		switch *c.Timezone {
		case TimestampWithTimeZone:
			sb.WriteString(" with time zone")
		case TimestampWithoutTimeZone:
			sb.WriteString(" without time zone")
		}
	}
	return NewDDL(sb.String(), args...)
}

type SQLDataTypeInterval struct {
	Digits *int
	Fields []string
}

func (c SQLDataTypeInterval) DDL() DDL {
	sb := &strings.Builder{}
	var args []interface{}
	sb.WriteString("interval")
	for _, field := range c.Fields {
		sb.WriteString(" ?")
		args = append(args, field)
	}

	if c.Digits != nil {
		sb.WriteString(" ?")
		args = append(args, *c.Digits)
	}

	return NewDDL(sb.String(), args...)
}

type SQLDataTypeArray struct {
	DataType SQLDataType
	Length   uint
}

func (s SQLDataTypeArray) DDL() DDL {
	arrDDL := s.DataType.DDL()
	arrDDL.Query += "["
	if s.Length != 0 {
		arrDDL.Query += strconv.Itoa(int(s.Length))
	}
	arrDDL.Query += "]"
	return arrDDL
}
