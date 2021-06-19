package memberships

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"github.com/nrc-no/core-kafka/pkg/parties/relationships"
	"github.com/nrc-no/core-kafka/pkg/parties/relationshiptypes"
	"github.com/nrc-no/core-kafka/pkg/teams"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type MembershipRelationship struct {
	relationships.Relationship
}

var RelationshipType = relationshiptypes.RelationshipType{
	ID:              "69fef57b-b37f-4803-a5fb-47e05282ac84",
	IsDirectional:   true,
	Name:            "teamMembership",
	FirstPartyRole:  "Is member of team",
	SecondPartyRole: "Has team member",
	Rules: []relationshiptypes.RelationshipTypeRule{
		{
			relationshiptypes.PartyTypeRule{
				FirstPartyType:  partytypes.IndividualPartyType.ID,
				SecondPartyType: teams.PartyType.ID,
			},
		},
	},
}

func Init(ctx context.Context, relationshipTypeStore *relationshiptypes.Store) error {
	logger := logrus.WithContext(ctx).WithField("logger", "teammembers.Init")
	logger.Infof("initializing TeamMembership relationship type")
	if err := relationshipTypeStore.Create(ctx, &RelationshipType); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			logger.WithError(err).Errorf("failed to create TeamMembership relationship type")
			return err
		}
		if err := relationshipTypeStore.Update(ctx, &RelationshipType); err != nil {
			logger.WithError(err).Errorf("failed to update TeamMembership relationship type")
			return err
		}
		return nil
	}
	return nil
}
