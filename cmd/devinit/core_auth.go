package devinit

import (
	"net"
	"path"
)

func (c *Config) makeCoreAuth() error {

	var err error

	c.coreAuthTlsKey, err = getOrCreatePrivateKey(path.Join(CoreAuthApiDir, "tls.key"))
	if err != nil {
		return err
	}

	c.coreAuthTlsCert, err = getOrCreateServerCert(
		path.Join(CoreAuthApiDir, "tls.crt"),
		c.coreAuthTlsKey,
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
