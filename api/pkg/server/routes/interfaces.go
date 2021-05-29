package routes

import "k8s.io/apimachinery/pkg/util/sets"

type ListedPathProvider interface {
	ListedPaths() []string
}

type ListedPathProviders []ListedPathProvider

// ListedPaths unions and sorts the included paths.
func (p ListedPathProviders) ListedPaths() []string {
	ret := sets.String{}
	for _, provider := range p {
		for _, path := range provider.ListedPaths() {
			ret.Insert(path)
		}
	}

	return ret.List()
}
