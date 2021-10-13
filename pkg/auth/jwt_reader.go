package auth

import (
	"fmt"
	"net/http"
	"strings"
)

type TokenSource interface {
	GetToken() (string, error)
}

type DelegateTokenSource struct {
	getToken func() (string, error)
}

func NewDelegateTokenSource(getToken func() (string, error)) *DelegateTokenSource {
	return &DelegateTokenSource{
		getToken: getToken,
	}
}
func (d *DelegateTokenSource) GetToken() (string, error) {
	return d.getToken()
}

func HeaderTokenSource(req *http.Request) TokenSource {
	return NewDelegateTokenSource(func() (string, error) {
		authorization := req.Header.Get("Authorization")
		if len(authorization) == 0 {
			return "", fmt.Errorf("authorization header not present")
		}
		parts := strings.Split(authorization, " ")
		if parts[0] != "Bearer" {
			return "", fmt.Errorf("malformed authorization header")
		}
		return parts[1], nil
	})
}
