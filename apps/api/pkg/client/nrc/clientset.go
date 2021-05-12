package nrc

import (
	"github.com/nrc-no/core/apps/api/pkg/client/nrc/typed/core/v1"
	"github.com/nrc-no/core/apps/api/pkg/client/rest"
)

type Interface interface {
	CoreV1() v1.CoreV1Interface
}

type Clientset struct {
	corev1 *v1.NrcCoreClient
}

func (c *Clientset) CoreV1() v1.CoreV1Interface {
	return c.corev1
}

func NewForConfig(c *rest.Config) (*Clientset, error) {
	shallowCopy := *c
	var cs Clientset
	var err error
	cs.corev1, err = v1.NewForConfig(&shallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.corev1 = v1.New(c)
	return &cs
}
