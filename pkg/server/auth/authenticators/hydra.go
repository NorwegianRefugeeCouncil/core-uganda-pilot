package authenticators

import (
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/utils/authorization"
	"github.com/ory/hydra-client-go/client/public"
	"net/http"
)

type hydraAuthenticator struct {
	hydraPublic public.ClientService
}

func NewHydraAuthenticator(hydraPublic public.ClientService) Authenticator {
	return &hydraAuthenticator{hydraPublic: hydraPublic}
}

type AuthenticationResponse struct {
	Birthdate           string `json:"birthdate,omitempty"`
	Email               string `json:"email,omitempty"`
	EmailVerified       bool   `json:"email_verified,omitempty"`
	FamilyName          string `json:"family_name,omitempty"`
	Gender              string `json:"gender,omitempty"`
	GivenName           string `json:"given_name,omitempty"`
	Locale              string `json:"locale,omitempty"`
	MiddleName          string `json:"middle_name,omitempty"`
	Name                string `json:"name,omitempty"`
	Nickname            string `json:"nickname,omitempty"`
	PhoneNumber         string `json:"phone_number,omitempty"`
	PhoneNumberVerified bool   `json:"phone_number_verified,omitempty"`
	Picture             string `json:"picture,omitempty"`
	PreferredUsername   string `json:"preferred_username,omitempty"`
	Profile             string `json:"profile,omitempty"`
	Sub                 string `json:"sub,omitempty"`
	UpdatedAt           int64  `json:"updated_at,omitempty"`
	Website             string `json:"website,omitempty"`
	Zoneinfo            string `json:"zoneinfo,omitempty"`
}

type Authenticator interface {
	Authenticate(req *http.Request) (AuthenticationResponse, error)
}

func (h *hydraAuthenticator) Authenticate(req *http.Request) (AuthenticationResponse, error) {

	ctx := req.Context()

	bearerToken, err := authorization.ExtractBearerToken(req)
	if err != nil {
		return AuthenticationResponse{}, err
	}

	userInfo, err := h.hydraPublic.Userinfo(&public.UserinfoParams{
		Context:    ctx,
		HTTPClient: nil,
	}, runtime.ClientAuthInfoWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
		if err := req.SetHeaderParam("Authorization", fmt.Sprintf("Bearer %s", bearerToken)); err != nil {
			return err
		}
		return nil
	}))
	if err != nil {
		return AuthenticationResponse{}, meta.NewUnauthorized("failed to get user info")
	}

	return AuthenticationResponse{
		Birthdate:           userInfo.Payload.Birthdate,
		Email:               userInfo.Payload.Email,
		EmailVerified:       userInfo.Payload.EmailVerified,
		FamilyName:          userInfo.Payload.FamilyName,
		Gender:              userInfo.Payload.Gender,
		GivenName:           userInfo.Payload.GivenName,
		Locale:              userInfo.Payload.Locale,
		MiddleName:          userInfo.Payload.MiddleName,
		Name:                userInfo.Payload.Name,
		Nickname:            userInfo.Payload.Nickname,
		PhoneNumber:         userInfo.Payload.PhoneNumber,
		PhoneNumberVerified: userInfo.Payload.PhoneNumberVerified,
		Picture:             userInfo.Payload.Picture,
		PreferredUsername:   userInfo.Payload.PreferredUsername,
		Profile:             userInfo.Payload.Profile,
		Sub:                 userInfo.Payload.Sub,
		UpdatedAt:           userInfo.Payload.UpdatedAt,
		Website:             userInfo.Payload.Website,
		Zoneinfo:            userInfo.Payload.Zoneinfo,
	}, nil

}
