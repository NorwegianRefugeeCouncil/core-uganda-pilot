package rest

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/http"
	"net/url"
	"strings"
)

type ContentConfig struct {
	AcceptContentTypes   string
	ContentType          string
	GroupVersion         *schema.GroupVersion
	NegotiatedSerializer runtime.NegotiatedSerializer
}

type Config struct {
	Host    string
	APIPath string
	ContentConfig
}

func NewRESTClient(baseURL *url.URL, versionedAPIPath string, config ClientContentConfig, client *http.Client) (*RESTClient, error) {
	if len(config.ContentType) == 0 {
		config.ContentType = "application/json"
	}

	base := *baseURL
	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}
	base.RawQuery = ""
	base.Fragment = ""

	return &RESTClient{
		base:             &base,
		versionedAPIPath: versionedAPIPath,
		content:          config,
		Client:           client,
	}, nil
}

func RESTClientFor(config *Config) (*RESTClient, error) {
	if config.GroupVersion == nil {
		return nil, fmt.Errorf("GroupVersion is required when initializing a RESTClient")
	}
	if config.NegotiatedSerializer == nil {
		return nil, fmt.Errorf("NegotiatedSerializer is required when initializing a RESTClient")
	}
	baseURL, versionedAPIPath, err := defaultServerURLFor(config)
	if err != nil {
		return nil, err
	}

	var gv schema.GroupVersion
	if config.GroupVersion != nil {
		gv = *config.GroupVersion
	}

	clientContent := ClientContentConfig{
		AcceptContentType: config.AcceptContentTypes,
		ContentType:       config.ContentType,
		GroupVersion:      gv,
		Negotiator:        runtime.NewClientNegotiator(config.NegotiatedSerializer, gv),
	}

	restClient, err := NewRESTClient(baseURL, versionedAPIPath, clientContent, http.DefaultClient)
	if err != nil {
		return nil, err
	}
	return restClient, nil
}

func defaultServerURLFor(config *Config) (*url.URL, string, error) {
	host := config.Host
	if host == "" {
		host = "localhost"
	}
	if config.GroupVersion != nil {
		return DefaultServerURL(host, config.APIPath, *config.GroupVersion)
	}
	return DefaultServerURL(host, config.APIPath, schema.GroupVersion{})
}
