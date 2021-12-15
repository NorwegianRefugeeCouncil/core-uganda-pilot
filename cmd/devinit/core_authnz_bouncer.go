package devinit

import (
	"net"
	"path"
)

func (c *Config) makeCoreAuth() error {

	var err error

	c.coreAuthnzBouncerTlsKey, err = getOrCreatePrivateKey(path.Join(CoreAuthnzBouncerDir, "tls.key"))
	if err != nil {
		return err
	}

	c.coreAuthnzBouncerTlsCert, err = getOrCreateServerCert(
		path.Join(CoreAuthnzBouncerDir, "tls.crt"),
		c.coreAuthnzBouncerTlsKey,
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
