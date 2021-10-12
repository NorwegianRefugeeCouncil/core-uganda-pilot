package cms

import (
	"github.com/nrc-no/core/pkg/rest"
)

type ClientSet struct {
	c *rest.Client
}

var _ Interface = &ClientSet{}

func NewClientSet(restConfig *rest.Config) *ClientSet {
	return &ClientSet{
		c: rest.NewClient(restConfig),
	}
}

func (c ClientSet) Cases() CaseClient {
	return &RESTCaseClient{
		c: c.c,
	}
}

func (c ClientSet) CaseTypes() CaseTypeClient {
	return &RESTCaseTypeClient{
		c: c.c,
	}
}

func (c ClientSet) Comments() CommentClient {
	return &RestCommentClient{
		c: c.c,
	}
}
