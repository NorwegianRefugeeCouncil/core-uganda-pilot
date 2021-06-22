package cases

type Case struct {
	ID          string `json:"id" bson:"id"`
	CaseTypeID  string `json:"caseTypeId" bson:"caseTypeId"`
	PartyID     string `json:"partyId" bson:"partyId"`
	Description string `json:"description" bson:"description"`
	Done        bool   `json:"done" bson:"done"`
	ParentID    string `json:"parentId" bson:"parentId"`
}

type CaseList struct {
	Items []*Case `json:"items" bson:"items"`
}
