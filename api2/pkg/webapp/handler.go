package webapp

import (
	"fmt"
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
	"github.com/nrc-no/core-kafka/pkg/teams"
	"html/template"
	"os"
)

type Handler struct {
	template                  *template.Template
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
) (*Handler, error) {

	if len(options.TemplateDirectory) == 0 {
		options.TemplateDirectory = "pkg/webapp/templates/"
	}

	e, err := os.Executable()
	if err != nil {
		return nil, err
	}
	fmt.Println(e)

	t, err := template.ParseGlob(options.TemplateDirectory + "/*.gohtml")
	if err != nil {
		return nil, err
	}
	h := &Handler{
		template:                  t,
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
	}
	return h, nil
}
