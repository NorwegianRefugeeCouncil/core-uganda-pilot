package devinit

import (
	"net"
	"path"
)

func (c *Config) makeCoreAdminApi() error {

	var err error

	c.coreAdminApiTlsKey, err = getOrCreatePrivateKey(path.Join(CoreAdminApiDir, "tls.key"))
	if err != nil {
		return err
	}

	c.coreAdminApiTlsCert, err = getOrCreateServerCert(
		path.Join(CoreAdminApiDir, "tls.crt"),
		c.coreAdminApiTlsKey,
		c.rootCa,
		c.rootCaKey,
		[]string{"localhost", "core.dev"},
		[]net.IP{net.IPv6loopback, net.ParseIP("127.0.0.1")},
	)
	if err != nil {
		return err
	}

	c.coreAdminApiBlockKey, err = getOrCreateRandomSecretStr(32, CoreAdminApiDir, "block-key")
	if err != nil {
		return err
	}
	c.coreAdminApiHashKey, err = getOrCreateRandomSecretStr(64, CoreAdminApiDir, "hash-key")
	if err != nil {
		return err
	}

	return nil
}
