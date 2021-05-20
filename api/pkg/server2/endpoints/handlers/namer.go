package handlers

import (
	"github.com/nrc-no/core/api/pkg/server2/endpoints/request"
	"net/http"
)

type Namer interface {
	Name(req *http.Request) (name string, err error)
}

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
