package validation

import (
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
)

func ValidateFooReads(foo *types.FooReads) validation.ErrorList {
	var result validation.ErrorList

	namePath := validation.NewPath("name")
	otherFieldPath := validation.NewPath("otherField")
	uuidFieldPath := validation.NewPath("uuidField")
	validPath := validation.NewPath("valid")

	if foo.Name == nil || len(*foo.Name) == 0 {
		result = append(result, validation.Required(namePath, "name is required"))
	}

	if foo.OtherField == nil || *foo.OtherField <= 0 || *foo.OtherField >= 100 {
		result = append(result, validation.Required(otherFieldPath, "otherField is required"))
	}

	if foo.UUIDField == nil || *foo.UUIDField == uuid.Nil {
		result = append(result, validation.Required(uuidFieldPath, "uuidField is required"))
	}

	if foo.Valid == nil {
		result = append(result, validation.Required(validPath, "validPath is required"))
	}

	return result
}

func ValidateFoo(foo *types.Foo) validation.ErrorList {
	var result validation.ErrorList

	idPath := validation.NewPath("id")
	namePath := validation.NewPath("name")
	uuidFieldPath := validation.NewPath("uuidField")
	otherFieldPath := validation.NewPath("otherField")
	// validPath := validation.NewPath("valid") // No validation required?

	if foo.ID == uuid.Nil {
		result = append(result, validation.Required(idPath, "id is required"))
	}

	if len(foo.Name) == 0 {
		result = append(result, validation.Required(namePath, "name is required"))
	}

	if foo.UUIDField == uuid.Nil {
		result = append(result, validation.Required(uuidFieldPath, "uuidField is required"))
	}

	if foo.OtherField <= 0 || foo.OtherField >= 100 {
		result = append(result, validation.Required(otherFieldPath, "otherField is out of range"))
	}

	return result
}
