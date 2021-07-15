package server

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

var (
	varPathToConfig = "config.path"
)

type Configuration struct {
	v *viper.Viper
}

func NewConfiguration() (*Configuration, error) {
	c := &Configuration{
		v: viper.New(),
	}
	c.v.SetConfigName("config")
	c.v.SetConfigType("yaml")
	viper.AddConfigPath("/etc/core/")
	viper.AddConfigPath("$HOME/.core")
	viper.AddConfigPath(".core")
	viper.AddConfigPath(".")

	c.v.AutomaticEnv()
	c.v.AddConfigPath(".")
	c.v.AddConfigPath("/core/config.yaml")
	c.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	c.v.SetTypeByDefaultValue(true)
	if err := c.v.ReadInConfig(); err != nil {
		return nil, err
	}
	c.v.WatchConfig()
	c.v.OnConfigChange(func(e fsnotify.Event) {
		logrus.WithField("file", e.Name).Warn("config file changed")
	})
	return c, nil
}

/**
--mongo-database=core --mongo-username=root --mongo-password=example --mongo-hosts=localhost:27017 --environment=Development --fresh=true --seed=true --hydra-admin-url=https://localhost:4445 --hydra-public-url=https://localhost:4444 --login-templates-directory=pkg/apps/login/templates --login-client-id=login --login-client-name=login --login-client-secret=somesecret --login-iam-host=localhost:9000 --login-iam-scheme=https --web-templates-directory=pkg/apps/webapp/templates --web-client-id=webapp --web-client-secret=webapp --web-client-name=webapp --web-iam-host=localhost:9000 --web-iam-scheme=https --web-cms-host=localhost:9000 --web-cms-scheme=https --listen-address=:9000 --base-url=https://localhost:9000 --tls-cert-path=certs/cert.pem --tls-key-path=certs/key.pem
*/
