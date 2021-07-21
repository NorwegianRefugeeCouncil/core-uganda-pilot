package iam

import (
	"github.com/nrc-no/core/pkg/validation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertRequired(field string) func(t *testing.T, errList validation.ErrorList) {
	return func(t *testing.T, errList validation.ErrorList) {
		assert.NotEmpty(t, errList)
		assert.Equal(t, errList.Find(field)[0].Type, validation.ErrorTypeRequired)
	}
}
func assertInvalid(field string) func(t *testing.T, errList validation.ErrorList) {
	return func(t *testing.T, errList validation.ErrorList) {
		assert.NotEmpty(t, errList)
		assert.Equal(t, errList.Find(field)[0].Type, validation.ErrorTypeInvalid)
	}
}
