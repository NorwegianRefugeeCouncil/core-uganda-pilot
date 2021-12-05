package pointers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInt(t *testing.T) {
	var val = 4
	assert.Equal(t, &val, Int(4))
}
func TestInt64(t *testing.T) {
	var val int64 = 4
	assert.Equal(t, &val, Int64(4))
}
func TestInt32(t *testing.T) {
	var val int32 = 4
	assert.Equal(t, &val, Int32(4))
}
func TestStr(t *testing.T) {
	var val = "abc"
	assert.Equal(t, &val, String("abc"))
}
