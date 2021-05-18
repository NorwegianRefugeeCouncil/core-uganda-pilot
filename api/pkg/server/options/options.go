package options

import (
	"github.com/nrc-no/core/api/pkg/server"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/storage/storagebackend"
)

type RecommendedOptions struct {
	Etcd *server.EtcdOptions
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
	if err := o.Etcd.ApplyTo(&config.Config); err != nil {
		return err
	}
	return nil
}

func (o *RecommendedOptions) Validate() []error {
	errors := []error{}
	errors = append(errors, o.Etcd.Validate()...)
	return errors
}
