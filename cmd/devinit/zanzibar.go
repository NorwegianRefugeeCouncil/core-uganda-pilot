package devinit

import (
	"github.com/nrc-no/core/pkg/zanzibar"
)

func (c *Config) makeZanzibarConfig() error {

	var err error
	c.zanzibarToken, err = getOrCreateRandomSecretStr(32, ZanzibarDir, "token")
	if err != nil {
		return err
	}

	var prefixErr error
	prefix, prefixErr := getOrCreateRandomSecretStr(32, ZanzibarDir, "prefix")
	if prefixErr != nil {
		return prefixErr
	}

	c.zanzibarConfig = zanzibar.ZanzibarClientConfig{
		Token: c.zanzibarToken,
		Prefix: prefix,
	}

	return nil
}
