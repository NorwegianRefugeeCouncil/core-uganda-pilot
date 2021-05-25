package options

import (
	"context"
	"fmt"
	v1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/api/pkg/authentication/oidc"
	restclient "github.com/nrc-no/core/api/pkg/client/rest"
	filters2 "github.com/nrc-no/core/api/pkg/endpoints/filters"
	"github.com/nrc-no/core/api/pkg/server"
	"github.com/spf13/pflag"
	mux2 "k8s.io/apiserver/pkg/server/mux"
	"net"
	"net/http"
	"strconv"
)

type Options struct {
	BindAddress      net.IP
	BindPort         int
	Listener         net.Listener
	StorageConfig    MongoOptions
	OidcIssuerUrl    string
	OidcClientID     string
	OidcClientSecret string
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

func (c *Options) WithOidcAuth(handler http.Handler) http.Handler {
	oidcHandler := oidc.NewOIDCHandler(server.Codecs, c.OidcClientID, c.OidcClientSecret, c.OidcIssuerUrl)
	mux := mux2.NewPathRecorderMux("oidc-auth")
	mux.Handle("/auth/login", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		oidcHandler.ServeLogin(writer, request)
	}))
	mux.Handle("/auth/logout", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		oidcHandler.ServeLogout(writer, request)
	}))
	mux.Handle("/auth/callback", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		oidcHandler.ServeCallback(writer, request)
	}))
	mux.NotFoundHandler(handler)
	return mux
}

func (c *Options) Config(ctx context.Context) (*server.Config, error) {

	serverConfig := &server.Config{
		ListenAddress: c.BindAddress,
		BuildHandlerChainFunc: func(apiHandler http.Handler, config *server.Config) http.Handler {
			handler := filters2.WithRequestInfo(apiHandler)
			handler = c.WithOidcAuth(handler)
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
	fs.StringVar(&o.OidcIssuerUrl, "oidc-issuer-url", o.OidcIssuerUrl, "Issuer url for the OIDC authenticator")
	fs.StringVar(&o.OidcClientID, "oidc-client-id", o.OidcClientID, "Client ID for the OIDC authenticator")
	fs.StringVar(&o.OidcClientSecret, "oidc-client-secret", o.OidcClientSecret, "Client secret for the OIDC authenticator")
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
