package v1

import (
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/api/pkg/client/rest"
	scheme2 "github.com/nrc-no/core/api/pkg/client/typed/scheme"
)

type CoreV1Interface interface {
	CustomResourceDefinitionsGetter
	FormDefinitionsGetter
}

// CoreV1Client is used to interact with features provided by the core.nrc.no group.
type CoreV1Client struct {
	restClient rest.Interface
}

func (c *CoreV1Client) CustomResourceDefinitions() CustomResourceDefinitionInterface {
	return newCustomResourceDefinitions(c)
}

func (c *CoreV1Client) FormDefinitions() FormDefinitionInterface {
	return newFormDefinitions(c)
}

// NewForConfig creates a new CoreV1Client for the given config.
func NewForConfig(c *rest.Config) (*CoreV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &CoreV1Client{client}, nil
}

// NewForConfigOrDie creates a new CoreV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *CoreV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new CoreV1Client for the given RESTClient.
func New(c rest.Interface) *CoreV1Client {
	return &CoreV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := corev1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme2.Codecs.WithoutConversion()
	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *CoreV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
