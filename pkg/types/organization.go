package types

type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type OrganizationList struct {
	Items []*Organization `json:"items"`
}
