package watch

import (
	"github.com/nrc-no/core/apps/api/pkg/runtime"
)

// EventType defines the possible types of events.
type EventType string

const (
	Added    EventType = "ADDED"
	Modified EventType = "MODIFIED"
	Deleted  EventType = "DELETED"
	Bookmark EventType = "BOOKMARK"
	Error    EventType = "ERROR"
)

// +k8s:deepcopy-gen=true
type Event struct {
	Type   EventType      `json:"type"`
	Object runtime.Object `json:"object"`
}

type Interface interface {
	Stop()
	ResultChan() <-chan Event
}
