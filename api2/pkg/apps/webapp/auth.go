package webapp

import (
	"context"
	"github.com/ory/hydra-client-go/client/admin"
	jwt2 "gopkg.in/square/go-jose.v2/jwt"
	"net/http"
)

type AccessToken struct {
	jwt2.Claims
}

func (a *AccessToken) GetClaims() *jwt2.Claims {
	return &a.Claims
}

type RefreshToken struct {
	jwt2.Claims
}

func (a *RefreshToken) GetClaims() *jwt2.Claims {
	return &a.Claims
}

func (s *Server) WithAuth(ctx context.Context) func(handler http.Handler) http.Handler {

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			if req.URL.Path == "/login" || req.URL.Path == "/callback" {
				handler.ServeHTTP(w, req)
				return
			}

			token := s.sessionManager.GetString(req.Context(), "access-token")

			res, err := s.HydraAdmin.IntrospectOAuth2Token(&admin.IntrospectOAuth2TokenParams{
				Token:      token,
				Context:    req.Context(),
				HTTPClient: nil,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !*res.Payload.Active {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			handler.ServeHTTP(w, req)
		})
	}

}
