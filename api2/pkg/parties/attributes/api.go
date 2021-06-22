package attributes

import "github.com/nrc-no/core-kafka/pkg/expressions"

type Attribute struct {
	ID                           string                 `json:"id" bson:"id"`
	Name                         string                 `json:"name" bson:"name"`
	ValueType                    expressions.ValueType  `json:"type" bson:"type"`
	PartyTypeIDs                 []string               `json:"partyTypeIds" bson:"partyTypeIds"`
	IsPersonallyIdentifiableInfo bool                   `json:"isPii" bson:"isPii"`
	Translations                 []AttributeTranslation `json:"translations" bson:"translations"`
}

type AttributeTranslation struct {
	Locale           string `json:"locale" bson:"locale"`
	LongFormulation  string `json:"longFormulation" bson:"longFormulation"`
	ShortFormulation string `json:"shortFormulation" bson:"shortFormulation"`
}

type AttributeList struct {
	Items []*Attribute `json:"items" bson:"items"`
}
