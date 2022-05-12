package types

type EntityDefinition struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Constraints EntityConstraints `json:"constraints"`
}
