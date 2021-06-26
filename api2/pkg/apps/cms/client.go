package cms

import "github.com/nrc-no/core-kafka/pkg/rest"

type ClientSet struct {
	c *rest.Client
}

var _ Interface = &ClientSet{}

func NewClientSet(restConfig *rest.RESTConfig) *ClientSet {
	return &ClientSet{
		c: rest.NewClient(restConfig),
	}
}

func (c ClientSet) Cases() CaseClient {
	panic("implement me")
}

func (c ClientSet) CaseTypes() CaseTypeClient {
	panic("implement me")
}
