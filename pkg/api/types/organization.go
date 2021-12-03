package types

// Organization represents a "tenant" in Core.
type Organization struct {
	// ID of the Organization
	ID string `json:"id"`
	// Name of the Organization
	Name string `json:"name"`
}

// OrganizationList represents a list of Organization
type OrganizationList struct {
	Items []*Organization `json:"items"`
}
