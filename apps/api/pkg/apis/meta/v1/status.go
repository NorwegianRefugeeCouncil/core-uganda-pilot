package v1

type StatusType string

const (
	StatusSuccess StatusType = "Success"
	StatusFailure StatusType = "Failure"
)

type Status struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`

	// Status is the status of the operation. One of "Failure" or "Success"
	Status StatusType `json:"status,omitempty"`

	// Message is a human-readable description of this operation
	Message string `json:"message,omitempty"`

	// Reason is a machine-readable description of why this operation is in the
	// "Failure" status.
	Reason StatusReason `json:"reason,omitempty"`

	// Details represents extended data associated with the reason.
	Details *StatusDetails `json:"details"`

	// Suggested HTTP status code.
	Code int32 `json:"code,omitempty"`
}

type StatusDetails struct {
	UID    string        `json:"uid,omitempty"`
	Group  string        `json:"group,omitempty"`
	Kind   string        `json:"kind,omitempty"`
	Causes []StatusCause `json:"causes,omitempty"`
}

type StatusReason string

const (
	StatusReasonUnknown          StatusReason = "Unknown"
	StatusReasonForbidden        StatusReason = "Forbidden"
	StatusReasonNotFound         StatusReason = "NotFound"
	StatusReasonAlreadyExists    StatusReason = "AlreadyExists"
	StatusReasonConflict         StatusReason = "Conflict"
	StatusReasonBadRequest       StatusReason = "BadRequest"
	StatusReasonMethodNotAllowed StatusReason = "MethodNotAllowed"
	StatusReasonNotAcceptable    StatusReason = "NotAcceptable"
	StatusReasonInternalError    StatusReason = "InternalError"
	StatusReasonInvalid          StatusReason = "Invalid"
)

type StatusCause struct {
	Type    CauseType `json:"reason,omitempty"`
	Message string    `json:"message,omitempty"`
	Field   string    `json:"field,omitempty"`
}

type CauseType string

const (
	CauseTypeFieldValueNotFound     CauseType = "FieldValueNotFound"
	CauseTypeFieldValueRequired     CauseType = "FieldValueRequired"
	CauseTypeFieldValueDuplicate    CauseType = "FieldValueDuplicate"
	CauseTypeFieldValueInvalid      CauseType = "FieldValueInvalid"
	CauseTypeFieldValueNotSupported CauseType = "FieldValueNotSupported"
)
