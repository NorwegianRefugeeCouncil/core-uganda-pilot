package devinit

import (
	"fmt"
	"os"
	"path"
)

func (c *Config) makePostgres() error {
	var err error

	c.postgresRootPassword, err = getOrCreateRandomSecretStr(32, PostgresDir, "password")
	if err != nil {
		return err
	}
	c.postgresUsername = "postgres"
	if err := os.WriteFile(
		path.Join(PostgresDir, "env"),
		[]byte(fmt.Sprintf("POSTGRES_USER=%s\nPOSTGRES_PASSWORD=%s", c.postgresUsername, c.postgresRootPassword)),
		os.ModePerm,
	); err != nil {
		return err
	}
	return nil
}
