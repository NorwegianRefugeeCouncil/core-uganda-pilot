package v1

import "github.com/nrc-no/coreapi/pkg/client/rest"

type CoreV1Interface interface {
	RESTClient() rest.Interface
	FormDefinitions() FormDefinitionInterface
}

type CoreV1Client struct {
	restClient rest.Interface
}

func (c *CoreV1Client) FormDefinitions() FormDefinitionInterface {
	return newFormDefinitions(c)
}

func (c *CoreV1Client) RESTClient() rest.Interface {
	return c.restClient
}

func New(c rest.Interface) *CoreV1Client {
	return &CoreV1Client{c}
}
