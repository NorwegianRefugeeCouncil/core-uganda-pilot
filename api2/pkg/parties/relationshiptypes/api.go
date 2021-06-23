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
	FirstPartyTypeID  string `json:"firstPartyTypeId" bson:"firstPartyTypeId"`
	SecondPartyTypeID string `json:"secondPartyTypeId" bson:"secondPartyTypeId"`
}

func (p PartyTypeRule) Mirror() PartyTypeRule {
	return PartyTypeRule{
		FirstPartyTypeID:  p.SecondPartyTypeID,
		SecondPartyTypeID: p.FirstPartyTypeID,
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
