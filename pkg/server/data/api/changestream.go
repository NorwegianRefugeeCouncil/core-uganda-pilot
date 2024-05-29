package api

import (
	"encoding/json"
	"fmt"
)

type Changes struct {
	Items []ChangeItem `json:"items"`
}

func (c Changes) String() string {
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return string(jsonBytes)
}

type ChangeItem struct {
	Sequence       int64    `json:"sequence"`
	TableName      string   `json:"table_name"`
	RecordID       string   `json:"record_id"`
	RecordRevision Revision `json:"record_revision"`
}

func (c ChangeItem) String() string {
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return string(jsonBytes)
}
