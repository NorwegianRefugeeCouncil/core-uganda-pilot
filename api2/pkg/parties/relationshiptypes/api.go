package relationshiptypes

type RelationshipTypeRule struct {
	PartyTypeRule `json:",inline" bson:",inline"`
}

type PartyTypeRule struct {
	FirstPartyType  string `json:"firstPartyType" bson:"firstPartyType"`
	SecondPartyType string `json:"secondPartyType" bson:"secondPartyType"`
}

type RelationshipType struct {
	ID              string                 `json:"id" bson:"id"`
	IsDirectional   bool                   `json:"isDirectional" bson:"isDirectional"`
	Name            string                 `json:"name" bson:"name"`
	FirstPartyRole  string                 `json:"firstPartyRole" bson:"firstPartyRole"`
	SecondPartyRole string                 `json:"secondPartyRole" bson:"secondPartyRole"`
	Rules           []RelationshipTypeRule `json:"rules"`
}

type RelationshipTypeList struct {
	Items []*RelationshipType `json:"items" bson:"items"`
}

func (r *RelationshipType) reversed() *RelationshipType {
	return &RelationshipType{
		ID:              r.ID,
		IsDirectional:   r.IsDirectional,
		Name:            r.Name,
		FirstPartyRole:  r.SecondPartyRole,
		SecondPartyRole: r.FirstPartyRole,
		Rules: []RelationshipTypeRule{
			{
				PartyTypeRule: PartyTypeRule{
					FirstPartyType:  r.Rules[0].SecondPartyType,
					SecondPartyType: r.Rules[0].FirstPartyType,
				},
			},
		},
	}
}
