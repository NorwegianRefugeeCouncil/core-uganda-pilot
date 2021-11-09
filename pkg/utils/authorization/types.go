package authorization

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ory/hydra-client-go/client/admin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type OIDCMetadata struct {
	AuthorizationEndpoint                     string   `json:"authorization_endpoint"`
	BackChannelLogoutSessionSupported         bool     `json:"back_channel_logout_session_supported"`
	BackChannelLogoutSupported                bool     `json:"back_channel_logout_supported"`
	ClaimsParameterSupported                  bool     `json:"claims_parameter_supported"`
	ClaimsSupported                           []string `json:"claims_supported"`
	CodeChallengeMethodsSupported             []string `json:"code_challenge_methods_supported"`
	EndSessionEndpoint                        string   `json:"end_session_endpoint"`
	FrontChannelLogoutSessionSupported        bool     `json:"front_channel_logout_session_supported"`
	FrontChannelLogoutSupported               bool     `json:"front_channel_logout_supported"`
	GrantTypesSupported                       []string `json:"grant_types_supported"`
	IdTokenSigningAlgValuesSupported          []string `json:"id_token_signing_alg_values_supported"`
	IntrospectionEndpoint                     string   `json:"introspection_endpoint"`
	IntrospectionEndpointAuthMethodsSupported []string `json:"introspection_endpoint_auth_methods_supported"`
	Issuer                                    string   `json:"issuer"`
	JWKsURI                                   string   `json:"jwks_uri"`
	RequestObjectSigningAlgValuesSupported    []string `json:"request_object_signing_alg_values_supported"`
	RequestParameterSupported                 bool     `json:"request_parameter_supported"`
	RequestURIParameterSupported              bool     `json:"request_uri_parameter_supported"`
	RequireRequestURIRegistration             bool     `json:"require_request_uri_registration"`
	ResponseModesSupported                    []string `json:"response_modes_supported"`
	ResponseTypesSupported                    []string `json:"response_types_supported"`
	RevocationEndpoint                        string   `json:"revocation_endpoint"`
	RevocationEndpointAuthMethodsSupported    []string `json:"revocation_endpoint_auth_methods_supported"`
	ScopesSupported                           []string `json:"scopes_supported"`
	SubjectTypesSupported                     []string `json:"subject_types_supported"`
	TokenEndpoint                             string   `json:"token_endpoint"`
	TokenEndpointAuthMethodsSupported         []string `json:"token_endpoint_auth_methods_supported"`
	UserInfoEndpoint                          string   `json:"user_info_endpoint"`
	UserInfoSigningAlgValuesSupported         []string `json:"user_info_signing_alg_values_supported"`
}

type TokenAudience []string

func (t *TokenAudience) UnmarshalJSON(v []byte) error {
	str := string(v)
	if strings.HasPrefix(str, "[") {
		var list []string
		if err := json.Unmarshal(v, &list); err != nil {
			return err
		}
		*t = list
		return nil
	} else {
		if len(str) != 0 {
			*t = []string{str}
		}
		return nil
	}
}

type TokenIntrospecter interface {
	Introspect(ctx context.Context, token string) (TokenIntrospectionResult, error)
}

type TokenIntrospectionResult struct {
	Active     bool          `json:"active"`
	Audience   TokenAudience `json:"aud"`
	ClientID   string        `json:"client_id"`
	Expiration int64         `json:"exp"`
	IssuedAt   int64         `json:"iat"`
	Issuer     string        `json:"iss"`
	JTI        string        `json:"jti"`
	NBF        int64         `json:"nbf"`
	Scope      string        `json:"scope"`
	Subject    string        `json:"sub"`
	TokenType  string        `json:"token_type"`
	Username   string        `json:"username"`
}

type oidcIntrospection struct {
	metadata    OIDCMetadata
	tokenSource oauth2.TokenSource
}

func NewOidcIntrospecter(
	ctx context.Context,
	issuer string,
	clientId string,
	clientSecret string,
) (TokenIntrospecter, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/.well-known/openid-configuration", issuer), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code")
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var meta = OIDCMetadata{}
	if err := json.Unmarshal(bodyBytes, &meta); err != nil {
		return nil, err
	}
	cfg := clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     meta.TokenEndpoint,
	}
	tokenSource := cfg.TokenSource(ctx)
	return &oidcIntrospection{
		metadata:    meta,
		tokenSource: tokenSource,
	}, nil
}

func (h *oidcIntrospection) Introspect(ctx context.Context, token string) (TokenIntrospectionResult, error) {
	form := url.Values{}
	form.Set("token", token)
	req, err := http.NewRequest("POST", h.metadata.IntrospectionEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return TokenIntrospectionResult{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	httpClient := oauth2.NewClient(ctx, h.tokenSource)
	resp, err := httpClient.Do(req)
	if err != nil {
		return TokenIntrospectionResult{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return TokenIntrospectionResult{}, errors.New("unexpected status code")
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TokenIntrospectionResult{}, err
	}
	result := TokenIntrospectionResult{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return TokenIntrospectionResult{}, err
	}
	return result, nil
}

func NewHydraIntrospecter(hydraAdmin admin.ClientService) TokenIntrospecter {
	return &hydraIntrospection{
		hydraAdmin: hydraAdmin,
	}
}

type hydraIntrospection struct {
	hydraAdmin admin.ClientService
}

func (h *hydraIntrospection) Introspect(ctx context.Context, token string) (TokenIntrospectionResult, error) {
	result, err := h.hydraAdmin.IntrospectOAuth2Token(&admin.IntrospectOAuth2TokenParams{
		Token:   token,
		Context: ctx,
	})
	if err != nil {
		return TokenIntrospectionResult{}, err
	}
	return TokenIntrospectionResult{
		Active:     *result.Payload.Active,
		Audience:   result.Payload.Aud,
		ClientID:   result.Payload.ClientID,
		Expiration: result.Payload.Exp,
		IssuedAt:   result.Payload.Iat,
		Issuer:     result.Payload.Iss,
		JTI:        "",
		NBF:        result.Payload.Nbf,
		Scope:      result.Payload.Scope,
		Subject:    result.Payload.Sub,
		TokenType:  result.Payload.TokenType,
		Username:   result.Payload.Username,
	}, nil
}
