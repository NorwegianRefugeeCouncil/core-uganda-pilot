package relationships

import "time"

type Relationship struct {
	ID                  string     `json:"id" bson:"id"`
	RelationshipTypeID  string     `json:"relationshipTypeId" bson:"relationshipTypeId"`
	FirstParty          string     `json:"firstParty" bson:"firstParty"`
	SecondParty         string     `json:"secondParty" bson:"secondParty"`
	StartOfRelationship time.Time  `json:"startOfRelationship" bson:"startOfRelationship"`
	EndOfRelationship   *time.Time `json:"endOfRelationship" bson:"endOfRelationship"`
}

type RelationshipList struct {
	Items []*Relationship `json:"items" bson:"items"`
}
