package devinit

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

func (c *Config) makeCore() error {

	var err error

	c.coreDbUsername = "core"
	c.coreDbName = "core"
	c.coreDbPassword, err = getOrCreateRandomSecretStr(32, CoreDir, "db-password")
	if err != nil {
		return err
	}

	c.dbUsers = append(c.dbUsers, dbUser{
		username: c.coreDbUsername,
		password: c.coreDbPassword,
		database: c.coreDbName,
	})

	coreConfig := map[string]interface{}{
		"serve": map[string]interface{}{
			"login": map[string]interface{}{
				"cache": map[string]interface{}{
					"redis": map[string]interface{}{
						"password": c.redisPassword,
					},
				},
				"secrets": map[string]interface{}{
					"hash": []string{
						c.loginHashKey,
					},
					"block": []string{
						c.loginBlockKey,
					},
				},
			},
		},
		"dsn": fmt.Sprintf("postgres://%s:%s@localhost:5433/%s?sslmode=disable", c.coreDbUsername, c.coreDbPassword, c.coreDbName),
		"hydra": map[string]interface{}{
			"admin": map[string]interface{}{
				"host":      "localhost:8443",
				"base_path": "hydra-admin/",
				"schemes":   []string{"https"},
			},
			"public": map[string]interface{}{
				"host":      "localhost:8443",
				"base_path": "hydra/",
				"schemes":   []string{"https"},
			},
		},
		"zanzibar": map[string]interface{}{
			"token": c.zanzibarToken,
			"prefix": c.zanzibarPrefix,
		},
	}

	yamlBytes, err := yaml.Marshal(coreConfig)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path.Join(CoreDir, "config.yaml"), yamlBytes, os.ModePerm); err != nil {
		return err
	}

	return nil
}
