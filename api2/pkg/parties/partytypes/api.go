package partytypes

type PartyType struct {
	ID        string `json:"id" bson:"id"`
	Name      string `json:"name" bson:"name"`
	IsBuiltIn bool   `json:"isBuiltIn" bson:"isBuiltIn"`
}

type PartyTypeSchema struct {
	ID    string                `json:"id" bson:"id"`
	Name  string                `json:"name" bson:"name"`
	Nodes []PartyTypeSchemaNode `json:"nodes" bson:"nodes"`
}

type PartyTypeSchemaList struct {
	Items []*PartyTypeSchema `json:"items" bson:"items"`
}

type PartyTypeSchemaNode struct {
	ID          string                `json:"id" bson:"id"`
	PartyTypeID string                `json:"partyTypeID" bson:"partyTypeID"`
	Children    []PartyTypeSchemaNode `json:"children" bson:"children"`
}

type PartyTypeList struct {
	Items []*PartyType
}
