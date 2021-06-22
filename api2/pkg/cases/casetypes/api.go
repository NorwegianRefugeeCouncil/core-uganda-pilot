package casetypes

type CaseType struct {
	ID          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	PartyTypeID string `json:"partyTypeId" bson:"partyTypeId"`
	TeamID      string `json:"teamId" bson:"teamId"`
}

type CaseTypeList struct {
	Items []*CaseType `json:"items" bson:"items"`
}

func (c *CaseType) String() string {
	return c.Name
}

func (l *CaseTypeList) FindByID(id string) *CaseType {
	for _, caseType := range l.Items {
		if caseType.ID == id {
			return caseType
		}
	}
	return nil
}
