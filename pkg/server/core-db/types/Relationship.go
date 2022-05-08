package types

type Cardinality string

const (
	OneToOne   Cardinality = "OneToOne"
	OneToMany  Cardinality = "OneToMany"
	ManyToOne  Cardinality = "ManyToOne"
	ManyToMany Cardinality = "ManyToMany"
)

type EntityRelationship struct {
	ID             string      `json:"id"`
	Cardinality    Cardinality `json:"cardinality"`
	SourceEntityID string      `json:"sourceEntityId"`
	TargetEntityID string      `json:"targetEntityId"`
}
