package cms

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/auth"
	jwt2 "gopkg.in/square/go-jose.v2/jwt"
	"net/http"
)

type RequestClaims struct {
	jwt2.Claims
}

func (c *RequestClaims) GetClaims() *jwt2.Claims {
	return &c.Claims
}

func (s *Server) WithAuth(ctx context.Context) func(handler http.Handler) http.Handler {
	tokenReader := auth.TokenReader(ctx)
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			c := &RequestClaims{}
			if err := tokenReader(req, c); err != nil {
				s.Error(w, err)
				return
			}
			handler.ServeHTTP(w, req)
		})
	}
}
