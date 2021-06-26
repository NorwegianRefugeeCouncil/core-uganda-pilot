package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/square/go-jose.v2"
	jwt2 "gopkg.in/square/go-jose.v2/jwt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var jwks jose.JSONWebKeySet

func (s *Server) WithAuth(ctx context.Context) func(handler http.Handler) http.Handler {

	req, err := http.NewRequest("GET", "http://localhost:4444/.well-known/jwks.json", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != http.StatusOK {
		panic("unexpected status code")
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(bodyBytes, &jwks); err != nil {
		panic(err)
	}

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			authorization := req.Header.Get("Authorization")
			if len(authorization) == 0 {
				s.Error(w, fmt.Errorf("unauthorized"))
				return
			}

			parts := strings.Split(authorization, " ")
			if len(parts) != 2 {
				s.Error(w, fmt.Errorf("unauthorized"))
				return
			}

			if parts[0] != "Bearer" {
				s.Error(w, fmt.Errorf("unauthorized"))
				return
			}

			t, err := jwt2.ParseSigned(parts[0])
			if err != nil {
				s.Error(w, fmt.Errorf("unauthorized"))
				return
			}

			var c jwt2.Claims
			if err := t.Claims(jwks.Keys[0].Key, &c); err != nil {
				s.Error(w, fmt.Errorf("unauthorized"))
				return
			}

			if err := c.Validate(jwt2.Expected{
				Issuer:   "http://localhost:4444",
				Audience: []string{},
				Time:     time.Now(),
			}); err != nil {
				s.Error(w, fmt.Errorf("unauthorized"))
				return
			}

			handler.ServeHTTP(w, req)
		})
	}
}
