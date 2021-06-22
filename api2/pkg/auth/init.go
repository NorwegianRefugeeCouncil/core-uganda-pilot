package auth

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/keycloak"
)

func Init(ctx context.Context, keycloakClient *keycloak.Client) error {

	token, err := keycloakClient.GetToken(keycloak.GetTokenOptions{
		GrantType: "client_credentials",
	})
	if err != nil {
		return err
	}

	users, err := keycloakClient.ListUsers(token, keycloak.ListUserOptions{
		Username: "admin",
	})
	if err != nil {
		return err
	}

	if len(users) == 0 {

		_, err := keycloakClient.CreateUser(token, &keycloak.User{
			Credentials: []keycloak.Credential{
				{
					Type:      "password",
					Temporary: false,
					Value:     "admin",
				},
			},
			EmailVerified: true,
			Enabled:       true,
			FirstName:     "Admin",
			Username:      "admin",
		})
		if err != nil {
			return err
		}
	}

	return nil

}
