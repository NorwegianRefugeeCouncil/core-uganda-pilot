package test

import "github.com/stretchr/testify/assert"

var ErrorIs = func(expect error) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, err error, i ...interface{}) bool {
		return assert.ErrorIs(t, err, expect, i...)
	}
}
