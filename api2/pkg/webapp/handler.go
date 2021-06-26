package webapp

import (
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/apps/iam"
	"github.com/nrc-no/core-kafka/pkg/auth"
	"github.com/nrc-no/core-kafka/pkg/cases/cases"
	"github.com/nrc-no/core-kafka/pkg/cases/casetypes"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"github.com/nrc-no/core-kafka/pkg/relationshipparties"
	"github.com/nrc-no/core-kafka/pkg/sessionmanager"
	"github.com/ory/hydra-client-go/client"
	"os"
)

type Handler struct {
	caseTypeClient            *casetypes.Client
	caseClient                *cases.Client
	relationshipPartiesClient *relationshipparties.Client
	hydraAdminClient          *client.OryHydra
	hydraPublicClient         *client.OryHydra
	renderFactory             *RendererFactory
	sessionManager            sessionmanager.Store
	credentialsClient         *auth.CredentialsClient
	partyStore                *parties.Store
	iam                       *iam.ClientSet
}

type Options struct {
	TemplateDirectory string
}

func NewHandler(
	options Options,
	CaseTypeClient *casetypes.Client,
	CaseClient *cases.Client,
	RelationshipPartiesClient *relationshipparties.Client,
	hydraAdminClient *client.OryHydra,
	hydraPublicClient *client.OryHydra,
	sessionManager sessionmanager.Store,
	credentialsClient *auth.CredentialsClient,
	partyStore *parties.Store,
	iamClient *iam.ClientSet,
) (*Handler, error) {

	renderFactory, err := NewRendererFactory(options.TemplateDirectory)
	if err != nil {
		return nil, err
	}

	e, err := os.Executable()
	if err != nil {
		return nil, err
	}
	fmt.Println(e)

	h := &Handler{
		caseTypeClient:            CaseTypeClient,
		caseClient:                CaseClient,
		relationshipPartiesClient: RelationshipPartiesClient,
		renderFactory:             renderFactory,
		hydraAdminClient:          hydraAdminClient,
		hydraPublicClient:         hydraPublicClient,
		sessionManager:            sessionManager,
		credentialsClient:         credentialsClient,
		partyStore:                partyStore,
	}
	return h, nil
}
