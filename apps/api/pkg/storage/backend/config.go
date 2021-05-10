package backend

import (
	"github.com/nrc-no/core/apps/api/pkg/runtime"
)

type Config struct {
	Codec           runtime.Codec
	EncodeVersioner runtime.GroupVersioner
	Prefix          string
	Transport       TransportConfig
}

type TransportConfig struct {
	ServerList []string
}
