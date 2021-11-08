package devinit

import (
	"net"
	"path"
)

func (c *Config) makeCoreApi() error {

	var err error

	c.coreApiTlsKey, err = getOrCreatePrivateKey(path.Join(CoreApiDir, "tls.key"))
	if err != nil {
		return err
	}

	c.coreApiTlsCert, err = getOrCreateServerCert(
		path.Join(CoreApiDir, "tls.crt"),
		c.coreApiTlsKey,
		c.rootCa,
		c.rootCaKey,
		[]string{"localhost", "core.dev"},
		[]net.IP{net.IPv6loopback, net.ParseIP("127.0.0.1")},
	)

	if err != nil {
		return err
	}

	c.coreApiBlockKey, err = getOrCreateRandomSecretStr(32, CoreApiDir, "block-key")
	if err != nil {
		return err
	}
	c.coreApiHashKey, err = getOrCreateRandomSecretStr(64, CoreApiDir, "hash-key")
	if err != nil {
		return err
	}

	return nil
}
