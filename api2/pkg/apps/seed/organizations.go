package seed

import "github.com/nrc-no/core-kafka/pkg/apps/iam"

var NRC = iam.Organization{
	Party: &iam.Party{
		ID:           "e2dc39b1-ba29-4322-942c-46f02d7a586d",
		PartyTypeIDs: []string{iam.OrganizationPartyType.ID},
		Attributes: map[string][]string{
			iam.LegalNameAttribute.ID: {"NRC"},
		},
	},
}
