package types

// Database is a namespace for FormDefinition. It acts as a scope for Access Control.
// For example, an Organization might create a Database for a specific country, region, or
// thematic area.
type Database struct {
	// ID is the id of the Database
	ID string `json:"id"`
	// Name is the name of the database
	Name string `json:"name"`
}

// DatabaseList represent a list of Database
type DatabaseList struct {
	Items []*Database `json:"items"`
}
