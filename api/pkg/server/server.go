package server

import (
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/apps/login"
	"github.com/nrc-no/core/pkg/apps/webapp"
	"github.com/ory/hydra-client-go/client"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	MongoClient       *mongo.Client
	WebAppServer      *webapp.Server
	HydraPublicClient *client.OryHydra
	HydraAdminClient  *client.OryHydra
	Router            *mux.Router
	IAMServer         *iam.Server
	LoginServer       *login.Server
	CMSServer         *cms.Server
}
