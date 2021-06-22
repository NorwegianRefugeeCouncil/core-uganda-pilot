package teamorganizations

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/organizations"
	"github.com/nrc-no/core-kafka/pkg/parties/relationships"
	"github.com/nrc-no/core-kafka/pkg/parties/relationshiptypes"
	"github.com/nrc-no/core-kafka/pkg/teams"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

var RelationshipType = relationshiptypes.RelationshipType{
	ID:   "ea1f077a-30ba-4ffb-a1c4-a51669f82e37",
	Name: "TeamOrganization",
	Rules: []relationshiptypes.RelationshipTypeRule{
		{
			PartyTypeRule: relationshiptypes.PartyTypeRule{
				FirstPartyType:  teams.PartyType.ID,
				SecondPartyType: organizations.PartyType.ID,
			},
		},
	},
	FirstPartyRole:  "Is part of organization",
	SecondPartyRole: "has team",
	IsDirectional:   true,
}

var KampalaResponseTeamOrg = relationships.Relationship{
	ID:                 "cf49cbad-3c5e-4a23-8217-00b19ecf00a0",
	RelationshipTypeID: RelationshipType.ID,
	FirstParty:         teams.KampalaResponseTeam.ID,
	SecondParty:        organizations.NRC.ID,
}

var KampalaICLATeamOrg = relationships.Relationship{
	ID:                 "e10afd64-ab3a-4d56-92af-48eee0f555ab",
	RelationshipTypeID: RelationshipType.ID,
	FirstParty:         teams.KampalaICLATeam.ID,
	SecondParty:        organizations.NRC.ID,
}

var NairobiResponseTeamOrg = relationships.Relationship{
	ID:                 "09635283-4f43-449a-825c-dfcae1dd3aff",
	RelationshipTypeID: RelationshipType.ID,
	FirstParty:         teams.NairobiResponseTeam.ID,
	SecondParty:        organizations.NRC.ID,
}

var NairobiICLATeamOrg = relationships.Relationship{
	ID:                 "d23ee6ad-0662-46c9-9313-fb7f0bf4dedd",
	RelationshipTypeID: RelationshipType.ID,
	FirstParty:         teams.NairobiICLATeam.ID,
	SecondParty:        organizations.NRC.ID,
}

func Init(
	ctx context.Context,
	relationshipTypeStore *relationshiptypes.Store,
	relationshipStore *relationships.Store,
) error {

	logger := logrus.WithContext(ctx).WithField("logger", "teamorganizations.Init")
	logger.Infof("initializing TeamOrganization relationship type")

	if err := relationshipTypeStore.Create(ctx, &RelationshipType); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			logger.WithError(err).Errorf("failed to create TeamOrganization relationship type")
			return err
		}
		if err := relationshipTypeStore.Update(ctx, &RelationshipType); err != nil {
			logger.WithError(err).Errorf("failed to update TeamOrganization relationship type")
			return err
		}
	}

	for _, relationship := range []*relationships.Relationship{
		&KampalaResponseTeamOrg,
		&KampalaICLATeamOrg,
		&NairobiResponseTeamOrg,
		&NairobiICLATeamOrg,
	} {
		logger.Infof("creating team relationship")
		if err := upsertRelationship(ctx, relationshipStore, relationship); err != nil {
			logger.WithError(err).Errorf("failed to create team relationship")
			return err
		}
	}

	return nil

}

func upsertRelationship(ctx context.Context, relationshipStore *relationships.Store, relationship *relationships.Relationship) error {
	if err := relationshipStore.Create(ctx, relationship); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			return err
		}
		if err := relationshipStore.Update(ctx, relationship); err != nil {
			return err
		}
	}
	return nil
}
