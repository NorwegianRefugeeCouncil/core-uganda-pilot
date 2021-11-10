package devinit

import (
	"path"
)

func (c *Config) makeRootCert() error {
	c.rootCaKeyPath = path.Join(CertsDir, "ca.key")
	rootCaKey, err := getOrCreatePrivateKey(c.rootCaKeyPath)
	if err != nil {
		return err
	}
	c.rootCaKey = rootCaKey

	c.rootCaPath = path.Join(CertsDir, "ca.crt")
	rootCa, err := getOrCreateCaRoot(c.rootCaPath, rootCaKey)
	if err != nil {
		return err
	}
	c.rootCa = rootCa

	return nil
}
