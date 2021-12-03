package types

// RecordRef represents a key that allows referencing a single Record.
type RecordRef struct {
	// ID is the Record.ID
	ID string `json:"id"`
	// DatabaseID is the Record.DatabaseID
	DatabaseID string `json:"databaseId"`
	// FormID is the Record.FormID
	FormID string `json:"formId"`
}

// FormRef represents a key that allows referencing a single FormDefinition.
type FormRef struct {
	// DatabaseID represents the FormDefinition.DatabaseID
	DatabaseID string `json:"databaseId"`
	// FormID represents the FormDefinition.ID
	FormID string `json:"formId"`
}

// Record represents an entry in a Form.
type Record struct {
	// ID of the record
	ID string `json:"id"`
	// Seq of the Record. This value is automatically increased by the database.
	// The presence of this field allows us to sort the table by insertion order.
	Seq int64 `json:"seq"`
	// DatabaseID of the Record
	DatabaseID string `json:"databaseId"`
	// FormID of the Record. Represents in which Form this record belongs.
	FormID string `json:"formId"`
	// ParentID represents the parent of the Record. In cases where
	// a Record is part of a SubForm, this field records the "Parent" form ID.
	ParentID *string `json:"parentId"`
	// Values is an arbitrary map of values that correspond to the FormDefinition.Fields.
	// The key of the map is the FieldDefinition.ID ! (not the FormDefinition.Name, not
	// the FormDefinition.Code)
	Values map[string]interface{} `json:"values"`
}

// RecordList represents a list of Record
type RecordList struct {
	Items []*Record `json:"items"`
}

// RecordListOptions represents the options for listing Record.
type RecordListOptions struct {
	DatabaseID string
	FormID     string
}
