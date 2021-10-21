package server

import (
	"github.com/nrc-no/core/pkg/attachments"
	"github.com/nrc-no/core/pkg/cms"
	"github.com/nrc-no/core/pkg/iam"
	"github.com/nrc-no/core/pkg/login"
	"github.com/nrc-no/core/pkg/storage"
	"github.com/nrc-no/core/pkg/webapp"
	"github.com/ory/hydra-client-go/client"
	"net/http"
)

type Server struct {
	MongoClientSrc    storage.MongoClientSrc
	WebAppServer      *webapp.Server
	HydraPublicClient *client.OryHydra
	HydraAdminClient  *client.OryHydra
	Router            http.Handler
	IAMServer         *iam.Server
	LoginServer       *login.Server
	CMSServer         *cms.Server
	AttachmentServer  *attachments.Server
}
