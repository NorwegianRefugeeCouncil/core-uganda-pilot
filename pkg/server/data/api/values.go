package api

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

var nullBytes = []byte("null")

// String is a nullable string.
// It does not consider empty strings to be null.
// It will decode to null, not "" when null.
// It implements json.Marshaler and json.Unmarshaler.
// It also implements sql.Scanner and sql.Valuer to marshal and unmarshal itself.
// So it is both database and json compatible.
type String struct {
	sql.NullString
}

// StringFrom creates a new String that will always be non-null.
func StringFrom(s string) String {
	return NewString(s, true)
}

// StringFromPtr creates a new String that be null if s is nil.
func StringFromPtr(s *string) String {
	if s == nil {
		return NewString("", false)
	}
	return NewString(*s, true)
}

// ValueOrZero returns the inner value if valid, otherwise empty string
func (ns String) ValueOrZero() string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

// UnmarshalJSON implements json.Unmarshaler.
func (ns *String) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	if err := json.Unmarshal(data, &ns.String); err != nil {
		return err
	}
	return nil
}

// MarshalJSON implements json.Marshaler.
func (ns String) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// NewString creates a new String
func NewString(s string, valid bool) String {
	return String{
		NullString: sql.NullString{
			String: s,
			Valid:  valid,
		},
	}
}

// MarshalText implements encoding.TextMarshaler.
func (ns String) MarshalText() ([]byte, error) {
	if ns.Valid {
		return []byte(ns.String), nil
	}
	return nil, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (ns *String) UnmarshalText(text []byte) error {
	ns.String = string(text)
	ns.Valid = true
	return nil
}

// SetValid changes this String's value and also sets it to be non-null.
func (ns *String) SetValid(v string) {
	ns.String = v
	ns.Valid = true
}

// IsZero returns true for invalid Strings, for Go omitempty tag support.
func (ns *String) IsZero() bool {
	return !ns.Valid
}

// Equal returns true if the other String is equal to this one.
func (ns String) Equal(other String) bool {
	if !ns.Valid && !other.Valid {
		return true
	}
	if !ns.Valid || !other.Valid {
		return false
	}
	return ns.String == other.String
}

// Ptr returns a pointer to this String's value, or a nil pointer if this String is invalid.
func (ns String) Ptr() *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

// Bool is a nullable bool.
// It does not default to false
// It will decode to null, not false when null.
// It implements json.Marshaler and json.Unmarshaler.
// It also implements sql.Scanner and sql.Valuer to marshal and unmarshal itself.
// So it is both database and json compatible.
type Bool struct {
	sql.NullBool
}

// BoolFrom creates a new Bool that will always be non-null.
func BoolFrom(s bool) Bool {
	return NewBool(s, true)
}

// BoolFromPtr creates a new Bool that be null if s is nil.
func BoolFromPtr(b *bool) Bool {
	if b == nil {
		return NewBool(false, false)
	}
	return NewBool(*b, true)
}

// ValueOrZero returns the inner value if valid, otherwise false
func (b Bool) ValueOrZero() bool {
	if !b.Valid {
		return false
	}
	return b.Bool
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Bool) UnmarshalJSON(data []byte) error {
	str := string(data)
	switch str {
	case "", "null":
		b.Valid = false
		return nil
	case "true":
		b.Valid = true
		b.Bool = true
		return nil
	case "false":
		b.Valid = true
		b.Bool = false
		return nil
	default:
		return fmt.Errorf("invalid boolean value: %s", str)
	}
}

// MarshalJSON implements json.Marshaler.
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	if b.Bool {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

// NewBool creates a new Bool
func NewBool(b bool, valid bool) Bool {
	return Bool{
		NullBool: sql.NullBool{
			Bool:  b,
			Valid: valid,
		},
	}
}

// MarshalText implements encoding.TextMarshaler.
func (b Bool) MarshalText() ([]byte, error) {
	if !b.Valid {
		return []byte{}, nil
	}
	if b.Bool {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Bool) UnmarshalText(text []byte) error {
	str := string(text)
	switch str {
	case "", "null":
		b.Valid = false
		return nil
	case "true":
		b.Valid = true
		b.Bool = true
		return nil
	case "false":
		b.Valid = true
		b.Bool = false
		return nil
	default:
		return fmt.Errorf("invalid boolean value: %s", str)
	}
}

// SetValid changes this Bool's value and also sets it to be non-null.
func (b *Bool) SetValid(v bool) {
	b.Bool = v
	b.Valid = true
}

// IsZero returns true for invalid Bools, for Go omitempty tag support.
func (b *Bool) IsZero() bool {
	return !b.Valid
}

// Equal returns true if the other Bool is equal to this one.
func (b Bool) Equal(other Bool) bool {
	if !b.Valid && !other.Valid {
		return true
	}
	if !b.Valid || !other.Valid {
		return false
	}
	return b.Bool == other.Bool
}

// Ptr returns a pointer to this Bool's value, or a nil pointer if this Bool is invalid.
func (b Bool) Ptr() *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

// Int is a nullable int.
// It does not default to 0
// It will decode to null, not 0 when null.
// It implements json.Marshaler and json.Unmarshaler.
// It also implements sql.Scanner and sql.Valuer to marshal and unmarshal itself.
// So it is both database and json compatible.
type Int struct {
	sql.NullInt64
}

// IntFrom creates a new Int that will always be non-null.
func IntFrom(s int64) Int {
	return NewInt(s, true)
}

// IntFromPtr creates a new Int that be null if s is nil.
func IntFromPtr(b *int64) Int {
	if b == nil {
		return NewInt(0, false)
	}
	return NewInt(*b, true)
}

// ValueOrZero returns the inner value if valid, otherwise false
func (b Int) ValueOrZero() int64 {
	if !b.Valid {
		return 0
	}
	return b.Int64
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Int) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		b.Valid = false
		return nil
	}
	if err := json.Unmarshal(data, &b.Int64); err != nil {
		var typeError *json.UnmarshalTypeError
		if errors.As(err, &typeError) {
			if typeError.Value != "string" {
				return err
			}
			var str string
			if err := json.Unmarshal(data, &str); err != nil {
				return err
			}
			b.Int64, err = strconv.ParseInt(str, 10, 64)
			if err != nil {
				return err
			}
			b.Valid = true
			return nil
		}
		return err
	}
	b.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (b Int) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(b.Int64, 10)), nil
}

