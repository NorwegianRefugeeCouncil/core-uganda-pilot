package devinit

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/dustinkirkland/golang-petname"
	"github.com/manifoldco/promptui"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/utils/files"
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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
	CoreFormsApiDir      string
	CoreAuthnzApiDir     string
	CoreAuthnzBouncerDir string
	LoginDir             string
	RedisDir             string
	PostgresDir          string
	OIDCDir              string
	ProxyDir             string
)

var (
	AdminURI      = "http://localhost:3001/app"
	AdminScope    = "openid email profile"
	PwaURI        = "http://localhost:3000/app"
	CoreLocalHost = "https://localhost:8443"

	WorkDir = ""

	// OidcIssuer is the oidc issuer of the fake-oidc-provider
	// This is initialized by the init() method
	OidcIssuer = ""

	// CoreHost is the core hostname
	// This is initialized by the init() method
	CoreHost = ""

	// HydraHost is the hydra hostname
	// This is initialized by the init() method
	HydraHost = ""
)

func initCredsDir() error {
	credsDir := path.Join(WorkDir, "creds")
	dirExists, err := files.DirectoryExists(credsDir)
	if err != nil {
		return err
	}
	if !dirExists {
		if err := os.MkdirAll(credsDir, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func getPetName() (string, error) {

	petNameFile := path.Join(WorkDir, "creds", "pet")
	petNameExists, err := files.FileExists(petNameFile)
	if err != nil {
		return "", err
	}
	if !petNameExists {
		prompt := promptui.Prompt{
			Label:   "Select Name for your Pet Computer:",
			Default: petname.Generate(3, "-"),
		}
		result, err := prompt.Run()
		if err != nil {
			return "", err
		}
		if err := os.WriteFile(petNameFile, []byte(result), os.ModePerm); err != nil {
			return "", err
		}
	}
	dir := filepath.Dir(petNameFile)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}
	petNameBytes, err := os.ReadFile(petNameFile)
	if err != nil {
		return "", err
	}
	petName := string(petNameBytes)
	return petName, nil
}

type dbUser struct {
	username string
	password string
	database string
}

type Config struct {
	coreAuthnzApiBlockKey     string
	coreAuthnzApiHashKey      string
	coreAuthnzApiTlsCert      *x509.Certificate
	coreAuthnzApiTlsKey       *rsa.PrivateKey
	coreAdminFrontendClientId string
	coreAdminFrontendTlsCert  *x509.Certificate
	coreAdminFrontendTlsKey   *rsa.PrivateKey
	coreApiBlockKey           string
	coreApiHashKey            string
	coreFormsApiTlsCert       *x509.Certificate
	coreFormsApiTlsKey        *rsa.PrivateKey
	coreAppFrontendClientId   string
	coreAppFrontendTlsCert    *x509.Certificate
	coreAppFrontendTlsKey     *rsa.PrivateKey
	coreAuthnzBouncerTlsCert  *x509.Certificate
	coreAuthnzBouncerTlsKey   *rsa.PrivateKey
	coreDbName                string
	coreDbPassword            string
	coreDbUsername            string
	coreNativeClientId        string
	dbUsers                   []dbUser
	hydraClients              []ClientConfig
	hydraCookieSecret         string
	hydraDbName               string
	hydraDbPassword           string
	hydraDbUsername           string
	hydraSystemSecret         string
	idpClientId               string
	idpClientSecret           string
	idpIssuer                 string
	loginBlockKey             string
	loginHashKey              string
	loginTlsCert              *x509.Certificate
	loginTlsKey               *rsa.PrivateKey
	oidcConfig                *OidcConfig
	oidcTlsCert               *x509.Certificate
	oidcTlsKey                *rsa.PrivateKey
	postgresRootPassword      string
	postgresUsername          string
	proxyTlsCert              *x509.Certificate
	proxyTlsKey               *rsa.PrivateKey
	redisPassword             string
	rootCa                    *x509.Certificate
	rootCaKey                 *rsa.PrivateKey
	rootCaKeyPath             string
	rootCaPath                string
	rootDir                   string
}

func Init() error {

	var err error
	l, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	WorkDir = wd

	if err := initCredsDir(); err != nil {
		panic(err)
	}

	petName, err := getPetName()
	if err != nil {
		return err
	}

	OidcIssuer = getOidcHost(petName)
	CoreHost = getServerHost(petName)
	HydraHost = fmt.Sprintf("%s/hydra", CoreHost)

	_, err = createConfig()
	return err
}

func StartTunnels() error {
	if err := Init(); err != nil {
		return err
	}

	petName, err := getPetName()
	if err != nil {
		return err
	}
	serverSubDomain := getServerTunnelName(petName)
	oidcSubDomain := getOidcTunnelName(petName)
	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() error {
		cmd := exec.Command("lt", "--port", "8444", "--subdomain", oidcSubDomain, "--local-https", "--allow-invalid-cert", "--print-requests")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	})
	g.Go(func() error {
		cmd := exec.Command("lt", "--port", "8443", "--subdomain", serverSubDomain, "--local-https", "--allow-invalid-cert", "--print-requests")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	})
	return g.Wait()
}

func getOidcTunnelName(petName string) string {
	return petName + "-oidc"
}

func getServerTunnelName(petName string) string {
	return petName + "-server"
}

func getOidcHost(petName string) string {
	return fmt.Sprintf("https://%s.loca.lt", getOidcTunnelName(petName))
}

func getServerHost(petName string) string {
	return fmt.Sprintf("https://%s.loca.lt", getServerTunnelName(petName))
}

func createConfig() (*Config, error) {
	config := &Config{
		oidcConfig: &OidcConfig{},
		rootDir:    WorkDir,
	}

	type errFunc func() error

	var funcs = []errFunc{
		config.makeRootCert,
		config.makeProxyConfig,
		config.makeRedis,
		config.makePostgres,
		config.makeIdp,
		config.makeLogin,
		config.makeHydra,
		config.makeCoreFormsApi,
		config.makeCoreAuthnzApi,
		config.makeAppFrontend,
		config.makeAdminFrontend,
		config.makeCoreAuth,
		config.makeNativeApp,
		config.makeOidcConfig,
		config.makeCore,
		config.makePostgresInit,
		config.makePostgres,
		config.makeHydraInit,
	}

	for _, f := range funcs {
		if err := f(); err != nil {
			return nil, err
		}
	}

	return config, nil
}

var l *zap.Logger

func Bootstrap() error {

	if err := Init(); err != nil {
		return err
	}

	l.Info("Bootstrapping...")

	envCfg, err := createConfig()
	if err != nil {
		return err
	}

	if err := createHydraClients(); err != nil {
		return err
	}

	factory, err := createDbFactory(err)
	if err != nil {
		return err
	}

	orgId, err := createOrganization(factory, err)
	if err != nil {
		return err
	}

	if err := createIdentityProviders(factory, err, orgId, envCfg); err != nil {
		return err
	}

	return nil

}

func createIdentityProviders(factory store.Factory, err error, orgId string, envCfg *Config) error {
	idpStore := store.NewIdentityProviderStore(factory)
	idps, err := idpStore.List(context.Background(), orgId, store.IdentityProviderListOptions{})
	if err != nil {
		return err
	}
	idp := &types.IdentityProvider{
		Name:           "Fake OIDC",
		OrganizationID: orgId,
		Domain:         OidcIssuer,
		ClientID:       envCfg.idpClientId,
		ClientSecret:   envCfg.idpClientSecret,
		EmailDomain:    "nrc.no",
		Scopes: 		"",
		Claim: 			types.Claim{Mappings: nil, Version: "0"},
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

func createOrganization(factory store.Factory, err error) (string, error) {
	orgStore := store.NewOrganizationStore(factory)
	existing, err := orgStore.List(context.Background())
	var orgId string
	if len(existing) == 0 {
		created, err := orgStore.Create(context.Background(), &types.Organization{
			Name: "Norwegian Refugee Council",
		})
		if err != nil {
			return "", err
		}
		orgId = created.ID
	} else {
		orgId = existing[0].ID
	}
	return orgId, nil
}

func createDbFactory(err error) (store.Factory, error) {
	coreConfig, err := parseCoreConfig(err)
	if err != nil {
		return nil, err
	}

	factory, err := store.NewFactory(coreConfig.DSN)
	if err != nil {
		return nil, err
	}
	return factory, nil
}

func parseCoreConfig(err error) (options.Options, error) {
	configBytes, err := os.ReadFile(path.Join(CoreDir, "config.yaml"))
	if err != nil {
		return options.Options{}, err
	}
	var cfg = options.Options{}
	if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
		return options.Options{}, err
	}
	return cfg, nil
}

func createHydraClients() error {

	adminCli := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Host:    "localhost:4445",
		Schemes: []string{"https"},
	}).Admin

	clients, err := parseHydraClients()
	if err != nil {
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
	return nil
}

func parseHydraClients() (HydraClients, error) {
	jsonBytes, err := os.ReadFile(path.Join(HydraCredsDir, "clients.json"))
	if err != nil {
		return HydraClients{}, err
	}
	var clients HydraClients
	if err := json.Unmarshal(jsonBytes, &clients); err != nil {
		return HydraClients{}, err
	}
	return clients, nil
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
	CoreFormsApiDir = path.Join(CoreDir, "forms_api")
	CoreAuthnzApiDir = path.Join(CoreDir, "authnz_api")
	CoreAuthnzBouncerDir = path.Join(CoreDir, "authnz_bouncer")
}
