package options

type CorsOptions struct {
	AllowedOrigins     []string `yaml:"allowed_origins"`
	AllowedMethods     []string `yaml:"allowed_methods"`
	AllowedHeaders     []string `yaml:"allowed_headers"`
	ExposedHeaders     []string `yaml:"exposed_headers"`
	AllowCredentials   bool     `yaml:"allow_credentials"`
	OptionsPassthrough bool     `yaml:"options_passthrough"`
	MaxAge             int      `yaml:"max_age"`
	Debug              bool     `yaml:"debug"`
	Enabled            bool     `yaml:"enabled"`
}

type SecretOptions struct {
	Hash  []string `yaml:"hash"`
	Block []string `yaml:"block"`
}

type URLOptions struct {
	Self string `yaml:"self"`
}

type OIDCOptions struct {
	Issuer       string   `yaml:"issuer"`
	ClientID     string   `yaml:"clientid"`
	ClientSecret string   `yaml:"clientsecret"`
	Scopes       []string `yaml:"scopes"`
}

type ServerOptions struct {
	Host    string        `yaml:"address"`
	Port    int           `yaml:"port"`
	Cors    CorsOptions   `yaml:"cors"`
	Secrets SecretOptions `yaml:"secrets"`
	URLs    URLOptions    `yaml:"urls"`
	Oidc    OIDCOptions   `yaml:"oidc"`
}

type ServeOptions struct {
	Public ServerOptions `yaml:"public"`
	Admin  ServerOptions `yaml:"admin"`
	Login  ServerOptions `yaml:"login"`
}

type Options struct {
	Serve ServeOptions `yaml:"serve"`
	DSN   string       `yaml:"dsn"`
}
