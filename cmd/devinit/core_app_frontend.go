package devinit

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
)

func (c *Config) makeAppFrontend() error {

	var err error

	c.coreAppFrontendTlsKey, err = getOrCreatePrivateKey(path.Join(CoreAppFrontendDir, "tls.key"))
	if err != nil {
		return err
	}

	c.coreAppFrontendTlsCert, err = getOrCreateServerCert(
		path.Join(CoreAppFrontendDir, "tls.crt"),
		c.coreAppFrontendTlsKey,
		c.rootCa,
		c.rootCaKey,
		[]string{"localhost", "core.dev"},
		[]net.IP{net.IPv6loopback, net.ParseIP("127.0.0.1")},
	)
	if err != nil {
		return err
	}

	c.coreAppFrontendClientId = "core-app-frontend"
	const scope = "openid profile email offline_access"

	// env file
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("REACT_APP_OIDC_ISSUER=%s\n", HydraHost))
	sb.WriteString(fmt.Sprintf("REACT_APP_OAUTH_SCOPE=%s\n", scope))
	sb.WriteString(fmt.Sprintf("REACT_APP_OAUTH_REDIRECT_URI=%s\n", PwaURI))
	sb.WriteString(fmt.Sprintf("REACT_APP_OAUTH_CLIENT_ID=%s\n", c.coreAppFrontendClientId))
	sb.WriteString(fmt.Sprintf("REACT_APP_SERVER_URL=https://localhost:8443\n"))

	envPath := path.Join(c.rootDir, "frontend", "apps", "pwa", ".env")
	fmt.Println(envPath)
	if err := ioutil.WriteFile(envPath, []byte(sb.String()), os.ModePerm); err != nil {
		return err
	}

	c.hydraClients = append(c.hydraClients, ClientConfig{
		ClientId: c.coreAppFrontendClientId,
		RedirectUris: []string{
			fmt.Sprintf(PwaURI),
		},
		GrantTypes: []string{
			"authorization_code",
			"refresh_token",
		},
		TokenEndpointAuthMethod: "none",
		Scope:                   scope,
		ResponseTypes:           []string{"code", "code id_token token"},
	})

	return nil
}