// NewInt creates a new Int
func NewInt(i int64, valid bool) Int {
	return Int{
		NullInt64: sql.NullInt64{
			Int64: i,
			Valid: valid,
		},
	}
}

// MarshalText implements encoding.TextMarshaler.
func (b Int) MarshalText() ([]byte, error) {
	if !b.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(b.Int64, 10)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Int) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		b.Valid = false
		return nil
	}
	var err error
	b.Int64, err = strconv.ParseInt(string(text), 10, 64)
	if err != nil {
		return err
	}
	b.Valid = true
	return nil
}

// SetValid changes this Int's value and also sets it to be non-null.
func (b *Int) SetValid(v int64) {
	b.Int64 = v
	b.Valid = true
}

// IsZero returns true for invalid Ints, for Go omitempty tag support.
func (b *Int) IsZero() bool {
	return !b.Valid
}

// Equal returns true if the other Int is equal to this one.
func (b Int) Equal(other Int) bool {
	if !b.Valid && !other.Valid {
		return true
	}
	if !b.Valid || !other.Valid {
		return false
	}
	return b.Int64 == other.Int64
}

// Ptr returns a point64er to this Int's value, or a nil point64er if this Int is invalid.
func (b Int) Ptr() *int64 {
	if !b.Valid {
		return nil
	}
	return &b.Int64
}

// Float is a nullable float.
// It does not default to 0
// It will decode to null, not 0 when null.
// It implements json.Marshaler and json.Unmarshaler.
// It also implements sql.Scanner and sql.Valuer to marshal and unmarshal itself.
// So it is both database and json compatible.
type Float struct {
	sql.NullFloat64
}

// FloatFrom creates a new Float that will always be non-null.
func FloatFrom(s float64) Float {
	return NewFloat(s, true)
}

// FloatFromPtr creates a new Float that be null if s is nil.
func FloatFromPtr(b *float64) Float {
	if b == nil {
		return NewFloat(0, false)
	}
	return NewFloat(*b, true)
}

