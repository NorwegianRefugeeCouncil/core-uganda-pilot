package fake

import (
  "context"
  "github.com/nrc-no/core/apps/api/apis/core"
  "github.com/nrc-no/core/apps/api/apis/exceptions"
  "github.com/nrc-no/core/apps/api/pkg/store"
  "github.com/stretchr/testify/assert"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/apimachinery/pkg/runtime"

  "testing"
)

func TestGetNotFound(t *testing.T) {
  ctx := context.TODO()
  s := NewFakeStore()
  var TestObj core.Model
  assert.ErrorIs(t, exceptions.ErrNotFound, s.Get(ctx, "hello", store.GetOptions{}, &TestObj))
}

func TestGet(t *testing.T) {
  ctx := context.TODO()
  s := NewFakeStore()
  obj := core.Model{}
  objPtr := core.Model{}
  assert.NoError(t, s.Create(ctx, "key", &obj, &obj))
  assert.NoError(t, s.Get(ctx, "key", store.GetOptions{}, &objPtr))
}

func TestGetVersion(t *testing.T) {
  ctx := context.TODO()
  s := NewFakeStore()
  obj := core.Model{}
  objPtr := core.Model{}
  assert.NoError(t, s.Create(ctx, "key", &obj, &obj))
  assert.NoError(t, s.Get(ctx, "key", store.GetOptions{}, &objPtr))
}

func TestUpdate(t *testing.T) {
  ctx := context.TODO()
  s := NewFakeStore()
  obj := core.Model{
    ObjectMeta: metav1.ObjectMeta{
      Name: "hello!",
    },
  }

  // Create object
  if !assert.NoError(t, s.Create(ctx, "key", &obj, &obj)){
    return
  }

  // Update object
  if !assert.NoError(t, s.Update(ctx, "key", &core.Model{}, func(obj runtime.Object) (runtime.Object, error) {
    var model = obj.(*core.Model)
    model.ObjectMeta.Name = "bla"
    return model, nil
  }, store.UpdateOptions{})) {
    return
  }

  // Verify
  var updated core.Model
  if !assert.NoError(t, s.Get(ctx, "key", store.GetOptions{}, &updated)) {
    return
  }

  assert.Equal(t, "bla", updated.ObjectMeta.Name)

}
