package devinit

import (
	"fmt"
	"os"
	"path"
)

func (c *Config) makeRedis() error {
	var err error
	c.redisPassword, err = getOrCreateRandomSecretStr(32, RedisDir, "password")
	if err != nil {
		return err
	}

	if err := os.WriteFile(
		path.Join(RedisDir, "env"),
		[]byte(fmt.Sprintf("REDIS_PASSWORD=%s", c.redisPassword)),
		os.ModePerm,
	); err != nil {
		return err
	}

	return nil
}