// ValueOrZero returns the inner value if valid, otherwise false
func (b Float) ValueOrZero() float64 {
	if !b.Valid {
		return 0
	}
	return b.Float64
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Float) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		b.Valid = false
		return nil
	}
	if err := json.Unmarshal(data, &b.Float64); err != nil {
		var typeError *json.UnmarshalTypeError
		if errors.As(err, &typeError) {
			if typeError.Value != "string" {
				return err
			}
			var str string
			if err := json.Unmarshal(data, &str); err != nil {
				return err
			}
			b.Float64, err = strconv.ParseFloat(str, 64)
			if err != nil {
				return err
			}
			b.Valid = true
			return nil
		}
		return err
	}
	b.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (b Float) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	if math.IsInf(b.Float64, 0) || math.IsNaN(b.Float64) {
		return nil, &json.UnsupportedValueError{
			Value: reflect.ValueOf(b.Float64),
			Str:   strconv.FormatFloat(b.Float64, 'g', -1, 64),
		}
	}
	return []byte(strconv.FormatFloat(b.Float64, 'f', -1, 64)), nil
}

// NewFloat creates a new Float
func NewFloat(i float64, valid bool) Float {
	return Float{
		NullFloat64: sql.NullFloat64{
			Float64: i,
			Valid:   valid,
		},
	}
}

// MarshalText implements encoding.TextMarshaler.
func (b Float) MarshalText() ([]byte, error) {
	if !b.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatFloat(b.Float64, 'f', -1, 64)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Float) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		b.Valid = false
		return nil
	}
	var err error
	b.Float64, err = strconv.ParseFloat(string(text), 64)
	if err != nil {
		return err
	}
	b.Valid = true
	return nil
}

// SetValid changes this Float's value and also sets it to be non-null.
func (b *Float) SetValid(v float64) {
	b.Float64 = v
	b.Valid = true
}

// IsZero returns true for invalid Floats, for Go omitempty tag support.
func (b *Float) IsZero() bool {
	return !b.Valid
}

// Equal returns true if the other Float is equal to this one.
func (b Float) Equal(other Float) bool {
	if !b.Valid && !other.Valid {
		return true
	}
	if !b.Valid || !other.Valid {
		return false
	}
	return b.Float64 == other.Float64
}

// Ptr returns a pofloat64er to this Float's value, or a nil pofloat64er if this Float is invalid.
func (b Float) Ptr() *float64 {
	if !b.Valid {
		return nil
	}
	return &b.Float64
}

// ValueKind represents the kind of value stored in a Value.
type ValueKind uint8

//go:generate go run golang.org/x/tools/cmd/stringer -type=ValueKind

const (
	// ValueKindNull represents a null value.
	ValueKindNull ValueKind = iota
	// ValueKindInt is a ValueKind representing an int64.
	ValueKindInt
	// ValueKindFloat is a ValueKind representing a float64.
	ValueKindFloat
	// ValueKindString is a ValueKind representing a string.
	ValueKindString
	// ValueKindBool is a ValueKind representing a bool.
	ValueKindBool
)

// Value represents a value of any kind.
type Value struct {
	// Kind is the kind of value this is.
	Kind ValueKind `json:"-"`
	// String is the string value if Kind is ValueKindString.
	String *String `json:"string,omitempty"`
	// Bool is the bool value if Kind is ValueKindBool.
	Bool *Bool `json:"bool,omitempty"`
	// Int is the int64 value if Kind is ValueKindInt.
	Int *Int `json:"int,omitempty"`
	// Float is the float64 value if Kind is ValueKindFloat.
	Float *Float `json:"float,omitempty"`
}

// Scan implements the sql.Scanner interface.
func (v *Value) Scan(value interface{}) error {
	switch v.Kind {
	case ValueKindInt:
		v.Int = &Int{}
		return v.Int.Scan(value)
	case ValueKindFloat:
		v.Float = &Float{}
		return v.Float.Scan(value)
	case ValueKindString:
		v.String = &String{}
		return v.String.Scan(value)
	case ValueKindBool:
		v.Bool = &Bool{}
		return v.Bool.Scan(value)
	default:
		return fmt.Errorf("unknown ValueKind %d", v.Kind)
	}
}

// Value implements the driver.Valuer interface.
func (v *Value) Value() (driver.Value, error) {
	switch v.Kind {
	case ValueKindNull:
		return nil, nil
	case ValueKindInt:
		return v.Int.Value()
	case ValueKindFloat:
		return v.Float.Value()
	case ValueKindString:
		return v.String.Value()
	case ValueKindBool:
		return v.Bool.Value()
	default:
		return nil, fmt.Errorf("unknown ValueKind %d", v.Kind)
	}
}

