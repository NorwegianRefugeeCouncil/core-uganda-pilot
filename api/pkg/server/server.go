package server

import (
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/attachments"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/apps/login"
	"github.com/nrc-no/core/pkg/apps/webapp"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client"
)

type Server struct {
	MongoClientFn     utils.MongoClientFn
	WebAppServer      *webapp.Server
	HydraPublicClient *client.OryHydra
	HydraAdminClient  *client.OryHydra
	Router            *mux.Router
	IAMServer         *iam.Server
	LoginServer       *login.Server
	CMSServer         *cms.Server
	AttachmentServer  *attachments.Server
}
