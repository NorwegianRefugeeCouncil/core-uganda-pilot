package teams

type Team struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type TeamList struct {
	Items []*Team `json:"items"`
}
