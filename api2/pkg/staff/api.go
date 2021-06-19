package staff

// Staff is a relationship between an organization and an individual
// that represents that the individual is working for an organization
type Staff struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organizationId"`
	IndividualID   string `json:"individualId"`
}

type StaffList struct {
	Items []*Staff `json:"items"`
}
