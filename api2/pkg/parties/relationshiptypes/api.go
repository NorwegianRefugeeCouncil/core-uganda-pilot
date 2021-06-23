package relationshiptypes

type RelationshipTypeRule struct {
	PartyTypeRule *PartyTypeRule `json:",inline" bson:",inline"`
}

func (r RelationshipTypeRule) Mirror() RelationshipTypeRule {
	ret := RelationshipTypeRule{}
	if r.PartyTypeRule != nil {
		rev := r.PartyTypeRule.Mirror()
		ret.PartyTypeRule = &rev
	}
	return ret
}

type PartyTypeRule struct {
	FirstPartyType  string `json:"firstPartyType" bson:"firstPartyType"`
	SecondPartyType string `json:"secondPartyType" bson:"secondPartyType"`
}

func (p PartyTypeRule) Mirror() PartyTypeRule {
	return PartyTypeRule{
		FirstPartyType:  p.SecondPartyType,
		SecondPartyType: p.FirstPartyType,
	}
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

func (r *RelationshipType) Mirror() *RelationshipType {
	rules := r.Rules
	for i, rule := range rules {
		rules[i] = rule.Mirror()
	}
	return &RelationshipType{
		ID:              r.ID,
		IsDirectional:   r.IsDirectional,
		Name:            r.Name,
		FirstPartyRole:  r.SecondPartyRole,
		SecondPartyRole: r.FirstPartyRole,
		Rules:           rules,
	}
}
