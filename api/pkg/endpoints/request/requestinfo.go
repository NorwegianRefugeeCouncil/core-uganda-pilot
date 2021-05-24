package request

import (
	"context"
	"net/http"
	"strings"
)

type RequestInfo struct {
	IsResourceRequest bool
	Path              string
	Verb              string
	APIPrefix         string
	APIGroup          string
	APIVersion        string
	Resource          string
	Name              string
	Parts             []string
}

func NewRequestInfo(req *http.Request) (*RequestInfo, error) {

	requestInfo := RequestInfo{
		IsResourceRequest: false,
		Path:              req.URL.Path,
		Verb:              strings.ToLower(req.Method),
	}

	switch req.Method {
	case "POST":
		requestInfo.Verb = "create"
	case "GET", "HEAD":
		requestInfo.Verb = "get"
	case "PUT":
		requestInfo.Verb = "update"
	case "DELETE":
		requestInfo.Verb = "delete"
	default:
		requestInfo.Verb = ""
	}

	currentParts := splitPath(req.URL.Path)
	if len(currentParts) < 4 {
		// return a non-resource request
		return &requestInfo, nil
	}

	requestInfo.APIPrefix = currentParts[0]
	currentParts = currentParts[1:]

	requestInfo.APIGroup = currentParts[0]
	currentParts = currentParts[1:]

	requestInfo.IsResourceRequest = true
	requestInfo.APIVersion = currentParts[0]
	currentParts = currentParts[1:]

	requestInfo.Parts = currentParts
	switch {
	case len(requestInfo.Parts) >= 2:
		requestInfo.Name = requestInfo.Parts[1]
		fallthrough
	case len(requestInfo.Parts) >= 1:
		requestInfo.Resource = requestInfo.Parts[0]
	}

	if len(requestInfo.Name) == 0 && requestInfo.Verb == "get" {
		requestInfo.Verb = "list"
	}

	return &requestInfo, nil
}

const requestInfoKey string = "RequestInfo"

// WithRequestInfo returns a copy of parent in which the request info value is set
func WithRequestInfo(parent context.Context, info *RequestInfo) context.Context {
	return WithValue(parent, requestInfoKey, info)
}

// RequestInfoFrom returns the value of the RequestInfo key on the ctx
func RequestInfoFrom(ctx context.Context) (*RequestInfo, bool) {
	info, ok := ctx.Value(requestInfoKey).(*RequestInfo)
	return info, ok
}

// splitPath returns the segments for a URL path.
func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}
