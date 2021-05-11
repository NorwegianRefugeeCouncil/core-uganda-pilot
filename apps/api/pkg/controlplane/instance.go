package controlplane

import (
	"fmt"
	v1 "github.com/nrc-no/core/apps/api/pkg/apis/core/v1"
	formdefinitionsstorage "github.com/nrc-no/core/apps/api/pkg/registry/core/formdefinitions/rest"
	"github.com/nrc-no/core/apps/api/pkg/registry/generic"
	"github.com/nrc-no/core/apps/api/pkg/server"
	"github.com/nrc-no/core/apps/api/pkg/server/storage"
	"github.com/sirupsen/logrus"
)

func DefaultAPIResourceConfigSource() *storage.ResourceConfig {
	ret := storage.NewResourceConfig()

	ret.EnableVersions(
		v1.SchemeGroupVersion,
	)

	return ret
}

type Config struct {
	GenericConfig *server.Config
	ExtraConfig   ExtraConfig
}

type completedConfig struct {
	GenericConfig server.CompletedConfig
	ExtraConfig   *ExtraConfig
}

type CompletedConfig struct {
	*completedConfig
}

type ExtraConfig struct {
	APIResourceConfigSource storage.APIResourceConfigSource
	StorageFactory          storage.StorageFactory
}

type Instance struct {
	GenericAPIServer *server.Server
}

func (c *Config) Complete() CompletedConfig {

	cfg := completedConfig{
		c.GenericConfig.Complete(),
		&c.ExtraConfig,
	}

	return CompletedConfig{
		&cfg,
	}

}

// RESTStorageProvider is a factory type for REST storage.
type RESTStorageProvider interface {
	GroupName() string
	NewRESTStorage(apiResourceConfigSource storage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (server.APIGroupInfo, bool, error)
}

func (c completedConfig) New(delegationTarget server.DelegationTarget) (*Instance, error) {
	s, err := c.GenericConfig.New(delegationTarget)
	if err != nil {
		return nil, err
	}

	m := &Instance{
		GenericAPIServer: s,
	}

	restStorageProviders := []RESTStorageProvider{
		formdefinitionsstorage.StorageProvider{},
	}
	if err := m.InstallAPIs(c.ExtraConfig.APIResourceConfigSource, c.GenericConfig.RESTOptionsGetter, restStorageProviders...); err != nil {
		return nil, err
	}

	return m, nil

}

func (m *Instance) InstallAPIs(apiResourceConfigSource storage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter, restStorageProviders ...RESTStorageProvider) error {

	apiGroupsInfo := []*server.APIGroupInfo{}

	for _, restStorageBuilder := range restStorageProviders {
		groupName := restStorageBuilder.GroupName()
		if !apiResourceConfigSource.AnyVersionForGroupEnabled(groupName) {
			logrus.Info("skipping disabled api group %q", groupName)
			continue
		}

		apiGroupInfo, enabled, err := restStorageBuilder.NewRESTStorage(apiResourceConfigSource, restOptionsGetter)
		if err != nil {
			return fmt.Errorf("problem initializing api group %s: %v", groupName, err)
		}

		if !enabled {
			logrus.Warnf("api group %s is not enabled, skipping", groupName)
		}

		logrus.Infof("enabling api group %q", groupName)

		apiGroupsInfo = append(apiGroupsInfo, &apiGroupInfo)

	}

	if err := m.GenericAPIServer.InstallAPIGroups(apiGroupsInfo...); err != nil {
		return fmt.Errorf("error in registering group versions: %v", err)
	}

	return nil

}
