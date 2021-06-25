package webapp

import (
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/auth"
	"github.com/nrc-no/core-kafka/pkg/cases/cases"
	"github.com/nrc-no/core-kafka/pkg/cases/casetypes"
	"github.com/nrc-no/core-kafka/pkg/individuals"
	"github.com/nrc-no/core-kafka/pkg/memberships"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"github.com/nrc-no/core-kafka/pkg/parties/relationships"
	"github.com/nrc-no/core-kafka/pkg/parties/relationshiptypes"
	"github.com/nrc-no/core-kafka/pkg/relationshipparties"
	"github.com/nrc-no/core-kafka/pkg/sessionmanager"
	"github.com/nrc-no/core-kafka/pkg/teams"
	"github.com/ory/hydra-client-go/client"
	"os"
)

type Handler struct {
	attributeClient           *attributes.Client
	individualClient          *individuals.Client
	relationshipTypeClient    *relationshiptypes.Client
	relationshipClient        *relationships.Client
	partyClient               *parties.Client
	partyTypeClient           *partytypes.Client
	caseTypeClient            *casetypes.Client
	caseClient                *cases.Client
	relationshipPartiesClient *relationshipparties.Client
	teamClient                *teams.Client
	membershipClient          *memberships.Client
	hydraAdminClient          *client.OryHydra
	hydraPublicClient         *client.OryHydra
	renderFactory             *RendererFactory
	sessionManager            sessionmanager.Store
	credentialsClient         *auth.CredentialsClient
	partyStore                *parties.Store
}

type Options struct {
	TemplateDirectory string
}

func NewHandler(
	options Options,
	AttributeClient *attributes.Client,
	IndividualClient *individuals.Client,
	RelationshipTypeClient *relationshiptypes.Client,
	RelationshipClient *relationships.Client,
	PartyClient *parties.Client,
	PartyTypeClient *partytypes.Client,
	CaseTypeClient *casetypes.Client,
	CaseClient *cases.Client,
	RelationshipPartiesClient *relationshipparties.Client,
	TeamClient *teams.Client,
	MembershipClient *memberships.Client,
	hydraAdminClient *client.OryHydra,
	hydraPublicClient *client.OryHydra,
	sessionManager sessionmanager.Store,
	credentialsClient *auth.CredentialsClient,
	partyStore *parties.Store,
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
		attributeClient:           AttributeClient,
		individualClient:          IndividualClient,
		relationshipTypeClient:    RelationshipTypeClient,
		relationshipClient:        RelationshipClient,
		partyClient:               PartyClient,
		partyTypeClient:           PartyTypeClient,
		caseTypeClient:            CaseTypeClient,
		caseClient:                CaseClient,
		relationshipPartiesClient: RelationshipPartiesClient,
		teamClient:                TeamClient,
		membershipClient:          MembershipClient,
		renderFactory:             renderFactory,
		hydraAdminClient:          hydraAdminClient,
		hydraPublicClient:         hydraPublicClient,
		sessionManager:            sessionManager,
		credentialsClient:         credentialsClient,
		partyStore:                partyStore,
	}
	return h, nil
}
