package api

import (
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/parties/api"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"strings"
	"time"
)

type Beneficiary struct {
	*api.Party `json:",inline"`
}

type BeneficiaryList struct {
	Items []*Beneficiary `json:"items" bson:"items"`
}

func NewBeneficiary(ID string) *Beneficiary {
	return &Beneficiary{
		Party: &api.Party{
			ID: ID,
			PartyTypes: []string{
				partytypes.BeneficiaryPartyType.ID,
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

	parsedDate, err := time.Parse("2006-01-02", birthDate.(string))
	if err != nil {
		return nil
	}

	diff := time.Now().UTC().Sub(parsedDate)
	years := diff.Hours() / 24 / 365
	rounded := int(years)

	return &rounded

}

func (b *Beneficiary) String() string {
	var hasFirstName bool
	var hasLastName bool
	var firstName string
	var lastName string

	if b.HasAttribute(attributes.FirstNameAttribute.ID) {
		v, ok := b.GetAttribute(attributes.FirstNameAttribute.ID).(string)
		if ok {
			firstName = v
			hasFirstName = false
		} else {
			hasFirstName = false
		}
	}

	if b.HasAttribute(attributes.LastNameAttribute.ID) {
		v, ok := b.GetAttribute(attributes.LastNameAttribute.ID).(string)
		if ok {
			lastName = v
			hasLastName = false
		} else {
			hasLastName = false
		}
	}

	if hasFirstName && hasLastName {
		return fmt.Sprintf("%s, %s",
			strings.ToUpper(lastName),
			firstName,
		)
	}
	if hasFirstName {
		return firstName
	}
	if hasLastName {
		return lastName
	}
	return b.ID
}
