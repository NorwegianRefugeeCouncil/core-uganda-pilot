package devinit

import (
	"net"
	"path"
)

func (c *Config) makeProxyConfig() error {
	var err error

	c.proxyTlsKey, err = getOrCreatePrivateKey(path.Join(ProxyDir, "tls.key"))
	if err != nil {
		return err
	}

	c.proxyTlsCert, err = getOrCreateServerCert(
		path.Join(ProxyDir, "tls.crt"),
		c.proxyTlsKey,
		c.rootCa,
		c.rootCaKey,
		[]string{"localhost", "core.dev", "oidc.dev"},
		[]net.IP{net.IPv6loopback, net.ParseIP("127.0.0.1")},
	)
	if err != nil {
		return err
	}

	return nil
}
