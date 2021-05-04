package exceptions

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestException(t *testing.T) {
	var err = NewAPIError("FakeErr", http.StatusNotFound)
	fmt.Printf("%v", err)
}

func TestExceptionIs(t *testing.T) {
	assert.True(t, errors.Is(ErrNotFound.WithError(fmt.Errorf("hello!")), ErrNotFound))
	assert.False(t, errors.Is(ErrConflict.WithError(fmt.Errorf("hello!")), ErrNotFound))
}
