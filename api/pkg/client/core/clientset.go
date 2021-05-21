package core

import (
	v1 "github.com/nrc-no/core/api/pkg/client/core/v1"
	"github.com/nrc-no/core/api/pkg/client/rest"
)

type Interface interface {
	CoreV1() v1.CoreV1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	coreV1 *v1.CoreV1Client
}

// CoreV1 retrieves the CoreV1Client
func (c *Clientset) CoreV1() v1.CoreV1Interface {
	return c.coreV1
}

// NewForConfig creates a new Clientset for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfig will generate a rate-limiter in configShallowCopy.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	var cs Clientset
	var err error
	cs.coreV1, err = v1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.coreV1 = v1.NewForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.coreV1 = v1.New(c)
	return &cs
}
