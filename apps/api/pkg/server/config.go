package server

import (
  "fmt"
  "k8s.io/apimachinery/pkg/runtime"
  "k8s.io/apimachinery/pkg/runtime/serializer"
  "k8s.io/apimachinery/pkg/util/sets"
  "net/http"
)

const (
  // APIGroupPrefix is where non-legacy API group will be located.
  APIGroupPrefix = "/apis"
)

type Config struct {
  Serializer             runtime.NegotiatedSerializer
  BuildHandlerChainFunc  func(apiHandler http.Handler, c *Config) (secure http.Handler)
  LegacyAPIGroupPrefixes sets.String
}

func NewConfig(codecs serializer.CodecFactory) *Config {
  return &Config{
    Serializer:             codecs,
    LegacyAPIGroupPrefixes: sets.NewString("/api"),
  }
}

func (c *Config) New(name string) (*Server, error) {
  if c.Serializer == nil {
    return nil, fmt.Errorf("config.Serializer == nil")
  }

  apiServerHandler := NewAPIServerHandler(name, func(apiHandler http.Handler) http.Handler {
    return apiHandler
  })

  s := &Server{
    Handler:                apiServerHandler,
    legacyAPIGroupPrefixes: c.LegacyAPIGroupPrefixes,
  }

  return s, nil

}
