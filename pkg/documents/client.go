package documents

import (
	"github.com/nrc-no/core/pkg/rest"
	"strings"
)

type Interface interface {
	Documents() Documents
	Buckets() Buckets
}

type client struct {
	c *rest.Client
}

func (c *client) Documents() Documents {
	return NewDocumentsClient(c.c)
}

func (c *client) Buckets() Buckets {
	return NewBucketsClient(c.c)
}

func NewFromClient(c *rest.Client) Interface {
	return &client{c: c}
}

func NewFromConfig(config *rest.Config) Interface {
	return &client{c: rest.NewClient(config)}
}

func normaliseKey(key string) string {
	if strings.HasPrefix(key, "/") {
		key = strings.TrimPrefix(key, "/")
	}
	if strings.HasSuffix(key, "/") {
		key = strings.TrimSuffix(key, "/")
	}
	return key
}
