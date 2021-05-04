package fake

import (
	"context"
	"github.com/nrc-no/core/apps/api/apis/core"
	"github.com/nrc-no/core/apps/api/pkg/store"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestGet(t *testing.T) {

	s := NewFakeStore()
	var TestObj core.Model
	assert.NoError(t, s.Get(context.TODO(), "hello", store.GetOptions{}, &TestObj))

}
