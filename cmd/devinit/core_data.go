package devinit

import (
	"net"
	"path"
)

func (c *Config) makeCoreDB() error {

	var err error

	c.coreDataTlsKey, err = getOrCreatePrivateKey(path.Join(CoreDataDir, "tls.key"))
	if err != nil {
		return err
	}

	c.coreDataTlsCert, err = getOrCreateServerCert(
		path.Join(CoreDataDir, "tls.crt"),
		c.coreDataTlsKey,
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
