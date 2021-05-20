package handlers

import (
	"github.com/nrc-no/core/api/pkg/server2/endpoints/request"
	"net/http"
)

// Namer is able to retrieve the name of a resource from the http request
type Namer interface {
	Name(req *http.Request) (name string, err error)
}

// ContextBasedNaming is an implementation of Namer and can retrieve
// the name using the request.NewRequestInfo
type ContextBasedNaming struct {
}

var _ Namer = &ContextBasedNaming{}

func (c ContextBasedNaming) Name(req *http.Request) (name string, err error) {
	info, err := request.NewRequestInfo(req)
	if err != nil {
		return "", err
	}
	return info.Name, nil
}
