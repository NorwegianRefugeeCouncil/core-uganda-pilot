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

var cmsPath = "/apis/cms/v1/"
var iamPath = "/apis/iam/v1/"

var Endpoints = map[string]string{
	"cases":             cmsPath + "cases",
	"casetypes":         cmsPath + "casetypes",
	"comments":          cmsPath + "comments",
	"attributes":        iamPath + "attributes",
	"individuals":       iamPath + "individuals",
	"memberships":       iamPath + "memberships",
	"parties":           iamPath + "parties",
	"partytypes":        iamPath + "partyTypes",
	"relationships":     iamPath + "relationships",
	"relationshiptypes": iamPath + "relationshiptypes",
	"teams":             iamPath + "teams",
}
