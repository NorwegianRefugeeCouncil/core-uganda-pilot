package eventstore

import (
	"context"
	goes "github.com/jetbasrawi/go.geteventstore"
	"github.com/nrc-no/core/apps/api/pkg/store"
	testing2 "github.com/nrc-no/core/apps/api/pkg/store/eventstore/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"net/http"
	"testing"
)

var (
	gv = schema.GroupVersion{Group: "test", Version: "v1"}
)

type StoreTestSuite struct {
	suite.Suite
	store  *Store
	scheme *runtime.Scheme
}

func (s *StoreTestSuite) SetupSuite() {
	s.scheme = runtime.NewScheme()
	codec := json.NewSerializer(json.DefaultMetaFactory, s.scheme, s.scheme, false)
	utilruntime.Must(metav1.AddMetaToScheme(s.scheme))
	s.scheme.AddKnownTypeWithName(gv.WithKind("TestObj"), &testing2.TestObj{})

	c, err := goes.NewClient(http.DefaultClient, "http://localhost:2113")
	if !s.Assert().NoError(err) {
		return
	}
	s.store = NewStore(c, codec)
}

func (s *StoreTestSuite) TestGet() {
	ctx := context.TODO()

	var obj = testing2.TestObj{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TestObj",
			APIVersion: "test/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "abc",
		},
	}
	var out = testing2.TestObj{}
	if !assert.NoError(s.T(), s.store.Create(ctx, "abc", &obj, &out)) {
		return
	}

	if !assert.NoError(s.T(), s.store.Get(ctx, "abc", store.GetOptions{}, &obj)) {
		return
	}

	s.T().Logf("%#v", obj)

}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}
