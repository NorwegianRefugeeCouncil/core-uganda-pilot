package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateAccessors(t *testing.T) {
	for _, kind := range GetAllFieldKinds() {
		t.Run(kind.String(), func(t *testing.T) {
			k := kind
			_, ok := FieldAccessors[k]
			assert.True(t, ok, "no accessor for field %s is defined in types.FieldAccessors", kind)
		})
	}
}
