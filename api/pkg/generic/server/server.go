package server

import (
	"github.com/boj/redistore"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client"
	"net/http"
	"path"
)

type GenericServerOptions struct {
	MongoClientFn     utils.MongoClientFn
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

type Endpoint = string

var (
	CasesEndpoint             Endpoint = path.Join(cmsPath, "cases")
	CaseTypesEndpoint         Endpoint = path.Join(cmsPath, "casetypes")
	CommentsEndpoint          Endpoint = path.Join(cmsPath, "comments")
	AttributesEndpoint        Endpoint = path.Join(iamPath, "attributes")
	IndividualsEndpoint       Endpoint = path.Join(iamPath, "individuals")
	MembershipsEndpoint       Endpoint = path.Join(iamPath, "memberships")
	PartiesEndpoint           Endpoint = path.Join(iamPath, "parties")
	PartyTypesEndpoint        Endpoint = path.Join(iamPath, "partytypes")
	RelationshipsEndpoint     Endpoint = path.Join(iamPath, "relationships")
	RelationshipTypesEndpoint Endpoint = path.Join(iamPath, "relationshiptypes")
	TeamsEndpoint             Endpoint = path.Join(iamPath, "teams")
	CountrysEndpoint		  Endpoint = path.Join(iamPath, "countrys")
)
