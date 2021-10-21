package server

import (
	"github.com/boj/redistore"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/nrc-no/core/pkg/storage"
	"github.com/ory/hydra-client-go/client"
	"net/http"
	"path"
)

type GenericServerOptions struct {
	MongoClientSrc    storage.MongoClientSrc
	MongoDatabase     string
	Environment       string
	HydraAdminClient  *client.OryHydra
	HydraPublicClient *client.OryHydra
	HydraHTTPClient   *http.Client
	OidcProvider      *oidc.Provider
	RedisStore        *redistore.RediStore
}

var cmsPath = "/apis/cms/v1"
var iamPath = "/apis/iam/v1"
var attachmentsPath = "/apis/attachments/v1"

type Endpoint = string

var (
	AttachmentsEndpoint                 = path.Join(attachmentsPath, "attachments")
	CasesEndpoint                       = path.Join(cmsPath, "cases")
	CaseTypesEndpoint                   = path.Join(cmsPath, "casetypes")
	CommentsEndpoint                    = path.Join(cmsPath, "comments")
	AttributesEndpoint                  = path.Join(iamPath, "attributes")
	IdentificationDocumentTypesEndpoint = path.Join(iamPath, "identificationdocumenttypes")
	IdentificationDocumentsEndpoint     = path.Join(iamPath, "identificationdocuments")
	IndividualsEndpoint                 = path.Join(iamPath, "individuals")
	MembershipsEndpoint                 = path.Join(iamPath, "memberships")
	PartiesEndpoint                     = path.Join(iamPath, "parties")
	PartyTypesEndpoint                  = path.Join(iamPath, "partytypes")
	RelationshipsEndpoint               = path.Join(iamPath, "relationships")
	RelationshipTypesEndpoint           = path.Join(iamPath, "relationshiptypes")
	TeamsEndpoint                       = path.Join(iamPath, "teams")
	CountriesEndpoint                   = path.Join(iamPath, "countries")
	NationalitiesEndpoint               = path.Join(iamPath, "nationalities")
)
