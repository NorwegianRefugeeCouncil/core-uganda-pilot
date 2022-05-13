package api

import (
	"encoding/json"
)

type Attributes map[string]Value

func NewAttributes() Attributes {
	return make(Attributes)
}

func (a Attributes) WithString(key string, value string) Attributes {
	a[key] = NewStringValue(value, true)
	return a
}

func (a Attributes) WithInt(key string, value int64) Attributes {
	a[key] = NewIntValue(value, true)
	return a
}

func (a Attributes) WithFloat(key string, value float64) Attributes {
	a[key] = NewFloatValue(value, true)
	return a
}

func (a Attributes) WithBool(key string, value bool) Attributes {
	a[key] = NewBoolValue(value, true)
	return a
}

// Record represents a record in a database
type Record struct {
	ID               string     `json:"id"`
	Table            string     `json:"table"`
	Revision         Revision   `json:"revision"`
	PreviousRevision Revision   `json:"-"`
	Attributes       Attributes `json:"attributes"`
}

func (r *Record) UnmarshalJSON(data []byte) error {
	type record struct {
		ID         string     `json:"id"`
		Table      string     `json:"table"`
		Revision   Revision   `json:"revision"`
		Attributes Attributes `json:"attributes"`
	}
	var rr record
	if err := json.Unmarshal(data, &rr); err != nil {
		return err
	}
	r.ID = rr.ID
	r.Table = rr.Table
	r.Revision = rr.Revision
	r.Attributes = rr.Attributes
	if r.Attributes == nil {
		r.Attributes = make(Attributes)
	}
	return nil
}

// String returns a string representation of the record
func (r Record) String() string {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// SetFieldValue sets the value of a field
func (r Record) SetFieldValue(name string, value Value) Record {
	if r.Attributes == nil {
		r.Attributes = make(map[string]Value)
	}
	r.Attributes[name] = value
	return r
}

// HasField returns true if the record has a field with the given name
func (r Record) HasField(name string) bool {
	_, ok := r.Attributes[name]
	return ok
}

// GetFieldValue returns the value of the field with the given name
// Or an error if the field does not exist
func (r Record) GetFieldValue(name string) (Value, error) {
	if r.Attributes == nil {
		return Value{}, ErrFieldNotFound
	}
	value, ok := r.Attributes[name]
	if !ok {
		return Value{}, ErrFieldNotFound
	}
	return value, nil
}

// GetID returns the ID of the record
// or empty string if the record does not have an ID field
func (r Record) GetID() string {
	return r.ID
}

// GetRevision returns the revision of the record
// or empty string if the record does not have a revision field
func (r Record) GetRevision() Revision {
	return r.Revision
}

// Table represents a database table
type Table struct {
	// Name of the table
	Name string `json:"name"`
	// Columns of the table
	Columns []Column `json:"columns"`
	// Constraints of the table
	Constraints []TableConstraint `json:"constraints"`
}

func (t Table) String() string {
	jsonBytes, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// TableList represents a list of tables
type TableList struct {
	// Items is the list of tables
	Items []Table `json:"items"`
}

// TableConstraint represents a SQL table constraint
type TableConstraint struct {
	PrimaryKey *PrimaryKeyTableConstraint `json:"primary_key"`
}

// PrimaryKeyTableConstraint represents a primary key table constraint
type PrimaryKeyTableConstraint struct {
	// Columns of the primary key
	Columns []string `json:"columns"`
}

// Column represents a database column
type Column struct {
	// Name of the column
	Name string `json:"name"`
	// Type is the data type of the column
	Type string `json:"type"`
	// Default value of the column
	Default string `json:"default"`
	// Constraints of the column
	Constraints []ColumnConstraint `json:"constraints"`
}

// ColumnConstraint represents a SQL column constraint
type ColumnConstraint struct {
	NotNull    *NotNullColumnConstraint    `json:"not_null"`
	PrimaryKey *PrimaryKeyColumnConstraint `json:"primary_key"`
}

// NotNullColumnConstraint represents a not null column constraint
type NotNullColumnConstraint struct{}

// PrimaryKeyColumnConstraint represents a primary key column constraint
type PrimaryKeyColumnConstraint struct{}
