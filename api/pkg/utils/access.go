package utils

import (
	"context"
	"github.com/nrc-no/core/pkg/sessionmanager"
	"golang.org/x/oauth2"
	"net/http"
)

func GetAccessAndRefreshTokens(s sessionmanager.Store, req *http.Request) (string, string, error) {
	accessToken, err := s.GetString(req, "access-token")
	if err != nil {
		return "", "", err
	}
	refreshToken, err := s.GetString(req, "refresh-token")
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func GetOauth2Token(s sessionmanager.Store, req *http.Request) (*oauth2.Token, error) {
	accessToken, refreshToken, err := GetAccessAndRefreshTokens(s, req)
	if err != nil {
		return nil, err
	}
	token := &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return token, nil
}

func GetOauth2HttpClient(
	s sessionmanager.Store,
	req *http.Request,
	oauth2Config *oauth2.Config,
	defaultClient *http.Client,
) (*http.Client, error) {
	token, err := GetOauth2Token(s, req)
	if err != nil {
		return nil, err
	}
	httpClient := defaultClient
	ctx := context.WithValue(req.Context(), oauth2.HTTPClient, httpClient)
	cli := oauth2Config.Client(ctx, token)
	if len(token.AccessToken) > 0 || len(token.RefreshToken) > 0 {
		httpClient = cli
	}
	return httpClient, nil
}
