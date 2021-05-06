package server

import (
	"context"
	"github.com/nrc-no/core/apps/api/pkg/registry/rest"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/sets"
	"net/http/httptest"
	"testing"
)

var (
	v1GroupVersion = schema.GroupVersion{Group: "", Version: "v1"}
	scheme         = runtime.NewScheme()
	codecs         = serializer.NewCodecFactory(scheme)
	parameterCodec = runtime.NewParameterCodec(scheme)
)

func init() {
	metav1.AddToGroupVersion(scheme, metav1.SchemeGroupVersion)
	scheme.AddUnversionedTypes(v1GroupVersion,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
}

type testGetterStorage struct {
	Version string
}

func (p *testGetterStorage) New() runtime.Object {
	return &metav1.APIGroup{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Getter",
			APIVersion: p.Version,
		},
	}
}

func (p *testGetterStorage) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return &metav1.APIGroup{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Getter",
			APIVersion: p.Version,
		},
		Name: "name",
	}, nil
}

func (p *testGetterStorage) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc) (runtime.Object, error) {
	return nil, nil
}

func TestServer(t *testing.T) {

	c := NewConfig(codecs)
	c.LegacyAPIGroupPrefixes = sets.NewString("/apiPrefix")
	s, err := c.New("test")
	if !assert.NoError(t, err) {
		return
	}

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(v1GroupVersion, &metav1.Status{})

	testApi := func(kind, name string, gv schema.GroupVersion) APIGroupInfo {
		getter := testGetterStorage{}

		scheme.AddKnownTypeWithName(gv.WithKind(kind), getter.New())

		metav1.AddToGroupVersion(scheme, v1GroupVersion)
		return APIGroupInfo{
			PrioritizedVersions: []schema.GroupVersion{gv},
			VersionedResourcesStorageMap: map[string]map[string]rest.Storage{
				gv.Version: {
					name: &testGetterStorage{Version: gv.Version},
				},
			},
			ParameterCodec:       parameterCodec,
			NegotiatedSerializer: codecs,
			Scheme:               scheme,
		}
	}

	apis := []APIGroupInfo{
		testApi("Getter", "getters", schema.GroupVersion{Group: "", Version: "v1"}),
		testApi("Getter", "getters", schema.GroupVersion{Group: "test", Version: "v1"}),
	}

	if !assert.NoError(t, s.InstallLegacyAPIGroup("/apiPrefix", &apis[0])) {
		return
	}
	for _, api := range apis[1:] {
		if !assert.NoError(t, s.InstallAPIGroups(&api)) {
			return
		}
	}

	server := httptest.NewServer(s.Handler)
	defer server.Close()

	api := testApi("Getter", "blabbers", schema.GroupVersion{Group: "test", Version: "v1"})
	if !assert.NoError(t, s.InstallAPIGroups(&api)) {
		return
	}

	// FormDefinitions
	// CustomResources

	req := httptest.NewRequest("GET", "/apis/test/v1/getters", nil)
	w := httptest.NewRecorder()
	s.Handler.ServeHTTP(w, req)
	bytes, err := ioutil.ReadAll(w.Body)
	if !assert.NoError(t, err) {
		return
	}
	t.Log(string(bytes))

	return

}
