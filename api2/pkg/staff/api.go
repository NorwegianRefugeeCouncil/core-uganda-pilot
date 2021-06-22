package staff

// Staff is a relationship between an organization and an individual
// that represents that the individual is working for that organization
type Staff struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organizationId"`
	IndividualID   string `json:"individualId"`
}

// StaffList is a list of Staff
type StaffList struct {
	Items []*Staff `json:"items"`
}
