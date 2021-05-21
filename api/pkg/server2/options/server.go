package options

import (
	"github.com/nrc-no/core/api/pkg/server2"
	"github.com/nrc-no/core/api/pkg/server2/endpoints/filters"
	"github.com/spf13/pflag"
	"net/http"
)

type Options struct {
	ListenAddress string
	StorageConfig MongoOptions
}

func (o *Options) Complete() error {
	return nil
}

func (o *Options) Validate() error {
	return nil
}

func (o *Options) Run(ch <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}
	server, err := config.Complete().New()
	if err != nil {
		return err
	}
	if err := server.Run(ch); err != nil {
		return err
	}
	return nil
}

func (c *Options) Config() (*server2.Config, error) {

	serverConfig := &server2.Config{
		ListenAddress: c.ListenAddress,
		BuildHandlerChainFunc: func(apiHandler http.Handler, config *server2.Config) http.Handler {
			handler := filters.WithRequestInfo(apiHandler)
			return handler
		},
	}

	if err := c.StorageConfig.ApplyTo(serverConfig); err != nil {
		return nil, err
	}

	serverConfig.CRDRestOptionsGetter = NewCRDRESTOptionsGetter(c.StorageConfig)

	return serverConfig, nil
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.StorageConfig.AddFlags(fs)
	fs.StringVar(&o.ListenAddress, "listen-address", o.ListenAddress,
		"http listen address")
}
