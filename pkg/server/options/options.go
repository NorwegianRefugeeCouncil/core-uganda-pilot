package options

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

type URLOptions struct {
	Self string `mapstructure:"self"`
}

type OIDCOptions struct {
	Issuer       string   `mapstructure:"issuer"`
	ClientID     string   `mapstructure:"clientid"`
	ClientSecret string   `mapstructure:"clientsecret"`
	Scopes       []string `mapstructure:"scopes"`
}

type ServerOptions struct {
	Host    string        `mapstructure:"address"`
	Port    int           `mapstructure:"port"`
	Cors    CorsOptions   `mapstructure:"cors"`
	Secrets SecretOptions `mapstructure:"secrets"`
	URLs    URLOptions    `mapstructure:"urls"`
	Oidc    OIDCOptions   `mapstructure:"oidc"`
}

type ServeOptions struct {
	Public ServerOptions `mapstructure:"public"`
	Admin  ServerOptions `mapstructure:"admin"`
	Login  ServerOptions `mapstructure:"login"`
}

type Options struct {
	Serve ServeOptions `mapstructure:"serve"`
	DSN   string       `mapstructure:"dsn"`
}
