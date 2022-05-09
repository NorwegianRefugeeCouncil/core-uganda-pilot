package types

type AttributeConstraintType string

type AttributeConstraints struct {
	Required bool     `json:"required"`
	Unique   bool     `json:"unique"`
	Min      *float32 `json:"min"`
	Max      *float32 `json:"max"`
	Pattern  *string  `json:"pattern"`
	Enum     []string `json:"enum"`
	Custom   []string `json:"custom"`
}
