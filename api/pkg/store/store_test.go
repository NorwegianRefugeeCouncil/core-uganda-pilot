package store

import (
	"context"
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"strings"
	"testing"
)

type Suite struct {
	suite.Suite
	ctx          context.Context
	cancel       context.CancelFunc
	client       *mongo.Client
	scheme       *runtime.Scheme
	codecFactory serializer.CodecFactory
	codec        runtime.Codec
	store        *MongoStore
}

func (s *Suite) TestCRUD() {

	fd := &corev1.FormDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "FormDefinition",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewV4().String(),
		},
	}

	// Test creating the form definition
	out := &corev1.FormDefinition{}
	key := getObjectKey(fd)
	if err := s.store.Create(s.ctx, key, CreateOptions{}, fd, out); !assert.NoError(s.T(), err) {
		return
	}

	// Test listing the form definitions
	list := &corev1.FormDefinitionList{}
	if err := s.store.List(s.ctx, getListKey(), ListOptions{}, list); !assert.NoError(s.T(), err) {
		return
	}
	var found = false
	for _, item := range list.Items {
		if item.Name == fd.Name {
			found = true
			break
		}
	}
	assert.True(s.T(), found)

	// Test getting back the form definition
	in := &corev1.FormDefinition{}
	if err := s.store.Get(s.ctx, key, GetOptions{}, in); !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), fd, in)

	// Test updating the form definition
	updated := &corev1.FormDefinition{}
	if err := s.store.Update(s.ctx, key, UpdateOptions{}, updated, func(input runtime.Object) (output runtime.Object, err error) {
		formDef := input.(*corev1.FormDefinition)
		formDef.Spec.Names.Kind = "Updated"
		return formDef, nil
	}); !assert.NoError(s.T(), err) {
		return
	}
	expected := fd.DeepCopy()
	expected.Spec.Names.Kind = "Updated"
	assert.Equal(s.T(), expected, updated)

	// Test deleting the form definition
	deleted := &corev1.FormDefinition{}
	if err := s.store.Delete(s.ctx, key, DeleteOptions{}, deleted); !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), expected, deleted)

}

func getObjectKey(fd *corev1.FormDefinition) string {
	return strings.Join([]string{
		corev1.SchemeGroupVersion.Group,
		"formdefinitions",
		fd.Name,
	}, "/")
}

func getListKey() string {
	return strings.Join([]string{
		corev1.SchemeGroupVersion.Group,
		"formdefinitions",
	}, "/")
}

func (s *Suite) SetupSuite() {

	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancel = cancel

	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetAuth(options.Credential{
			Username: "root",
			Password: "pass12345",
		}))
	if err != nil {
		s.T().Fatal(err)
	}
	s.client = client

	scheme := runtime.NewScheme()
	s.scheme = scheme
	codecFactory := serializer.NewCodecFactory(scheme)
	s.codecFactory = codecFactory

	if err := corev1.AddToScheme(scheme); err != nil {
		s.T().Fatal(err)
	}

	codec := codecFactory.LegacyCodec(corev1.SchemeGroupVersion)
	s.codec = codec

	store := NewMongoStore(client, codec, "test")
	s.store = store

}

func (s *Suite) TearDownSuite() {
	s.cancel()
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
