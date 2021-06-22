package organizations

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

var LegalNameAttribute = attributes.Attribute{
	ID:   "7afb0744-c764-4c5b-9dc6-b341d9b320b4",
	Name: "legalName",
	Translations: []attributes.AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Legal Name",
			ShortFormulation: "Legal Name",
		},
	},
	PartyTypes: []string{
		PartyType.ID,
	},
}

var NRC = Organization{
	Party: &parties.Party{
		ID: "d0dc08b4-5e6b-461b-a444-09e02e69a8e1",
		PartyTypes: []string{
			PartyType.ID,
		},
		Attributes: map[string][]string{
			LegalNameAttribute.ID: {
				"NRC",
			},
		},
	},
}

func Init(
	ctx context.Context,
	partyTypeStore *partytypes.Store,
	attributeStore *attributes.Store,
	partyStore *parties.Store) error {

	logger := logrus.WithContext(ctx).WithField("logger", "organizations.Init")
	logger.Infof("initializing Organization party type")

	if err := partyTypeStore.Create(ctx, &PartyType); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			logger.WithError(err).Errorf("failed to create Organization party type")
			return err
		}
		if err := partyTypeStore.Update(ctx, &PartyType); err != nil {
			logger.WithError(err).Errorf("failed to update Organization party type")
			return err
		}
	}

	logger.Infof("initializing LegalName attribute")

	if err := attributeStore.Create(ctx, &LegalNameAttribute); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			logger.WithError(err).Errorf("failed to create LegalName attribute")
		}
		if err := attributeStore.Update(ctx, &LegalNameAttribute); err != nil {
			logger.WithError(err).Errorf("failed to update LegalName attribute")
			return err
		}
	}

	logger.Infof("initializing NRC organization")

	if err := partyStore.Create(ctx, NRC.Party); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			logger.WithError(err).Errorf("failed to create NRC organization")
		}
		if err := partyStore.Update(ctx, NRC.Party); err != nil {
			logger.WithError(err).Errorf("failed to update NRC organization")
			return err
		}
	}

	return nil
}
