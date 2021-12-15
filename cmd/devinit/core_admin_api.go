package devinit

import (
	"net"
	"path"
)

func (c *Config) makeCoreAuthnzApi() error {

	var err error

	c.coreAuthnzApiTlsKey, err = getOrCreatePrivateKey(path.Join(CoreAuthnzApiDir, "tls.key"))
	if err != nil {
		return err
	}

	c.coreAuthnzApiTlsCert, err = getOrCreateServerCert(
		path.Join(CoreAuthnzApiDir, "tls.crt"),
		c.coreAuthnzApiTlsKey,
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
