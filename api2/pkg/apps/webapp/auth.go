package webapp

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/auth"
	jwt2 "gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"time"
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

	reader := auth.TokenReader(ctx)

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			ctx := req.Context()

			if req.URL.Path == "/login" || req.URL.Path == "/callback" {
				handler.ServeHTTP(w, req)
				return
			}

			accessToken := s.sessionManager.GetString(ctx, "access-token")
			refreshToken := s.sessionManager.GetString(ctx, "refresh-token")

			if len(accessToken) == 0 {
				http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
				return
			}

			accessTokenClaims := AccessToken{}
			if err := reader(auth.StaticTokenSource(accessToken), &accessTokenClaims); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			refreshTokenClaims := RefreshToken{}
			if err := reader(auth.StaticTokenSource(refreshToken), &refreshTokenClaims); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			refreshTokenExpired := false
			if err := refreshTokenClaims.GetClaims().Validate(jwt2.Expected{
				Time: time.Now().Add(-1 * time.Minute),
			}); err != nil {
				refreshTokenExpired = true
			}

			accessTokenExpired := false
			if err := accessTokenClaims.GetClaims().Validate(jwt2.Expected{
				Time: time.Now().Add(-1 * time.Minute),
			}); err != nil {
				accessTokenExpired = true
			}

			if !accessTokenExpired {
				handler.ServeHTTP(w, req)
				return
			}

			if accessTokenExpired && !refreshTokenExpired {

				// use refresh token
			}

			if accessTokenExpired && refreshTokenExpired {
				// ask for login again
			}

		})
	}
}
