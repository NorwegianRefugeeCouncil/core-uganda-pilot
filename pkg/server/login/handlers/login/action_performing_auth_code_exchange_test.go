package login

import (
	"context"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/stretchr/testify/assert"
	"testing"
)

var idp = types.IdentityProvider{
	ID: "id",
	Name: "name",
	OrganizationID: "organizationId",
	Domain: "domain",
	ClientID: "clientId",
	ClientSecret: "secret",
	EmailDomain: "emailDomain",
	Scopes: "scope1 scope2",
	ClaimMappings: types.ClaimMappings{
		Version       : "VersionFieldName",
		Subject       : "SubjectFieldName",
		DisplayName   : "DisplayNameFieldName",
		FullName      : "FullNameFieldName",
		Email         : "EmailFieldName",
		EmailVerified : "EmailVerifiedFieldName",
	},
}


func TestExtractIdentityProfile(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name          string
		claims        interface{}
		expect        authrequest.Claims
	}{
		{
			name: "maps correctly",
			claims: map[string]interface{}{
				"SubjectFieldName": "subjectValue",
				"DisplayNameFieldName": "displayNameValue",
				"FullNameFieldName": "fullNameValue",
				"EmailFieldName": "emailValue",
				"EmailVerifiedFieldName": true,
			},
			expect: authrequest.Claims{
				Subject: "subjectValue",
				DisplayName: "displayNameValue",
				FullName: "fullNameValue",
				Email: "emailValue",
				EmailVerified: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractIdentityProfile(ctx, tt.claims, &idp)

			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tt.expect, *got)
		})
	}
}
