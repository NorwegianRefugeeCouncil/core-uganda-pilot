package beneficiaries

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/expressions"
	i81n "github.com/nrc-no/core-kafka/pkg/i81n/api/v1"
	"github.com/nrc-no/core-kafka/pkg/parties/api"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/testhelpers"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStore(t *testing.T) {

	ctx := context.Background()
	mongoClient, err := testhelpers.NewMongoClient(ctx)

	if !assert.NoError(t, err) {
		return
	}
	store := NewStore(mongoClient)

	ID := uuid.NewV4().String()
	beneficiary := &Beneficiary{
		ID: ID,
		Attributes: map[string]api.AttributeValue{
			"firstName": api.AttributeValue{
				Attribute: attributes.Attribute{
					Name:        "firstName",
					ValueType:   expressions.ValueType{},
					SubjectType: api.BeneficiaryType,
					LongFormulation: []i81n.Translation{
						{
							Locale:      "en",
							Translation: "What is the beneficiary first name?",
						},
					},
					ShortFormulation: []i81n.Translation{
						{
							Locale:      "en",
							Translation: "First name",
						},
					},
					IsPersonallyIdentifiableInfo: false,
				},
				Value: "John",
			},
		},
	}
	err = store.Create(ctx, beneficiary)
	if !assert.NoError(t, err) {
		return
	}

	got, err := store.Get(ctx, ID)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, *beneficiary, *got)

}
