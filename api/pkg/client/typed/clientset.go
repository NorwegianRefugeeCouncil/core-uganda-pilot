package typed

import (
	"github.com/nrc-no/core/api/pkg/client/rest"
	corev1client "github.com/nrc-no/core/api/pkg/client/typed/core/v1"
	discoveryv1 "github.com/nrc-no/core/api/pkg/client/typed/discovery/v1"
)

type Interface interface {
	DiscoveryV1() discoveryv1.DiscoveryV1Interface
	CoreV1() corev1client.CoreV1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	coreV1      *corev1client.CoreV1Client
	discoveryV1 *discoveryv1.DiscoveryV1Client
}

// CoreV1 retrieves the DiscoveryV1Client
func (c *Clientset) CoreV1() corev1client.CoreV1Interface {
	return c.coreV1
}

// DiscoveryV1 retrieves the DiscoveryV1Client
func (c *Clientset) DiscoveryV1() discoveryv1.DiscoveryV1Interface {
	return c.discoveryV1
}

// NewForConfig creates a new Clientset for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfig will generate a rate-limiter in configShallowCopy.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	var cs Clientset
	var err error
	cs.coreV1, err = corev1client.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.discoveryV1, err = discoveryv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.coreV1 = corev1client.NewForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.coreV1 = corev1client.New(c)
	return &cs
}
