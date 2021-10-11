package iam

import (
	"github.com/nrc-no/core/internal/validation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertRequired(field string) func(t *testing.T, errList validation.ErrorList) {
	return func(t *testing.T, errList validation.ErrorList) {
		assert.NotEmpty(t, errList)
		e := *errList.Find(field)
		assert.Equal(t, e[0].Type, validation.ErrorTypeRequired)
	}
}
func assertInvalid(field string) func(t *testing.T, errList validation.ErrorList) {
	return func(t *testing.T, errList validation.ErrorList) {
		e := *errList.Find(field)
		assert.Equal(t, e[0].Type, validation.ErrorTypeInvalid)
	}
}

func assertNoError(field string) func(t *testing.T, errList validation.ErrorList) {
	return func(t *testing.T, errList validation.ErrorList) {
		e := *errList.Find(field)
		assert.Empty(t, e)
	}
}
