package server

import (
	webapp2 "github.com/nrc-no/core/pkg/apps/webapp"
	"github.com/ory/hydra-client-go/client"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type Server struct {
	MongoClient       *mongo.Client
	WebAppHandler     *webapp2.Server
	HttpServer        *http.Server
	HydraPublicClient *client.OryHydra
	HydraAdminClient  *client.OryHydra
}
