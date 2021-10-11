package utils

import (
	"context"
	"github.com/nrc-no/core/internal/sessionmanager"
	"golang.org/x/oauth2"
	"net/http"
)

func GetAccessAndRefreshTokens(s sessionmanager.Store, req *http.Request) (string, string) {
	accessToken, _ := s.FindString(req, "access-token")
	refreshToken, _ := s.FindString(req, "refresh-token")
	return accessToken, refreshToken
}

func GetOauth2Token(s sessionmanager.Store, req *http.Request) *oauth2.Token {
	accessToken, refreshToken := GetAccessAndRefreshTokens(s, req)
	if len(accessToken) > 0 || len(refreshToken) > 0 {
		return &oauth2.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
	}
	return nil
}

func GetOauth2HttpClient(
	s sessionmanager.Store,
	req *http.Request,
	oauth2Config *oauth2.Config,
	defaultClient *http.Client,
) (*http.Client, error) {
	ctx := req.Context()
	httpClient := defaultClient
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	token := GetOauth2Token(s, req)
	if token != nil {
		httpClient = oauth2Config.Client(ctx, token)
	}

	return httpClient, nil
}
