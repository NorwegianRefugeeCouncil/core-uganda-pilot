package types

// FormReference is an interface for accessing a form ID and DatabaseID
type FormReference interface {
	// GetFormID returns the form ID
	GetFormID() string
	// GetDatabaseID returns the form Database ID
	GetDatabaseID() string
}

