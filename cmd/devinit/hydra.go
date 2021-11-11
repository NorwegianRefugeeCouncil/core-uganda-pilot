package devinit

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"net"
	"os"
	"path"
)

func (c *Config) makeHydra() error {
	if err := c.makeHydraSecrets(); err != nil {
		return err
	}
	if err := c.makeHydraCerts(); err != nil {
		return err
	}
	return nil
}

func (c *Config) makeHydraSecrets() error {

	var err error

	c.hydraSystemSecret, err = getOrCreateRandomSecretStr(32, HydraCredsDir, "system-secret")
	if err != nil {
		return err
	}
	c.hydraCookieSecret, err = getOrCreateRandomSecretStr(32, HydraCredsDir, "cookie-secret")
	if err != nil {
		return err
	}
	c.hydraDbPassword, err = getOrCreateRandomSecretStr(32, HydraCredsDir, "db-password")
	if err != nil {
		return err
	}
	c.hydraDbName = "hydra"
	c.hydraDbUsername = "hydra"

	c.dbUsers = append(c.dbUsers, dbUser{
		username: c.hydraDbUsername,
		password: c.hydraDbPassword,
		database: c.hydraDbName,
	})

	hydraConfig := map[string]interface{}{
		"urls": map[string]interface{}{
			"self": map[string]interface{}{
				"public": fmt.Sprintf(HydraHost),
				"issuer": fmt.Sprintf(HydraHost),
			},
			"consent": fmt.Sprintf("%s/login/consent", CoreHost),
			"login":   fmt.Sprintf("%s/login/identify", CoreHost),
			"logout":  fmt.Sprintf("%s/login/logout", CoreHost),
		},
		"secrets": map[string]interface{}{
			"system": []string{
				c.hydraSystemSecret,
			},
			"cookie": []string{
				c.hydraCookieSecret,
			},
		},
		"dsn": fmt.Sprintf("postgres://%s:%s@db:5432/%s?sslmode=disable&max_conns=20&max_idle_conns=4",
			c.hydraDbUsername,
			c.hydraDbPassword,
			c.hydraDbName),
	}

	hydraConfigBytes, err := yaml.Marshal(hydraConfig)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(HydraCredsDir, "config.yaml"), hydraConfigBytes, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (c *Config) makeHydraCerts() error {
	hydraTlsKey, err := getOrCreatePrivateKey(path.Join(HydraCredsDir, "public-tls.key"))
	if err != nil {
		return err
	}

	_, err = getOrCreateServerCert(
		path.Join(HydraCredsDir, "public-tls.crt"),
		hydraTlsKey,
		c.rootCa,
		c.rootCaKey,
		[]string{"localhost", "core.dev"},
		[]net.IP{net.IPv6loopback, net.ParseIP("127.0.0.1")},
	)
	if err != nil {
		return err
	}

	hydraAdminTlsKey, err := getOrCreatePrivateKey(path.Join(HydraCredsDir, "admin-tls.key"))
	if err != nil {
		return err
	}

	_, err = getOrCreateServerCert(
		path.Join(HydraCredsDir, "admin-tls.crt"),
		hydraAdminTlsKey,
		c.rootCa,
		c.rootCaKey,
		[]string{"localhost", "core.dev"},
		[]net.IP{net.IPv6loopback, net.ParseIP("127.0.0.1")},
	)
	if err != nil {
		return err
	}

	return nil
}
