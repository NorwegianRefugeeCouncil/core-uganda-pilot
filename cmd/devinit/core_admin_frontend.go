package devinit

import (
	"net"
	"path"
)

func (c *Config) makeAdminFrontend() error {

	var err error

	c.coreAdminFrontendTlsKey, err = getOrCreatePrivateKey(path.Join(CoreAdminFrontendDir, "tls.key"))
	if err != nil {
		return err
	}

	c.coreAdminFrontendTlsCert, err = getOrCreateServerCert(
		path.Join(CoreAdminFrontendDir, "tls.crt"),
		c.coreAdminFrontendTlsKey,
		c.rootCa,
		c.rootCaKey,
		[]string{"localhost", "core.dev"},
		[]net.IP{net.IPv6loopback, net.ParseIP("127.0.0.1")},
	)

	if err != nil {
		return err
	}

	c.coreAdminFrontendClientId = "core-admin-frontend"

	c.oidcConfig.Clients = append(c.oidcConfig.Clients, ClientConfig{
		ClientId:                c.coreAdminFrontendClientId,
		RedirectUris:            []string{},
		GrantTypes:              []string{"authorization_code", "refresh_token"},
		TokenEndpointAuthMethod: "none",
		Scope:                   "openid email profile",
		ResponseTypes:           []string{"code"},
	})

	return nil
}
