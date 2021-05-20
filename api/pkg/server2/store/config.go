package store

import "k8s.io/apimachinery/pkg/runtime"

type TransportConfig struct {
	ServerList []string
	Database   string
	Username   string
	Password   string
}

type Config struct {
	Codec     runtime.Codec
	Transport TransportConfig
}
