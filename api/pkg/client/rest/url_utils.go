package rest

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/url"
	"path"
)

func DefaultServerURL(host, apiPath string, groupVersion schema.GroupVersion) (*url.URL, string, error) {
	if host == "" {
		return nil, "", fmt.Errorf("host must be a URL or a host:port pair")
	}
	base := host
	hostURL, err := url.Parse(base)
	if err != nil || hostURL.Scheme == "" || hostURL.Host == "" {
		scheme := "http://"
		hostURL, err = url.Parse(scheme + base)
		if err != nil {
			return nil, "", err
		}
		if hostURL.Path != "" && hostURL.Path != "/" {
			return nil, "", fmt.Errorf("host must be a URL or a host:port pair: %q", base)
		}
	}
	versionedAPIPath := DefaultVersionedAPIPath(apiPath, groupVersion)
	return hostURL, versionedAPIPath, nil
}

func DefaultVersionedAPIPath(apiPath string, groupVersion schema.GroupVersion) string {
	versionedAPIPath := path.Join("/", apiPath)
	if len(groupVersion.Group) > 0 {
		versionedAPIPath = path.Join(versionedAPIPath, groupVersion.Group, groupVersion.Version)
	} else {
		versionedAPIPath = path.Join(versionedAPIPath, groupVersion.Version)
	}
	return versionedAPIPath
}
