package teams

type Team struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type TeamList struct {
	Items []*Team `json:"items"`
}

func (l *TeamList) FindByID(id string) *Team {
	for _, team := range l.Items {
		if team.ID == id {
			return team
		}
	}
	return nil
}
