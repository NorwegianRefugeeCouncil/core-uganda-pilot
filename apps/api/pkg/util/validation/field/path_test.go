package field

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleToString(t *testing.T) {
	assert.Equal(t, "name", NewPath("name").String())
}

func TestChildToString(t *testing.T) {
	assert.Equal(t, "name.bla", NewPath("name").Child("bla").String())
}

func TestIndexedToString(t *testing.T) {
	assert.Equal(t, "name[2]", NewPath("name").Index(2).String())
}

func TestKeyedToString(t *testing.T) {
	assert.Equal(t, "name[abc]", NewPath("name").Key("abc").String())
}

func TestManyToString(t *testing.T) {
	assert.Equal(t, "name.a.b.c", NewPath("name", "a", "b", "c").String())
}
