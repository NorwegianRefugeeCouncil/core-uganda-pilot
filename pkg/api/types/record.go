package types

type RecordRef struct {
	ID         string
	DatabaseID string
	FormID     string
}

type FormRef struct {
	DatabaseID string
	FormID     string
}

type Record struct {
	ID         string                 `json:"id"`
	Seq        int64                  `json:"seq"`
	DatabaseID string                 `json:"databaseId"`
	FormID     string                 `json:"formId"`
	ParentID   *string                `json:"parentId"`
	Values     map[string]interface{} `json:"values"`
}

type RecordList struct {
	Items []*Record `json:"items"`
}

type RecordListOptions struct {
	DatabaseID string
	FormID     string
}
