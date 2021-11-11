package devinit

import "fmt"

func (c *Config) makeIdp() error {

	var err error
	c.idpIssuer = OidcIssuer
	c.idpClientId = "nrc-external-idp"
	c.idpClientSecret, err = getOrCreateRandomSecretStr(32, IDPDir, "client-secret")
	if err != nil {
		return err
	}

	c.oidcConfig.Clients = append(c.oidcConfig.Clients, ClientConfig{
		ClientId: c.idpClientId,
		RedirectUris: []string{
			fmt.Sprintf("%s/login/oidc/callback", CoreLocalHost),
			fmt.Sprintf("%s/login/oidc/callback", CoreHost),
			"https://nrc-core-server.loca.lt/login/oidc/callback",
		},
		GrantTypes:              []string{"authorization_code", "refresh_token"},
		TokenEndpointAuthMethod: "client_secret_post",
		Scope:                   "openid email profile",
		ResponseTypes:           []string{"code"},
		ClientSecret:            c.idpClientSecret,
	})

	return nil
}
