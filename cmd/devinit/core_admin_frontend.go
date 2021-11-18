package devinit

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
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
		RedirectUris:            []string{"http://localhost:3001/app"},
		GrantTypes:              []string{"authorization_code", "refresh_token"},
		TokenEndpointAuthMethod: "none",
		Scope:                   "openid email profile",
		ResponseTypes:           []string{"code"},
	})

	// env file
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("REACT_APP_ISSUER=%s\n", OidcIssuer))
	sb.WriteString(fmt.Sprintf("REACT_APP_CLIENT_ID=%s\n", c.coreAdminFrontendClientId))
	envPath := path.Join(c.rootDir, "web", "apps", "admin", ".env")
	fmt.Println(envPath)
	if err := ioutil.WriteFile(envPath, []byte(sb.String()), os.ModePerm); err != nil {
		return err
	}

	return nil
}
