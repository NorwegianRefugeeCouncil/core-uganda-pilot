package devinit

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/nrc-no/core/pkg/store"
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"math/rand"
	"os"
	"path"
	"time"
)

var (
	RootDir              string
	CredsDir             string
	CertsDir             string
	HydraCredsDir        string
	IDPDir               string
	CoreDir              string
	CoreAppFrontendDir   string
	CoreAdminFrontendDir string
	CoreApiDir           string
	CoreAdminApiDir      string
	CoreAuthApiDir       string
	LoginDir             string
	RedisDir             string
	PostgresDir          string
	OIDCDir              string
	ProxyDir             string
)

const (
	OidcHost  string = "oidc.dev"
	LoginHost string = "core.dev/login"
	ProxyPort int    = 8443
)

type dbUser struct {
	username string
	password string
	database string
}

type Config struct {
	coreAdminApiBlockKey      string
	coreAdminApiHashKey       string
	coreAdminApiTlsCert       *x509.Certificate
	coreAdminApiTlsKey        *rsa.PrivateKey
	coreAdminFrontendClientId string
	coreAdminFrontendTlsCert  *x509.Certificate
	coreAdminFrontendTlsKey   *rsa.PrivateKey
	coreApiBlockKey           string
	coreApiHashKey            string
	coreApiTlsCert            *x509.Certificate
	coreApiTlsKey             *rsa.PrivateKey
	coreAppFrontendTlsCert    *x509.Certificate
	coreAppFrontendTlsKey     *rsa.PrivateKey
	coreDbName                string
	coreDbPassword            string
	coreDbUsername            string
	dbUsers                   []dbUser
	hydraCookieSecret         string
	hydraDbName               string
	hydraDbPassword           string
	hydraDbUsername           string
	hydraSystemSecret         string
	idpClientId               string
	idpClientSecret           string
	idpIssuer                 string
	idpRedirectUri            string
	loginBlockKey             string
	loginHashKey              string
	loginTlsCert              *x509.Certificate
	loginTlsKey               *rsa.PrivateKey
	oidcConfig                *OidcConfig
	hydraClients              []ClientConfig
	oidcTlsCert               *x509.Certificate
	oidcTlsKey                *rsa.PrivateKey
	postgresRootPassword      string
	postgresUsername          string
	redisPassword             string
	rootCa                    *x509.Certificate
	rootCaKey                 *rsa.PrivateKey
	rootCaKeyPath             string
	rootCaPath                string
	proxyTlsKey               *rsa.PrivateKey
	proxyTlsCert              *x509.Certificate
	rootDir                   string
	coreAuthTlsKey            *rsa.PrivateKey
	coreAuthTlsCert           *x509.Certificate
	coreAppFrontendClientId   string
}

func Init() error {
	_, err := createConfig()
	return err
}

func createConfig() (*Config, error) {

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	config := &Config{
		oidcConfig: &OidcConfig{},
		rootDir:    wd,
	}

	if err := config.makeRootCert(); err != nil {
		return nil, err
	}

	if err := config.makeProxyConfig(); err != nil {
		return nil, err
	}

	if err := config.makeRedis(); err != nil {
		return nil, err
	}

	if err := config.makePostgres(); err != nil {
		return nil, err
	}

	if err := config.makeIdp(); err != nil {
		return nil, err
	}

	if err := config.makeLogin(); err != nil {
		return nil, err
	}

	if err := config.makeHydra(); err != nil {
		return nil, err
	}

	if err := config.makeCoreApi(); err != nil {
		return nil, err
	}

	if err := config.makeCoreAdminApi(); err != nil {
		return nil, err
	}

	if err := config.makeAppFrontend(); err != nil {
		return nil, err
	}

	if err := config.makeAdminFrontend(); err != nil {
		return nil, err
	}

	if err := config.makeCoreAuth(); err != nil {
		return nil, err
	}

	if err := config.makeOidcConfig(); err != nil {
		return nil, err
	}

	if err := config.makeCore(); err != nil {
		return nil, err
	}

	if err := config.makePostgresInit(); err != nil {
		return nil, err
	}

	if err := config.makeHydraInit(); err != nil {
		return nil, err
	}

	return config, nil
}

