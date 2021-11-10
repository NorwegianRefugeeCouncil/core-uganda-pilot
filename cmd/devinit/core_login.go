package devinit

import (
	"net"
	"path"
)

func (c *Config) makeLogin() error {

	var err error

	c.loginBlockKey, err = getOrCreateRandomSecretStr(32, LoginDir, "block-key")
	if err != nil {
		return err
	}

	c.loginHashKey, err = getOrCreateRandomSecretStr(64, LoginDir, "hash-key")
	if err != nil {
		return err
	}

	c.loginTlsKey, err = getOrCreatePrivateKey(path.Join(LoginDir, "tls.key"))
	if err != nil {
		return err
	}

	c.loginTlsCert, err = getOrCreateServerCert(
		path.Join(LoginDir, "tls.crt"),
		c.loginTlsKey,
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
