package server

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gomodule/redigo/redis"
	"github.com/ory/hydra-client-go/client"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type GenericServerOptions struct {
	MongoClient       *mongo.Client
	MongoDatabase     string
	Environment       string
	HydraAdminClient  *client.OryHydra
	HydraPublicClient *client.OryHydra
	HydraHTTPClient   *http.Client
	RedisPool         *redis.Pool
	OidcProvider      *oidc.Provider
}
