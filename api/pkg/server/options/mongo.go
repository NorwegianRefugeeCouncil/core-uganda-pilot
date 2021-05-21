package options

import (
	generic2 "github.com/nrc-no/core/api/pkg/registry/generic"
	"github.com/nrc-no/core/api/pkg/server2"
	store2 "github.com/nrc-no/core/api/pkg/store"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type MongoOptions struct {
	StorageConfig store2.Config
}

func (m *MongoOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&m.StorageConfig.Transport.ServerList, "mongo-servers", m.StorageConfig.Transport.ServerList,
		"List of mongo servers to connect with (scheme://p:port), comma separated")
	fs.StringVar(&m.StorageConfig.Transport.Username, "mongo-username", m.StorageConfig.Transport.Username,
		"username for mongo server")
	fs.StringVar(&m.StorageConfig.Transport.Password, "mongo-password", m.StorageConfig.Transport.Password,
		"password for mongo server")
	fs.StringVar(&m.StorageConfig.Transport.Database, "mongo-database", m.StorageConfig.Transport.Database,
		"mongo database name")
}

func (m *MongoOptions) ApplyTo(c *server.Config) error {
	if m == nil {
		return nil
	}
	c.RESTOptionsGetter = &SimpleRestOptionsFactory{
		Options: *m,
	}
	return nil
}

type SimpleRestOptionsFactory struct {
	Options MongoOptions
}

func (f *SimpleRestOptionsFactory) GetRESTOptions(resource schema.GroupResource) (generic2.RESTOptions, error) {
	ret := generic2.RESTOptions{
		StorageConfig: &f.Options.StorageConfig,
	}
	return ret, nil
}
