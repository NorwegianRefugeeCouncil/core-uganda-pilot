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

	// env file
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("REACT_APP_ISSUER=%s\n", HydraHost))
	sb.WriteString(fmt.Sprintf("REACT_APP_CLIENT_ID=%s\n", c.coreAppFrontendClientId))
	sb.WriteString(fmt.Sprintf("REACT_APP_SERVER_URL=https://localhost:8443\n"))
	envPath := path.Join(c.rootDir, "web", "apps", "pwa", ".env")
	fmt.Println(envPath)
	if err := ioutil.WriteFile(envPath, []byte(sb.String()), os.ModePerm); err != nil {
		return err
	}

	c.hydraClients = append(c.hydraClients, ClientConfig{
		ClientId: c.coreAppFrontendClientId,
		RedirectUris: []string{
			fmt.Sprintf("%s/app", CoreHost),
			fmt.Sprintf("http://localhost:3000/app"),
		},
		GrantTypes: []string{
			"authorization_code",
			"refresh_token",
		},
		TokenEndpointAuthMethod: "none",
		Scope:                   "openid profile email offline_access",
		ResponseTypes:           []string{"code", "code id_token token"},
	})

	return nil
}
