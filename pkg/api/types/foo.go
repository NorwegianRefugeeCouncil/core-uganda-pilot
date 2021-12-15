package types

import (
	uuid "github.com/satori/go.uuid"
)

type FooReads struct {
	Name       *string    `json:"name"`
	OtherField *int       `json:"otherField"`
	UUIDField  *uuid.UUID `json:"uuidField"`
	Valid      *bool      `json:"valid"`
}

type Foo struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	OtherField int       `json:"otherField"`
	UUIDField  uuid.UUID `json:"uuidField"`
	Valid      bool      `json:"valid"`
}