func (v *Value) UnmarshalJSON(data []byte) error {
	v.String = nil
	v.Bool = nil
	v.Int = nil
	v.Float = nil
	type value struct {
		Null   bool   `json:"null"`
		String string `json:"string,omitempty"`
		Int    string `json:"int,omitempty"`
		Float  string `json:"float,omitempty"`
		Bool   string `json:"bool,omitempty"`
	}
	var val value
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	if val.Null {
		v.Kind = ValueKindNull
		return nil
	}
	if val.String != "" {
		v.Kind = ValueKindString
		v.String = &String{}
		return v.String.UnmarshalText([]byte(val.String))
	}
	if val.Int != "" {
		v.Kind = ValueKindInt
		v.Int = &Int{}
		return v.Int.UnmarshalText([]byte(val.Int))
	}
	if val.Float != "" {
		v.Kind = ValueKindFloat
		v.Float = &Float{}
		return v.Float.UnmarshalText([]byte(val.Float))
	}
	if val.Bool != "" {
		v.Kind = ValueKindBool
		v.Bool = &Bool{}
		return v.Bool.UnmarshalText([]byte(val.Bool))
	}
	return fmt.Errorf("invalid value: %s", data)
}

// MarshalJSON implements json.Marshaler.
func (v Value) MarshalJSON() ([]byte, error) {
	var ret = map[string]interface{}{}
	switch v.Kind {
	case ValueKindNull:
		return marshalNullValue(ret)
	case ValueKindInt:
		return marshalIntValue(v, ret)
	case ValueKindFloat:
		return marshalFloatValue(v, ret)
	case ValueKindString:
		return marshalStringValue(v, ret)
	case ValueKindBool:
		return marshalBoolValue(v, ret)
	}
	return nil, fmt.Errorf("unsupported value kind %d", v.Kind)
}

func marshalNullValue(ret map[string]interface{}) ([]byte, error) {
	ret["null"] = true
	return json.Marshal(ret)
}

func marshalBoolValue(v Value, ret map[string]interface{}) ([]byte, error) {
	if v.Bool.IsZero() {
		ret["null"] = true
	} else {
		text, err := v.Bool.MarshalText()
		if err != nil {
			return nil, err
		}
		ret["bool"] = string(text)
	}
	return json.Marshal(ret)
}

func marshalStringValue(v Value, ret map[string]interface{}) ([]byte, error) {
	if v.String.IsZero() {
		ret["null"] = true
	} else {
		text, err := v.String.MarshalText()
		if err != nil {
			return nil, err
		}
		ret["string"] = string(text)
	}
	return json.Marshal(ret)
}

func marshalFloatValue(v Value, ret map[string]interface{}) ([]byte, error) {
	if v.Float.IsZero() {
		ret["null"] = true
	} else {
		text, err := v.Float.MarshalText()
		if err != nil {
			return nil, err
		}
		ret["float"] = string(text)
	}
	return json.Marshal(ret)
}

func marshalIntValue(v Value, ret map[string]interface{}) ([]byte, error) {
	if v.Int.IsZero() {
		ret["null"] = true
	} else {
		text, err := v.Int.MarshalText()
		if err != nil {
			return nil, err
		}
		ret["int"] = string(text)
	}
	return json.Marshal(ret)
}

// NewStringValue creates a new Value with a string value.
func NewStringValue(s string, valid bool) Value {
	val := NewString(s, valid)
	return Value{
		Kind:   ValueKindString,
		String: &val,
	}
}

// NewBoolValue creates a new Value with a bool value.
func NewBoolValue(b bool, valid bool) Value {
	val := NewBool(b, valid)
	return Value{
		Kind: ValueKindBool,
		Bool: &val,
	}
}

// NewIntValue creates a new Value with an int64 value.
func NewIntValue(i int64, valid bool) Value {
	val := NewInt(i, valid)
	return Value{
		Kind: ValueKindInt,
		Int:  &val,
	}
}

// NewFloatValue creates a new Value with a float64 value.
func NewFloatValue(f float64, valid bool) Value {
	val := NewFloat(f, valid)
	return Value{
		Kind:  ValueKindFloat,
		Float: &val,
	}
}

// NewNullValue creates a new null Value
func NewNullValue() Value {
	return Value{
		Kind: ValueKindNull,
	}
}
