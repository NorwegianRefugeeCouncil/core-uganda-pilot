package keycloak

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestCreateClientScope(t *testing.T) {

	keycloakUrl, _ := url.Parse("http://localhost:8080")

	cli := &KeycloakClient{
		baseUrl:      keycloakUrl,
		clientID:     "api",
		clientSecret: "a9a1b0f8-3920-4a6a-9c5f-2ed9e49e833e",
	}

	token, err := cli.GetToken("Core")
	if !assert.NoError(t, err) {
		return
	}

	list, err := cli.ListClientScopes(token, "Core")
	if !assert.NoError(t, err) {
		return
	}

	marshaled, err := json.MarshalIndent(list, "", "  ")
	if !assert.NoError(t, err) {
		return
	}

	t.Logf("%s", string(marshaled))

	clientScope := &ClientScope{
		Name:        "test",
		Description: "Test Client Scope",
		ProtocolMappers: []*ProtocolMapper{
			{
				Protocol:       "openid-connect",
				ProtocolMapper: "oidc-usermodel-attribute-mapper",
				Name:           "no.nrc:status",
				Configuration: map[string]string{
					"access.token.claim":   "true",
					"claim.name":           "core.nrc.no/status",
					"id.token.claim":       "true",
					"jsonType.label":       "String",
					"user.attribute":       "core.nrc.no/status",
					"userinfo.token.claim": "true",
				},
			},
		},
	}

	created, err := cli.CreateClientScope(token, "Core", clientScope)
	if !assert.NoError(t, err) {
		return
	}

	t.Logf("%#v", created)

}

func TestCreateUser(t *testing.T) {

	keycloakUrl, _ := url.Parse("http://localhost:8080")

	user := &User{
		Username: "test",
		Attributes: map[string][]string{
			"someattribute": {"abc", "def"},
		},
	}

	cli := &KeycloakClient{
		baseUrl:      keycloakUrl,
		clientID:     "api",
		clientSecret: "a9a1b0f8-3920-4a6a-9c5f-2ed9e49e833e",
	}

	token, err := cli.GetToken("Core")
	if !assert.NoError(t, err) {
		return
	}

	listUsers, err := cli.ListUsers(token, "Core", ListUserOptions{
		Username: "test",
	})
	if !assert.NoError(t, err) {
		return
	}

	for _, listedUser := range listUsers {
		if listedUser.Username == "test" {
			if err := cli.DeleteUser(token, "Core", listedUser.ID); !assert.NoError(t, err) {
				return
			}
		}
	}

	gotUser, err := cli.CreateUser(token, "Core", user)
	if !assert.NoError(t, err) {
		return
	}

	t.Logf("%#v", gotUser)

	if err := cli.DeleteUser(token, "Core", gotUser.ID); !assert.NoError(t, err) {
		return
	}

}
