package types

import "time"

// Database is a namespace for FormDefinition. It acts as a scope for Access Control.
// For example, an Organization might create a Database for a specific country, region, or
// thematic area.
type Database struct {
	// ID is the id of the Database
	ID string `json:"id"`
	// Name is the name of the database
	Name string `json:"name"`
	// CreatedAt is the creation time of the database
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt is the last update time of the database
	UpdatedAt time.Time `json:"updatedAt"`
}

// DatabaseList represent a list of Database
type DatabaseList struct {
	Items []*Database `json:"items"`
}
