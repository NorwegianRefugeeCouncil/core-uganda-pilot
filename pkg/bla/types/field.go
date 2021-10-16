package types

type FieldDefinition struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Required  bool      `json:"required"`
	FieldType FieldType `json:"fieldType"`
}
