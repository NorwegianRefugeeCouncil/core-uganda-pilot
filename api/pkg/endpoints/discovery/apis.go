package discovery

import (
	discoveryv1 "github.com/nrc-no/core/api/pkg/apis/discovery/v1"
	listers "github.com/nrc-no/core/api/pkg/client/listers/discovery/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"net/http"
	"sort"
)

type apisHandler struct {
	codecs serializer.CodecFactory
	lister listers.APIServiceLister
}

func NewApisHandler(
	codecs serializer.CodecFactory,
	lister listers.APIServiceLister,
) http.Handler {
	return &apisHandler{
		codecs: codecs,
		lister: lister,
	}
}

func (r *apisHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	discoveryGroupList := &discoveryv1.APIGroupList{
		Groups: []discoveryv1.APIGroup{},
	}

	apiServices, err := r.lister.List(labels.Everything())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	grouped := groupByAPIGroup(apiServices)

	var groupNames []string
	for k, _ := range grouped {
		groupNames = append(groupNames, k)
	}
	sort.Strings(groupNames)

	for _, groupName := range groupNames {
		group := convertToDiscoveryAPIGroup(grouped[groupName])
		if group != nil {
			discoveryGroupList.Groups = append(discoveryGroupList.Groups, *group)
		}
	}

	responsewriters.WriteObjectNegotiated(
		r.codecs,
		negotiation.DefaultEndpointRestrictions,
		discoveryv1.SchemeGroupVersion,
		w,
		req,
		http.StatusOK,
		discoveryGroupList)
}

func groupByAPIGroup(apiServices []*discoveryv1.APIService) map[string][]*discoveryv1.APIService {
	ret := map[string][]*discoveryv1.APIService{}
	for _, service := range apiServices {
		if _, ok := ret[service.Spec.Group]; !ok {
			ret[service.Spec.Group] = []*discoveryv1.APIService{}
		}
		ret[service.Spec.Group] = append(ret[service.Spec.Group], service)
	}
	return ret
}

// convertToDiscoveryAPIGroup takes apiservices in a single group and returns a discovery compatible object.
// if none of the services are available, it will return nil.
func convertToDiscoveryAPIGroup(apiServices []*discoveryv1.APIService) *discoveryv1.APIGroup {

	var discoveryGroup *discoveryv1.APIGroup

	for _, apiService := range apiServices {
		// the first APIService which is valid becomes the default
		if discoveryGroup == nil {
			discoveryGroup = &discoveryv1.APIGroup{
				Name: apiService.Spec.Group,
				PreferredVersion: discoveryv1.GroupVersionForDiscovery{
					GroupVersion: apiService.Spec.Group + "/" + apiService.Spec.Version,
					Version:      apiService.Spec.Version,
				},
			}
		}

		discoveryGroup.Versions = append(discoveryGroup.Versions,
			discoveryv1.GroupVersionForDiscovery{
				GroupVersion: apiService.Spec.Group + "/" + apiService.Spec.Version,
				Version:      apiService.Spec.Version,
			},
		)
	}

	return discoveryGroup
}
