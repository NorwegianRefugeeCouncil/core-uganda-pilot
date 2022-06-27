package devinit

func (c *Config) makeZanzibarConfig() error {

	var err error
	c.zanzibarToken, err = getOrCreateRandomSecretStr(32, ZanzibarDir, "token")
	if err != nil {
		return err
	}

	c.zanzibarPrefix, err = getOrCreateRandomSecretStr(32, ZanzibarDir, "prefix")
	if err != nil {
		return err
	}

	return nil
}
