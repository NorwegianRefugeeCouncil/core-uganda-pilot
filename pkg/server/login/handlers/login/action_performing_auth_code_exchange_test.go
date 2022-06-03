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
		Version       : "v1",
		Subject       : "{{.SubjectFieldName}}",
		DisplayName   : "{{.Title}} {{.FirstName}} {{.LastName}}",
		FullName      : "{{.FirstName}} {{.LastName}}",
		Email         : "{{.EmailFieldName}}",
		EmailVerified : "{{.EmailVerifiedFieldName}}",
	},
}
var idp2 = types.IdentityProvider{
	ID: "id",
	Name: "name",
	OrganizationID: "organizationId",
	Domain: "domain",
	ClientID: "clientId",
	ClientSecret: "secret",
	EmailDomain: "emailDomain",
	Scopes: "scope1 scope2",
	ClaimMappings: types.ClaimMappings{
		Version       : "v1",
		Subject       : "SubjectFieldName",
	},
}


func TestExtractIdentityProfile(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name          string
		claims        interface{}
		idp           types.IdentityProvider
		expect        authrequest.Claims
		wantErr       bool
	}{
		{
			name: "maps correctly",
			claims: map[string]interface{}{
				"SubjectFieldName": "subjectValue",
				"Title": "Prof.",
				"FirstName": "Neil",
				"LastName": "Armstrong",
				"EmailFieldName": "emailValue",
				"EmailVerifiedFieldName": true,
			},
			idp: idp,
			expect: authrequest.Claims{
				Subject: "subjectValue",
				DisplayName: "Prof. Neil Armstrong",
				FullName: "Neil Armstrong",
				Email: "emailValue",
				EmailVerified: true,
			},
			wantErr: false,
		},
		{
			name: "fails to map data if IDP claims are not setup correctly",
			claims: map[string]interface{}{
				"SubjectFieldName": "subjectValue",
				"Title": "Prof.",
				"FirstName": "Neil",
				"LastName": "Armstrong",
				"EmailFieldName": "emailValue",
				"EmailVerifiedFieldName": true,
			},
			idp: idp2,
			expect: authrequest.Claims{
				Subject: "SubjectFieldName",
				DisplayName: "",
				FullName: "",
				Email: "",
				EmailVerified: false,
			},
			wantErr: false,
		},
		{
			name: "fills profile with placeholders for missing data",
			claims: map[string]interface{}{
				"SubjectFieldName": "subjectValue",
			},
			expect: authrequest.Claims{
				Subject: "subjectValue",
				DisplayName: "<no value> <no value> <no value>",
				FullName: "<no value> <no value>",
				Email: "<no value>",
				EmailVerified: false,
			},
			wantErr: false,
		},
		{
			name: "fails if claims is not a map[string]interface{}",
			claims: 3,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractIdentityProfileTemplateVersion(ctx, tt.claims, &tt.idp)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tt.expect, *got)
		})
	}
}
