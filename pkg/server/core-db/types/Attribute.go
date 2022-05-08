package types

type AttributeType string

const (
	String     AttributeType = "string"
	Number     AttributeType = "number"
	Boolean    AttributeType = "boolean"
	Date       AttributeType = "date"
	Time       AttributeType = "time"
	DateTime   AttributeType = "datetime"
	Month      AttributeType = "month"
	Week       AttributeType = "week"
	Coordinate AttributeType = "coordinate"
	File       AttributeType = "file"
)

type Attribute struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	List        bool                 `json:"list"`
	Type        AttributeType        `json:"type"`
	EntityID    string               `json:"entityId"`
	Constraints AttributeConstraints `json:"constraints"`
}
