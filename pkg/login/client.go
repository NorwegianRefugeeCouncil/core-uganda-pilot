package login

import (
	"github.com/nrc-no/core/pkg/rest"
)

type ClientSet struct {
	c *rest.Client
}

var _ Interface = &ClientSet{}

func NewClientSet(restConfig *rest.RESTConfig) *ClientSet {
	return &ClientSet{
		c: rest.NewClient(restConfig),
	}
}

func (c *ClientSet) Login() Client {
	return &RESTLoginClient{
		c: c.c,
	}
}
