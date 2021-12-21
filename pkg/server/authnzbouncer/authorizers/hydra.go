package authorizers

import (
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/utils/authorization"
	"github.com/ory/hydra-client-go/client/admin"
	"net/http"
)

type hydraAuthorizer struct {
	hydraAdmin admin.ClientService
}

func NewHydraAuthorizer(hydraAdmin admin.ClientService) Authorizer {
	return &hydraAuthorizer{hydraAdmin: hydraAdmin}
}

type AuthorizationResponse struct {
	Active            bool        `json:"active"`
	Aud               []string    `json:"aud"`
	ClientID          string      `json:"client_id,omitempty"`
	Exp               int64       `json:"exp,omitempty"`
	Ext               interface{} `json:"ext,omitempty"`
	Iat               int64       `json:"iat,omitempty"`
	Iss               string      `json:"iss,omitempty"`
	Nbf               int64       `json:"nbf,omitempty"`
	ObfuscatedSubject string      `json:"obfuscated_subject,omitempty"`
	Scope             string      `json:"scope,omitempty"`
	Sub               string      `json:"sub,omitempty"`
	TokenType         string      `json:"token_type,omitempty"`
	TokenUse          string      `json:"token_use,omitempty"`
	Username          string      `json:"username,omitempty"`
}

type Authorizer interface {
	Authorize(req *http.Request) (AuthorizationResponse, error)
}

func (h *hydraAuthorizer) Authorize(req *http.Request) (AuthorizationResponse, error) {

	ctx := req.Context()

	bearerToken, err := authorization.ExtractBearerToken(req)
	if err != nil {
		return AuthorizationResponse{}, err
	}

	introspectionResponse, err := h.hydraAdmin.IntrospectOAuth2Token(&admin.IntrospectOAuth2TokenParams{
		Token:   bearerToken,
		Context: ctx,
	})
	if err != nil {
		return AuthorizationResponse{}, meta.NewInternalServerError(err)
	}

	return AuthorizationResponse{
		Active:            *introspectionResponse.Payload.Active,
		Aud:               introspectionResponse.Payload.Aud,
		ClientID:          introspectionResponse.Payload.ClientID,
		Exp:               introspectionResponse.Payload.Exp,
		Ext:               introspectionResponse.Payload.Ext,
		Iat:               introspectionResponse.Payload.Iat,
		Iss:               introspectionResponse.Payload.Iss,
		Nbf:               introspectionResponse.Payload.Nbf,
		ObfuscatedSubject: introspectionResponse.Payload.ObfuscatedSubject,
		Scope:             introspectionResponse.Payload.Scope,
		Sub:               introspectionResponse.Payload.Sub,
		TokenType:         introspectionResponse.Payload.TokenType,
		TokenUse:          introspectionResponse.Payload.TokenUse,
		Username:          introspectionResponse.Payload.Username,
	}, nil

}
