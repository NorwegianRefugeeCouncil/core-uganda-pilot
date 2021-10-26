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

var (
	cmsPath         = "/apis/cms/v1"
	iamPath         = "/apis/iam/v1"
	attachmentsPath = "/apis/attachments/v1"
	documentsPath   = "/apis/documents/v1"
)

var (
	AttachmentsEndpoint                 = path.Join(attachmentsPath, "attachments")
	AttributesEndpoint                  = path.Join(iamPath, "attributes")
	BucketsEndpoint                     = path.Join(documentsPath, "buckets")
	CaseTypesEndpoint                   = path.Join(cmsPath, "casetypes")
	CasesEndpoint                       = path.Join(cmsPath, "cases")
	CommentsEndpoint                    = path.Join(cmsPath, "comments")
	CountriesEndpoint                   = path.Join(iamPath, "countries")
	DocumentsEndpoint                   = path.Join(documentsPath, "documents")
	FormDefinitionsEndpoint             = path.Join(documentsPath, "form_definitions")
	IdentificationDocumentTypesEndpoint = path.Join(iamPath, "identificationdocumenttypes")
	IdentificationDocumentsEndpoint     = path.Join(iamPath, "identificationdocuments")
	IndividualsEndpoint                 = path.Join(iamPath, "individuals")
	MembershipsEndpoint                 = path.Join(iamPath, "memberships")
	NationalitiesEndpoint               = path.Join(iamPath, "nationalities")
	PartiesEndpoint                     = path.Join(iamPath, "parties")
	PartyTypesEndpoint                  = path.Join(iamPath, "partytypes")
	RelationshipTypesEndpoint           = path.Join(iamPath, "relationshiptypes")
	RelationshipsEndpoint               = path.Join(iamPath, "relationships")
	TeamsEndpoint                       = path.Join(iamPath, "teams")
)
