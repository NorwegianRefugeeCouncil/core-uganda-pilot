package client

import v1 "github.com/nrc-no/core/api/pkg/client/core/v1"

type Interface interface {
	CoreV1() v1.CoreV1Interface
}

type ClientSet struct {
	corev1 *v1.CoreV1Client
}

func (c *ClientSet) CoreV1() v1.CoreV1Interface {
	return c.corev1
}
