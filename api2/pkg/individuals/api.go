package individuals

import (
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"strings"
	"time"
)

type Individual struct {
	*parties.Party `json:",inline" bson:",inline"`
}

type IndividualList struct {
	Items []*Individual `json:"items" bson:"items"`
}

func NewIndividual(ID string) *Individual {
	return &Individual{
		Party: &parties.Party{
			ID: ID,
			PartyTypeIDs: []string{
				partytypes.IndividualPartyType.ID,
			},
			Attributes: map[string][]string{},
		},
	}
}

func (b *Individual) FindAge() *int {

	birthDate := b.FindAttributeValue(BirthDateAttribute.ID)
	if birthDate == nil {
		return nil
	}

	birthDateStrs, ok := birthDate.([]string)
	if !ok {
		return nil
	}

	if len(birthDateStrs) == 0 {
		return nil
	}

	birthDateStr := birthDateStrs[0]
	parsedDate, err := time.Parse("2006-01-02", birthDateStr)
	if err != nil {
		return nil
	}

	diff := time.Now().UTC().Sub(parsedDate)
	years := diff.Hours() / 24 / 365
	rounded := int(years)

	return &rounded

}

func (b *Individual) String() string {

	firstNames, hasFirstNames := b.Attributes[FirstNameAttribute.ID]
	lastNames, hasLastNames := b.Attributes[LastNameAttribute.ID]

	if hasFirstNames && hasLastNames {
		return strings.ToUpper(lastNames[0]) + ", " + firstNames[0]
	}

	return b.ID
}
