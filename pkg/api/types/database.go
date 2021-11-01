package types

type Database struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DatabaseList struct {
	Items []*Database `json:"items"`
}

func NewDatabase(name string) Database {
	return Database{Name: name}
}
