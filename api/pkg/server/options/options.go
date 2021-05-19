package options

import (
	"github.com/nrc-no/core/api/pkg/server"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	restclient "k8s.io/client-go/rest"
)

type RecommendedOptions struct {
	Etcd                 *server.EtcdOptions
	LoopbackClientConfig *restclient.Config
}

func NewRecommendedOptions(prefix string, codec runtime.Codec) *RecommendedOptions {
	o := &RecommendedOptions{
		Etcd: server.NewEtcdOptions(storagebackend.NewDefaultConfig(prefix, codec)),
	}
	return o
}

func (o *RecommendedOptions) AddFlags(fs *pflag.FlagSet) {
	o.Etcd.AddFlags(fs)
}

func (o *RecommendedOptions) ApplyTo(config *server.RecommendedConfig) error {
	c := &config.Config

	storageFactoryConfig := server.NewStorageFactoryConfig()
	completedStorageFactoryConfig, err := storageFactoryConfig.Complete(o.Etcd)
	if err != nil {
		return err
	}
	storageFactory, err := completedStorageFactoryConfig.New()
	if err != nil {
		return err
	}

	if err := o.Etcd.ApplyWithStorageFactoryTo(c, storageFactory); err != nil {
		return err
	}
	return nil
}

func (o *RecommendedOptions) Validate() []error {
	errors := []error{}
	errors = append(errors, o.Etcd.Validate()...)
	return errors
}
