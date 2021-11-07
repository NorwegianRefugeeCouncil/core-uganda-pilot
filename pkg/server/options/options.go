package options

import (
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
)

type CorsOptions struct {
	AllowedOrigins     []string `mapstructure:"allowed_origins"`
	AllowedMethods     []string `mapstructure:"allowed_methods"`
	AllowedHeaders     []string `mapstructure:"allowed_headers"`
	ExposedHeaders     []string `mapstructure:"exposed_headers"`
	AllowCredentials   bool     `mapstructure:"allow_credentials"`
	OptionsPassthrough bool     `mapstructure:"options_passthrough"`
	MaxAge             int      `mapstructure:"max_age"`
	Debug              bool     `mapstructure:"debug"`
	Enabled            bool     `mapstructure:"enabled"`
}

type SecretOptions struct {
	Hash  []string `mapstructure:"hash"`
	Block []string `mapstructure:"block"`
}

type CacheOptions struct {
	Cookie *CookieOptions `mapstructure:"cookie,omitempty"`
	Redis  *RedisOptions  `mapstructure:"redis,omitempty"`
}

type CookieOptions struct {
}

type RedisOptions struct {
	MaxIdleConnections int    `mapstructure:"max_idle_connections"`
	Address            string `mapstructure:"address"`
	Password           string `mapstructure:"password"`
	MaxLength          int    `mapstructure:"max_length"`
}

type URLOptions struct {
	Self string `mapstructure:"self"`
}

type OIDCOptions struct {
	Issuer       string   `mapstructure:"issuer"`
	ClientID     string   `mapstructure:"client_id"`
	ClientSecret string   `mapstructure:"client_secret"`
	Scopes       []string `mapstructure:"scopes"`
	RedirectURI  string   `mapstructure:"redirect_uri"`
}

type ServerOptions struct {
	Host    string        `mapstructure:"address"`
	Port    int           `mapstructure:"port"`
	Cors    CorsOptions   `mapstructure:"cors"`
	Secrets SecretOptions `mapstructure:"secrets"`
	URLs    URLOptions    `mapstructure:"urls"`
	Oidc    OIDCOptions   `mapstructure:"oidc"`
	Cache   CacheOptions  `mapstructure:"cache"`
	TLS     TLSOptions    `mapstructure:"tls"`
}

type ServeOptions struct {
	Public ServerOptions `mapstructure:"public"`
	Admin  ServerOptions `mapstructure:"admin"`
	Login  ServerOptions `mapstructure:"login"`
}

type CertOptions struct {
	Path string `mapstructure:"path"`
}

type TLSOptions struct {
	Enabled bool        `mapstructure:"enabled"`
	Cert    CertOptions `mapstructure:"cert"`
	Key     CertOptions `mapstructure:"key"`
}

type LogOptions struct {
	Level string `mapstructure:"level"`
}

type HydraEndpoint struct {
	Schemes  []string `mapstructure:"schemes"`
	Host     string   `mapstructure:"host"`
	BasePath string   `mapstructure:"base_path"`
}

func (h HydraEndpoint) AdminClient() admin.ClientService {
	return client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Host:     h.Host,
		BasePath: h.BasePath,
		Schemes:  h.Schemes,
	}).Admin
}

type HydraOptions struct {
	Admin HydraEndpoint `mapstructure:"admin"`
}

type Options struct {
	Serve ServeOptions `mapstructure:"serve"`
	DSN   string       `mapstructure:"dsn"`
	Log   LogOptions   `mapstructure:"log"`
	Hydra HydraOptions `mapstructure:"hydra"`
}
