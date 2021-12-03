package store

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestClearProtection(t *testing.T) {
	err := checkClearProtection(context.Background())
	assert.Error(t, err)
	if err := os.Setenv(clearProtectionEnvVar, "yes"); !assert.NoError(t, err) {
		return
	}
	err = checkClearProtection(context.Background())
	assert.NoError(t, err)
}
