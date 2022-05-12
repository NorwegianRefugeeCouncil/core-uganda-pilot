package types

type Entity struct {
	EntityDefinition
	Attributes    []Attribute          `json:"attributes"`
	Contraints    EntityConstraints    `json:"constraints"`
	Relationships []EntityRelationship `json:"relationships"`
}
