package teams

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/expressions"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

var PartyType = partytypes.PartyType{
	ID:        "dacd6e08-3e3d-495b-8655-ea1d8e822cf3",
	Name:      "Team",
	IsBuiltIn: true,
}

var TeamNameAttribute = attributes.Attribute{
	ID:                           "18f410a3-6fde-45ce-80c7-fc5d92b85870",
	Name:                         "teamName",
	ValueType:                    expressions.ValueType{},
	PartyTypes:                   []string{PartyType.ID},
	IsPersonallyIdentifiableInfo: false,
	Translations: []attributes.AttributeTranslation{
		{
			Locale:           "en",
			ShortFormulation: "Team name",
			LongFormulation:  "Team name",
		},
	},
}

var KampalaResponseTeam = Team{
	ID:   "83ff03b1-f96a-48ca-9041-f5340fc25d60",
	Name: "Kampala Response Team",
}

var KampalaICLATeam = Team{
	ID:   "aa5a6ca9-b590-40cb-9521-3d7713a5f37b",
	Name: "Kampala ICLA Team",
}

var NairobiResponseTeam = Team{
	ID:   "f9aa5076-0add-463b-b236-95e8b3e5b7b2",
	Name: "Nairobi Response Team",
}

var NairobiICLATeam = Team{
	ID:   "a7d2f2cc-df53-4bb2-a525-7d1d8a1725f1",
	Name: "Nairobi ICLA Team",
}

func Init(
	ctx context.Context,
	teamStore *Store,
	partyTypeStore *partytypes.Store,
	attributeStore *attributes.Store,
) error {
	logger := logrus.WithContext(ctx).WithField("logger", "teams.Init")

	logger.Infof("initializing TeamName attribute")
	if err := attributeStore.Create(ctx, &TeamNameAttribute); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			return err
		}
		if err := attributeStore.Update(ctx, &TeamNameAttribute); err != nil {
			return err
		}
	}

	logger.Infof("initializing Team party type")
	if err := partyTypeStore.Create(ctx, &PartyType); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			logger.WithError(err).Errorf("unable to create Team party type")
			return err
		}
		if err := partyTypeStore.Update(ctx, &PartyType); err != nil {
			logger.WithError(err).Errorf("unable to update Team party type")
			return err
		}
	}

	for _, team := range []*Team{
		&KampalaResponseTeam,
		&KampalaICLATeam,
		&NairobiResponseTeam,
		&NairobiICLATeam,
	} {
		logger.Infof("Creating team %s", team.Name)
		if err := upsertTeam(ctx, teamStore, team); err != nil {
			logger.WithError(err).Errorf("failed to upsert team")
			return err
		}
	}

	return nil
}

func upsertTeam(ctx context.Context, teamStore *Store, team *Team) error {
	if err := teamStore.Create(ctx, team); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			return err
		}
		if err := teamStore.Update(ctx, team); err != nil {
			return err
		}
	}
	return nil
}
