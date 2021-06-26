package seed

import "github.com/nrc-no/core-kafka/pkg/apps/iam"

var KampalaResponseTeamOrg = iam.Relationship{
	ID:                 "cf49cbad-3c5e-4a23-8217-00b19ecf00a0",
	RelationshipTypeID: iam.TeamOrganizationRelationshipType.ID,
	FirstParty:         KampalaResponseTeam.ID,
	SecondParty:        NRC.ID,
}

var KampalaICLATeamOrg = iam.Relationship{
	ID:                 "e10afd64-ab3a-4d56-92af-48eee0f555ab",
	RelationshipTypeID: iam.TeamOrganizationRelationshipType.ID,
	FirstParty:         KampalaICLATeam.ID,
	SecondParty:        NRC.ID,
}

var NairobiResponseTeamOrg = iam.Relationship{
	ID:                 "09635283-4f43-449a-825c-dfcae1dd3aff",
	RelationshipTypeID: iam.TeamOrganizationRelationshipType.ID,
	FirstParty:         NairobiResponseTeam.ID,
	SecondParty:        NRC.ID,
}

var NairobiICLATeamOrg = iam.Relationship{
	ID:                 "d23ee6ad-0662-46c9-9313-fb7f0bf4dedd",
	RelationshipTypeID: iam.TeamOrganizationRelationshipType.ID,
	FirstParty:         NairobiICLATeam.ID,
	SecondParty:        NRC.ID,
}
