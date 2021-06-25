package partytypes

type PartyType struct {
	ID        string `json:"id" bson:"id"`
	Name      string `json:"name" bson:"name"`
	IsBuiltIn bool   `json:"isBuiltIn" bson:"isBuiltIn"`
}

type PartyTypeList struct {
	Items []*PartyType
}

func (p *PartyTypeList) FindByID(id string) *PartyType {
	for _, item := range p.Items {
		if item.ID == id {
			return item
		}
	}
	return nil
}
