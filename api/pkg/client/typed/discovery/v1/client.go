package v1

import (
	discoveryv1 "github.com/nrc-no/core/api/pkg/apis/discovery/v1"
	"github.com/nrc-no/core/api/pkg/client/rest"
	"github.com/nrc-no/core/api/pkg/client/typed/scheme"
)

type DiscoveryV1Interface interface {
	APIServicesGetter
}

// DiscoveryV1Client is used to interact with features provided by the discovery.nrc.no group.
type DiscoveryV1Client struct {
	restClient rest.Interface
}

func (c *DiscoveryV1Client) APIServices() APIServiceInterface {
	return newAPIServices(c)
}

// NewForConfig creates a new DiscoveryV1Client for the given config.
func NewForConfig(c *rest.Config) (*DiscoveryV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &DiscoveryV1Client{client}, nil
}

// NewForConfigOrDie creates a new DiscoveryV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *DiscoveryV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new DiscoveryV1Client for the given RESTClient.
func New(c rest.Interface) *DiscoveryV1Client {
	return &DiscoveryV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := discoveryv1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *DiscoveryV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
