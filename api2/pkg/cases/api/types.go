package api

type CaseType struct {
	ID          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	PartyTypeID string `json:"partyTypeId" bson:"partyTypeId"`
}

type Case struct {
	ID          string `json:"id" bson:"id"`
	CaseTypeID  string `json:"caseTypeId" bson:"caseTypeId"`
	PartyID     string `json:"partyId" bson:"partyId"`
	Description string `json:"description" bson:"description"`
	Done        bool   `json:"done" bson:"done"`
}

func (c *CaseType) String() string {
	return c.Name
}

type CaseList struct {
	Items []*Case `json:"items" bson:"items"`
}

type CaseTypeList struct {
	Items []*CaseType `json:"items" bson:"items"`
}

func (l *CaseTypeList) FindByID(id string) *CaseType {
	for _, caseType := range l.Items {
		if caseType.ID == id {
			return caseType
		}
	}
	return nil
}
