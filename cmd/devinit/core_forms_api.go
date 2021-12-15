package devinit

import (
	"net"
	"path"
)

func (c *Config) makeCoreFormsApi() error {

	var err error

	c.coreFormsApiTlsKey, err = getOrCreatePrivateKey(path.Join(CoreFormsApiDir, "tls.key"))
	if err != nil {
		return err
	}

	c.coreFormsApiTlsCert, err = getOrCreateServerCert(
		path.Join(CoreFormsApiDir, "tls.crt"),
		c.coreFormsApiTlsKey,
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
