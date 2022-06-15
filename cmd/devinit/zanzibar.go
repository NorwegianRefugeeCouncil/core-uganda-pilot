package devinit

import (
	"fmt"
	"os"
	"path"
)


func (c *Config) makeZanzibarClient() error {

	var err error
	c.zanzibarToken, err = getOrCreateRandomSecretStr(32, ZanzibarDir, "token")
	if err != nil {
		return err
	}

	if err := os.WriteFile(
		path.Join(RedisDir, "env"),
		[]byte(fmt.Sprintf("ZANZIBAR_TOKEN=%s", c.zanzibarToken)),
		os.ModePerm,
	); err != nil {
		return err
	}

	return nil
}
