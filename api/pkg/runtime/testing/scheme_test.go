package testing

import (
	"github.com/nrc-no/core/apps/api/pkg/apis/example"
	v1 "github.com/nrc-no/core/apps/api/pkg/apis/example/v1"
	v2 "github.com/nrc-no/core/apps/api/pkg/apis/example/v2"
	"github.com/nrc-no/core/apps/api/pkg/apis/meta/v1/unstructured"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestRecognizes(t *testing.T) {
	s := runtime.NewScheme()
	if !assert.NoError(t, v1.AddToScheme(s)) {
		return
	}
	if !assert.NoError(t, v2.AddToScheme(s)) {
		return
	}

	assert.False(t, s.Recognizes(v1.SchemeGroupVersion.WithKind("RandomKind")))
	assert.False(t, s.Recognizes(v2.SchemeGroupVersion.WithKind("RandomKind")))
	assert.True(t, s.Recognizes(v1.SchemeGroupVersion.WithKind("TestModel")))
	assert.True(t, s.Recognizes(v2.SchemeGroupVersion.WithKind("TestModel")))
}

func TestConvertSelf(t *testing.T) {
	s := runtime.NewScheme()
	if !assert.NoError(t, v1.AddToScheme(s)) {
		return
	}

	a := &v1.TestModel{Spec: v1.TestModelSpec{SomeProperty: "abc"}}
	b := &v1.TestModel{}

	assert.NoError(t, s.Convert(a, b, nil))
	assert.Equal(t, "abc", b.Spec.SomeProperty)
}

func TestConvertToUnstructured(t *testing.T) {
	s := runtime.NewScheme()
	if !assert.NoError(t, v1.AddToScheme(s)) {
		return
	}

	a := &v1.TestModel{Spec: v1.TestModelSpec{SomeProperty: "abc"}}
	b := &unstructured.Unstructured{}

	assert.NoError(t, s.Convert(a, b, nil))
	assert.Equal(t, "abc", (b.Object["spec"].(map[string]interface{}))["someProperty"].(string))
}

func TestConvertBothToUnstructured(t *testing.T) {
	s := runtime.NewScheme()
	if !assert.NoError(t, v1.AddToScheme(s)) {
		return
	}

	a := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": v1.SchemeGroupVersion.String(),
			"kind":       "TestModel",
			"spec": map[string]interface{}{
				"someProperty": "abc",
			},
		},
	}

	b := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": v1.SchemeGroupVersion.String(),
			"kind":       "TestModel2",
		},
	}

	assert.NoError(t, s.Convert(a, b, nil))
	assert.Equal(t, "abc", (b.Object["spec"].(map[string]interface{}))["someProperty"].(string))
}

func TestConvertFromUnstructured(t *testing.T) {
	s := runtime.NewScheme()
	if !assert.NoError(t, v1.AddToScheme(s)) {
		return
	}

	s.AddKnownTypes(v1.SchemeGroupVersion, &v1.TestModel{})
	a := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": v1.SchemeGroupVersion.String(),
			"kind":       "TestModel",
			"spec": map[string]interface{}{
				"someProperty": "abc",
			},
		},
	}
	b := &example.TestModel{}
	assert.NoError(t, s.Convert(a, b, nil))
	assert.Equal(t, "abc", b.Spec.SomeProperty)
}

func TestConvertVersion(t *testing.T) {
	s := runtime.NewScheme()
	if !assert.NoError(t, v1.AddToScheme(s)) {
		return
	}
	if !assert.NoError(t, v2.AddToScheme(s)) {
		return
	}

	v1 := &v1.TestModel{Spec: v1.TestModelSpec{SomeProperty: "abc"}}
	v2 := &example.TestModel{}

	assert.NoError(t, s.Convert(v1, v2, nil))
	assert.Equal(t, "abc", v2.Spec.SomeProperty)

}

func TestConvertToVersion(t *testing.T) {
	s := runtime.NewScheme()
	if !assert.NoError(t, v1.AddToScheme(s)) {
		return
	}
	if !assert.NoError(t, example.AddToScheme(s)) {
		return
	}

	v1 := &example.TestModel{Spec: example.TestModelSpec{SomeProperty: "abc"}}
	out, err := s.ConvertToVersion(v1, example.SchemeGroupVersion)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.IsType(t, &example.TestModel{}, out) {
		return
	}
	converted := out.(*example.TestModel)
	assert.Equal(t, "abc", converted.Spec.SomeProperty)
}

func TestConvertUnstructuredToVersion(t *testing.T) {
	s := runtime.NewScheme()
	if !assert.NoError(t, v1.AddToScheme(s)) {
		return
	}
	if !assert.NoError(t, example.AddToScheme(s)) {
		return
	}

	u := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": v1.SchemeGroupVersion.String(),
			"kind":       "TestModel",
			"spec": map[string]interface{}{
				"someProperty": "abc",
			},
		},
	}

	out, err := s.ConvertToVersion(u, example.SchemeGroupVersion)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.IsType(t, &example.TestModel{}, out) {
		return
	}
	converted := out.(*example.TestModel)
	assert.Equal(t, "abc", converted.Spec.SomeProperty)
}

func TestConvertURLValuesToModel(t *testing.T) {
	s := runtime.NewScheme()
	if !assert.NoError(t, example.AddToScheme(s)) {
		return
	}

	codec := runtime.NewParameterCodec(s)
	values := url.Values{}
	values.Add("abc", "def")
	var out = &example.TestModelUrlValues{}
	if !assert.NoError(t, codec.DecodeParameters(values, example.SchemeGroupVersion, out)) {
		return
	}
	assert.Equal(t, "def", out.Abc)
}
