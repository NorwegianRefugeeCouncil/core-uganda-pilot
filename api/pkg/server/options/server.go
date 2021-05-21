package options

import (
	"context"
	"fmt"
	v1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	restclient "github.com/nrc-no/core/api/pkg/client/rest"
	filters2 "github.com/nrc-no/core/api/pkg/endpoints/filters"
	"github.com/nrc-no/core/api/pkg/server2"
	"github.com/spf13/pflag"
	"net"
	"net/http"
	"strconv"
)

type Options struct {
	BindAddress   net.IP
	BindPort      int
	Listener      net.Listener
	StorageConfig MongoOptions
}

func (o *Options) Complete() error {
	return nil
}

func (o *Options) Validate() error {
	return nil
}

func (o *Options) Run(ctx context.Context) error {

	config, err := o.Config(ctx)
	if err != nil {
		return err
	}

	server, err := config.Complete().New(ctx)
	if err != nil {
		return err
	}

	if err := server.Run(); err != nil {
		return err
	}

	return nil
}

func (c *Options) Config(ctx context.Context) (*server.Config, error) {

	serverConfig := &server.Config{
		ListenAddress: c.BindAddress,
		BuildHandlerChainFunc: func(apiHandler http.Handler, config *server.Config) http.Handler {
			handler := filters2.WithRequestInfo(apiHandler)
			return handler
		},
	}

	c.StorageConfig.StorageConfig.Codec = server.Codecs.LegacyCodec(v1.SchemeGroupVersion)

	if err := c.StorageConfig.ApplyTo(serverConfig); err != nil {
		return nil, err
	}

	addr := net.JoinHostPort(c.BindAddress.String(), strconv.Itoa(c.BindPort))
	lnConf := net.ListenConfig{}

	ln, port, err := CreateListener(ctx, addr, lnConf)
	if err != nil {
		return nil, err
	}
	c.Listener = ln
	c.BindPort = port

	serverConfig.LoopbackClientConfig = &restclient.Config{
		Host: "http://" + addr,
	}

	serverConfig.Listener = c.Listener
	serverConfig.CRDRestOptionsGetter = NewCRDRESTOptionsGetter(c.StorageConfig)

	return serverConfig, nil
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.StorageConfig.AddFlags(fs)
	fs.IPVar(&o.BindAddress, "bind-address", o.BindAddress, "http listen address")
	fs.IntVar(&o.BindPort, "bind-port", o.BindPort, "The port on which to serve traffic")
}

func CreateListener(ctx context.Context, addr string, config net.ListenConfig) (net.Listener, int, error) {
	ln, err := config.Listen(ctx, "tcp", addr)
	if err != nil {
		return nil, 0, err
	}
	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		ln.Close()
		return nil, 0, fmt.Errorf("invalid listen address: %q", ln.Addr().String())
	}
	return ln, tcpAddr.Port, nil
}
