package devinit

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
)

func (c *Config) makeOidcConfig() error {
	var err error

	c.oidcTlsKey, err = getOrCreatePrivateKey(path.Join(OIDCDir, "tls.key"))
	if err != nil {
		return err
	}

	c.oidcTlsCert, err = getOrCreateServerCert(
		path.Join(OIDCDir, "tls.crt"),
		c.oidcTlsKey,
		c.rootCa,
		c.rootCaKey,
		[]string{"localhost", "oidc.dev"},
		[]net.IP{net.IPv6loopback, net.ParseIP("127.0.0.1")},
	)
	if err != nil {
		return err
	}

	c.oidcConfig.Scopes = []string{
		"openid",
		"profile",
		"email",
		"offline_access",
	}

	jsonBytes, err := json.MarshalIndent(c.oidcConfig, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path.Join(OIDCDir, "config.json"), jsonBytes, os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(
		path.Join(OIDCDir, "env"),
		[]byte(fmt.Sprintf("ISSUER=%s", OidcIssuer)),
		os.ModePerm); err != nil {
		return err
	}

	return nil
}

type ClientConfig struct {
	ClientId                string   `json:"client_id"`
	RedirectUris            []string `json:"redirect_uris"`
	GrantTypes              []string `json:"grant_types"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
	Scope                   string   `json:"scope"`
	ResponseTypes           []string `json:"response_types"`
	ClientSecret            string   `json:"client_secret,omitempty"`
}

type OidcConfig struct {
	Scopes  []string       `json:"scopes"`
	Clients []ClientConfig `json:"clients"`
}
