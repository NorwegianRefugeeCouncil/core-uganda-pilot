package api

import (
	"time"

	"github.com/nrc-no/core-kafka/pkg/parties/api"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
)

type Beneficiary struct {
	*api.Party `json:",inline" bson:",inline"`
}

type BeneficiaryList struct {
	Items []*Beneficiary `json:"items" bson:"items"`
}

func NewBeneficiary(ID string) *Beneficiary {
	return &Beneficiary{
		Party: &api.Party{
			ID: ID,
			PartyTypes: []string{
				partytypes.IndividualPartyType.ID,
			},
			Attributes: map[string][]string{},
		},
	}
}

func (b *Beneficiary) FindAge() *int {

	birthDate := b.FindAttributeValue(attributes.BirthDateAttribute.ID)
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

func (b *Beneficiary) String() string {
	return b.ID
}
