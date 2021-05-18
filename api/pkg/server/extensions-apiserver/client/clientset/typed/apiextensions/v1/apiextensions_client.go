package v1

import (
	v1 "github.com/nrc-no/core/apps/api/pkg/apis/apiextensions/v1"
	"github.com/nrc-no/core/apps/api/pkg/client/rest"
	"github.com/nrc-no/core/apps/api/pkg/server/extensions-apiserver/client/clientset/scheme"
)

type ApiextensionsV1Interface interface {
	RESTClient() rest.Interface
	CustomResourceDefinitionsGetter
}

// ApiextensionsV1Client is used to interact with features provided by the apiextensions.k8s.io group.
type ApiextensionsV1Client struct {
	restClient rest.Interface
}

func (c *ApiextensionsV1Client) CustomResourceDefinitions() CustomResourceDefinitionInterface {
	return newCustomResourceDefinitions(c)
}

// NewForConfig creates a new ApiextensionsV1Client for the given config.
func NewForConfig(c *rest.Config) (*ApiextensionsV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &ApiextensionsV1Client{client}, nil
}

// NewForConfigOrDie creates a new ApiextensionsV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *ApiextensionsV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new ApiextensionsV1Client for the given RESTClient.
func New(c rest.Interface) *ApiextensionsV1Client {
	return &ApiextensionsV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	//if config.UserAgent == "" {
	//  config.UserAgent = rest.DefaultKubernetesUserAgent()
	//}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *ApiextensionsV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
