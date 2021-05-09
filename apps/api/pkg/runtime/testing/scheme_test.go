package testing

import (
	"github.com/nrc-no/core/apps/api/pkg/apis/meta/v1/unstructured"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/testing/testscheme"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

var (
	gvk1_v1 = schema.GroupVersion{
		Group:   "group1",
		Version: "v1",
	}
	gvk1_v2 = schema.GroupVersion{
		Group:   "group1",
		Version: "v2",
	}
	gvk2_v1 = schema.GroupVersion{
		Group:   "group2",
		Version: "v1",
	}
)

func TestRecognizes(t *testing.T) {
	s := runtime.NewScheme()
	s.AddKnownTypes(gvk1_v1, &testscheme.TestModel{})
	assert.False(t, s.Recognizes(gvk1_v1.WithKind("RandomKind")))
	assert.False(t, s.Recognizes(gvk1_v2.WithKind("TestModel")))
	assert.True(t, s.Recognizes(gvk1_v1.WithKind("TestModel")))
}

func TestConvertSelf(t *testing.T) {
	s := runtime.NewScheme()
	s.AddKnownTypes(gvk1_v1, &testscheme.TestModel{})

	a := &testscheme.TestModel{Spec: testscheme.TestModelSpec{SomeProperty: "abc"}}
	b := &testscheme.TestModel{}

	assert.NoError(t, s.Convert(a, b, nil))
	assert.Equal(t, "abc", b.Spec.SomeProperty)
}

func TestConvertToUnstructured(t *testing.T) {
	s := runtime.NewScheme()
	s.AddKnownTypes(gvk1_v1, &testscheme.TestModel{})

	a := &testscheme.TestModel{Spec: testscheme.TestModelSpec{SomeProperty: "abc"}}
	b := &unstructured.Unstructured{}

	assert.NoError(t, s.Convert(a, b, nil))
	assert.Equal(t, "abc", (b.Object["spec"].(map[string]interface{}))["someProperty"].(string))
}

func TestConvertBothToUnstructured(t *testing.T) {
	s := runtime.NewScheme()

	a := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": gvk1_v1.String(),
			"kind":       "TestModel",
			"spec": map[string]interface{}{
				"someProperty": "abc",
			},
		},
	}

	b := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": gvk1_v1.String(),
			"kind":       "TestModel2",
		},
	}

	if !assert.NoError(t, testscheme.AddConversionFunc(s)) {
		return
	}

	assert.NoError(t, s.Convert(a, b, nil))
	assert.Equal(t, "abc", (b.Object["spec"].(map[string]interface{}))["someProperty"].(string))
}

func TestConvertFromUnstructured(t *testing.T) {
	s := runtime.NewScheme()
	s.AddKnownTypes(gvk1_v1, &testscheme.TestModel{})
	a := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": gvk1_v1.String(),
			"kind":       "TestModel",
			"spec": map[string]interface{}{
				"someProperty": "abc",
			},
		},
	}
	b := &testscheme.TestModel{}
	assert.NoError(t, s.Convert(a, b, nil))
	assert.Equal(t, "abc", b.Spec.SomeProperty)
}

func TestConvertVersion(t *testing.T) {
	s := runtime.NewScheme()
	s.AddKnownTypeWithName(gvk1_v1.WithKind("TestModel"), &testscheme.TestModel{})
	s.AddKnownTypeWithName(gvk1_v2.WithKind("TestModel"), &testscheme.TestModel2{})

	v1 := &testscheme.TestModel{Spec: testscheme.TestModelSpec{SomeProperty: "abc"}}
	v2 := &testscheme.TestModel2{}
	if !assert.NoError(t, testscheme.AddConversionFunc(s)) {
		return
	}

	assert.NoError(t, s.Convert(v1, v2, nil))
	assert.Equal(t, "abc", v2.Spec.SomeProperty)

}

func TestConvertToVersion(t *testing.T) {
	s := runtime.NewScheme()
	s.AddKnownTypeWithName(gvk1_v1.WithKind("TestModel"), &testscheme.TestModel{})
	s.AddKnownTypeWithName(gvk1_v2.WithKind("TestModel"), &testscheme.TestModel2{})
	if !assert.NoError(t, testscheme.AddConversionFunc(s)) {
		return
	}
	v1 := &testscheme.TestModel{Spec: testscheme.TestModelSpec{SomeProperty: "abc"}}
	out, err := s.ConvertToVersion(v1, gvk1_v2)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.IsType(t, &testscheme.TestModel2{}, out) {
		return
	}
	converted := out.(*testscheme.TestModel2)
	assert.Equal(t, "abc", converted.Spec.SomeProperty)
}

func TestConvertUnstructuredToVersion(t *testing.T) {
	s := runtime.NewScheme()
	s.AddKnownTypeWithName(gvk1_v1.WithKind("TestModel"), &testscheme.TestModel{})
	s.AddKnownTypeWithName(gvk1_v2.WithKind("TestModel"), &testscheme.TestModel2{})
	if !assert.NoError(t, testscheme.AddConversionFunc(s)) {
		return
	}

	u := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": gvk1_v1.String(),
			"kind":       "TestModel",
			"spec": map[string]interface{}{
				"someProperty": "abc",
			},
		},
	}

	out, err := s.ConvertToVersion(u, gvk1_v2)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.IsType(t, &testscheme.TestModel2{}, out) {
		return
	}
	converted := out.(*testscheme.TestModel2)
	assert.Equal(t, "abc", converted.Spec.SomeProperty)
}

func TestConvertURLValuesToModel(t *testing.T) {
	s := runtime.NewScheme()
	if !assert.NoError(t, testscheme.AddToScheme(s)) {
		return
	}
	codec := runtime.NewParameterCodec(s)
	values := url.Values{}
	values.Add("abc", "def")
	var out = &testscheme.TestModelUrlValues{}
	if !assert.NoError(t, codec.DecodeParameters(values, testscheme.SchemeGroupVersion, out)) {
		return
	}
	assert.Equal(t, "def", out.Abc)
}
