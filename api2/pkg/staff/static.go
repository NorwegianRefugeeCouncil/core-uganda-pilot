package staff

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/organizations"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"github.com/nrc-no/core-kafka/pkg/parties/relationshiptypes"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// RelationshipType represents the built-in Staff relationship type
var RelationshipType = relationshiptypes.RelationshipType{
	ID:              "53478121-23af-4ed8-a367-2e0de6d60271",
	Name:            "staff",
	FirstPartyRole:  "Is working for",
	SecondPartyRole: "Has staff",
	Rules: []relationshiptypes.RelationshipTypeRule{
		{
			PartyTypeRule: &relationshiptypes.PartyTypeRule{
				FirstPartyType:  partytypes.IndividualPartyType.ID,
				SecondPartyType: organizations.PartyType.ID,
			},
		},
	},
}

// Init makes sure that the built-in entities exist and are up to date in the database
func Init(ctx context.Context, relationshipTypeStore *relationshiptypes.Store) error {

	logger := logrus.WithContext(ctx).
		WithField("logger", "staff.Init").
		WithField("relationshipType", RelationshipType)

	logger.Infof("initializing Staff party type")
	if err := relationshipTypeStore.Create(ctx, &RelationshipType); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			logger.WithError(err).Errorf("unable to create Staff relationship type")
			return err
		}
		if err := relationshipTypeStore.Update(ctx, &RelationshipType); err != nil {
			logger.WithError(err).Errorf("unable to update Staff relationship type")
			return err
		}
		return nil
	}
	return nil
}
