package endpoints

import (
	"net/http"
	"strings"
)

type RequestInfo struct {
	ResourceID  string
	APIGroup    string
	APIVersion  string
	APIResource string
	Verb        string
	Path        string
}

type RequestInfoResolver interface {
	NewRequestInfo(req *http.Request) (*RequestInfo, error)
}

type RequestInfoFactory struct {
}

var DefaultRequestInfoFactory = &RequestInfoFactory{}

func (r *RequestInfoFactory) NewRequestInfo(req *http.Request) (*RequestInfo, error) {
	requestInfo := &RequestInfo{
		Verb: strings.ToLower(req.Method),
		Path: req.URL.Path,
	}

	var parts []string
	path := strings.Trim(req.URL.Path, "/")
	if path == "" {
		parts = []string{}
	} else {
		parts = strings.Split(path, "/")
	}

	if len(parts) > 1 {
		requestInfo.APIGroup = parts[1]
	}
	if len(parts) > 2 {
		requestInfo.APIVersion = parts[2]
	}
	if len(parts) > 3 {
		requestInfo.APIResource = parts[3]
	}
	if len(parts) > 4 {
		requestInfo.ResourceID = parts[4]
	}

	return requestInfo, nil

}