var l *zap.Logger

func init() {
	var err error
	l, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}

func Bootstrap() error {

	l.Info("Bootstrapping...")

	envCfg, err := createConfig()
	if err != nil {
		return err
	}

	adminCli := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Host:    "localhost:4445",
		Schemes: []string{"https"},
	}).Admin

	jsonBytes, err := os.ReadFile(path.Join(HydraCredsDir, "clients.json"))
	if err != nil {
		return err
	}
	var clients HydraClients
	if err := json.Unmarshal(jsonBytes, &clients); err != nil {
		return err
	}

	for _, config := range clients.Clients {
		l.Info("Creating Hydra Client", zap.String("client_id", config.ClientId))
		hydraClient := &models.OAuth2Client{
			ClientID:                config.ClientId,
			ClientName:              "",
			ClientSecret:            config.ClientSecret,
			GrantTypes:              config.GrantTypes,
			RedirectUris:            config.RedirectUris,
			ResponseTypes:           config.ResponseTypes,
			Scope:                   config.Scope,
			TokenEndpointAuthMethod: config.TokenEndpointAuthMethod,
		}
		if _, err := adminCli.CreateOAuth2Client(&admin.CreateOAuth2ClientParams{
			Body:    hydraClient,
			Context: context.Background(),
		}); err != nil {
			if _, err := adminCli.UpdateOAuth2Client(&admin.UpdateOAuth2ClientParams{
				Body:    hydraClient,
				ID:      config.ClientId,
				Context: context.Background(),
			}); err != nil {
				return err
			}
		}
	}

	configBytes, err := os.ReadFile(path.Join(CoreDir, "config.yaml"))
	if err != nil {
		return err
	}
	var cfg = options.Options{}
	if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
		return err
	}
	factory, err := store.NewFactory(cfg.DSN)
	if err != nil {
		return err
	}

	orgStore := store.NewOrganizationStore(factory)

	existing, err := orgStore.List(context.Background())
	var orgId string
	if len(existing) == 0 {
		created, err := orgStore.Create(context.Background(), &types.Organization{
			Name: "Norwegian Refugee Council",
		})
		if err != nil {
			return err
		}
		orgId = created.ID
	} else {
		orgId = existing[0].ID
	}

	idpStore := store.NewIdentityProviderStore(factory)

	idps, err := idpStore.List(context.Background(), orgId, store.IdentityProviderListOptions{})
	if err != nil {
		return err
	}
	idp := &types.IdentityProvider{
		Name:           "Fake OIDC",
		OrganizationID: orgId,
		Domain:         "https://oidc.dev:8443",
		ClientID:       envCfg.idpClientId,
		ClientSecret:   envCfg.idpClientSecret,
		EmailDomain:    "nrc.no",
	}
	if len(idps) == 0 {
		_, err := idpStore.Create(context.Background(), idp, store.IdentityProviderCreateOptions{})
		if err != nil {
			return err
		}
	} else {
		idp.ID = idps[0].ID
		_, err := idpStore.Update(context.Background(), idp, store.IdentityProviderUpdateOptions{})
		if err != nil {
			return err
		}
	}

	return nil

}

func init() {
	rand.Seed(time.Now().UnixNano())
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	RootDir = cwd
	CredsDir = path.Join(RootDir, "creds")
	CertsDir = path.Join(RootDir, "certs")
	RedisDir = path.Join(CredsDir, "redis")
	ProxyDir = path.Join(CredsDir, "envoy")
	PostgresDir = path.Join(CredsDir, "postgres")
	HydraCredsDir = path.Join(CredsDir, "hydra")
	OIDCDir = path.Join(CredsDir, "oidc")
	IDPDir = path.Join(CredsDir, "nrc_idp")
	CoreDir = path.Join(CredsDir, "core")
	LoginDir = path.Join(CoreDir, "login")
	CoreAppFrontendDir = path.Join(CoreDir, "app_frontend")
	CoreAdminFrontendDir = path.Join(CoreDir, "admin_frontend")
	CoreApiDir = path.Join(CoreDir, "api")
	CoreAdminApiDir = path.Join(CoreDir, "admin_api")
	CoreAuthApiDir = path.Join(CoreDir, "auth")
}