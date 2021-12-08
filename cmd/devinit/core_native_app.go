package devinit

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func (c *Config) makeNativeApp() error {

	c.coreNativeClientId = "core-native-frontend"

	// env file
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("NODE_ENV=dev\n"))
	sb.WriteString(fmt.Sprintf("REACT_APP_OIDC_ISSUER=%s\n", HydraHost))
	sb.WriteString(fmt.Sprintf("REACT_APP_OAUTH_SCOPE=%s\n", "openid profile offline_access"))
	sb.WriteString(fmt.Sprintf("REACT_APP_OAUTH_CLIENT_ID=%s\n", c.coreNativeClientId))
	sb.WriteString(fmt.Sprintf("REACT_APP_SERVER_URL=https://localhost:8443\n"))

	envPath := path.Join(c.rootDir, "frontend", "apps", "native", ".env")
	fmt.Println(envPath)
	if err := ioutil.WriteFile(envPath, []byte(sb.String()), os.ModePerm); err != nil {
		return err
	}

	c.hydraClients = append(c.hydraClients, ClientConfig{
		ClientId: c.coreNativeClientId,
		RedirectUris: []string{
			"http://localhost:19006",
			"exp://127.0.0.1:19000",
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
