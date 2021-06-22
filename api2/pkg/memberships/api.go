package memberships

type Membership struct {
	ID           string `json:"id"`
	TeamID       string `json:"teamId"`
	IndividualID string `json:"individualId"`
}

type MembershipList struct {
	Items []*Membership `json:"items"`
}
