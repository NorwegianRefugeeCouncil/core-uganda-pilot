package rest

import (
	"github.com/nrc-no/core/apps/api/pkg/api/defaultscheme"
	"github.com/nrc-no/core/apps/api/pkg/apis/core"
	v1 "github.com/nrc-no/core/apps/api/pkg/apis/core/v1"
	formdefinitionsstorage "github.com/nrc-no/core/apps/api/pkg/registry/core/formdefinitions/storage"
	"github.com/nrc-no/core/apps/api/pkg/registry/generic"
	"github.com/nrc-no/core/apps/api/pkg/registry/rest"
	"github.com/nrc-no/core/apps/api/pkg/server"
	"github.com/nrc-no/core/apps/api/pkg/server/storage"
)

type StorageProvider struct{}

func (p StorageProvider) NewRESTStorage(cs storage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (server.APIGroupInfo, bool, error) {
	apiGroupInfo := server.NewDefaultAPIGroupInfo(core.GroupName, defaultscheme.Scheme, defaultscheme.ParameterCodec, defaultscheme.Codecs)

	if cs.VersionEnabled(v1.SchemeGroupVersion) {
		if storageMap, err := p.v1Storage(restOptionsGetter); err != nil {
			return server.APIGroupInfo{}, false, err
		} else {
			apiGroupInfo.VersionedResourcesStorageMap[v1.SchemeGroupVersion.Version] = storageMap
		}
	}

	return apiGroupInfo, true, nil
}

func (p StorageProvider) v1Storage(restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}
	formdefinitionsStorage, err := formdefinitionsstorage.NewREST(restOptionsGetter)
	if err != nil {
		return nil, err
	}
	storage["formdefinitions"] = formdefinitionsStorage
	return storage, nil
}

func (p StorageProvider) GroupName() string {
	return core.GroupName
}
