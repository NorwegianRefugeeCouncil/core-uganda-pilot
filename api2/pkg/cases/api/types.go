package api

type CaseType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PartyTypeID string `json:"partyTypeId"`
}

type Case struct {
	ID          string `json:"id"`
	CaseTypeID  string `json:"caseTypeId"`
	PartyID     string `json:"partyId"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func (c *CaseType) String() string {
	return c.Name
}

type CaseList struct {
	Items []*Case `json:"items"`
}

type CaseTypeList struct {
	Items []*CaseType `json:"items"`
}
