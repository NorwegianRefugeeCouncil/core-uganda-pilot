package v1

type TypeMeta struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
}

type ObjectMeta struct {
	Name string `json:"name,omitempty"`
}

type ListMeta struct {
}

type StatusReason string

// Values of Status.Status
const (
	StatusSuccess = "Success"
	StatusFailure = "Failure"
)

const (
	// StatusReasonUnknown means the server has declined to indicate a specific reason.
	// The details field may contain other information about this error.
	// Status code 500.
	StatusReasonUnknown StatusReason = ""

	// StatusReasonUnauthorized means the server can be reached and understood the request, but requires
	// the user to present appropriate authorization credentials (identified by the WWW-Authenticate header)
	// in order for the action to be completed. If the user has specified credentials on the request, the
	// server considers them insufficient.
	// Status code 401
	StatusReasonUnauthorized StatusReason = "Unauthorized"

	// StatusReasonForbidden means the server can be reached and understood the request, but refuses
	// to take any further action.  It is the result of the server being configured to deny access for some reason
	// to the requested resource by the client.
	// Details (optional):
	//   "kind" string - the kind attribute of the forbidden resource
	//                   on some operations may differ from the requested
	//                   resource.
	//   "id"   string - the identifier of the forbidden resource
	// Status code 403
	StatusReasonForbidden StatusReason = "Forbidden"

	// StatusReasonNotFound means one or more resources required for this operation
	// could not be found.
	// Details (optional):
	//   "kind" string - the kind attribute of the missing resource
	//                   on some operations may differ from the requested
	//                   resource.
	//   "id"   string - the identifier of the missing resource
	// Status code 404
	StatusReasonNotFound StatusReason = "NotFound"

	// StatusReasonAlreadyExists means the resource you are creating already exists.
	// Details (optional):
	//   "kind" string - the kind attribute of the conflicting resource
	//   "id"   string - the identifier of the conflicting resource
	// Status code 409
	StatusReasonAlreadyExists StatusReason = "AlreadyExists"

	// StatusReasonConflict means the requested operation cannot be completed
	// due to a conflict in the operation. The client may need to alter the
	// request. Each resource may define custom details that indicate the
	// nature of the conflict.
	// Status code 409
	StatusReasonConflict StatusReason = "Conflict"

	// StatusReasonGone means the item is no longer available at the server and no
	// forwarding address is known.
	// Status code 410
	StatusReasonGone StatusReason = "Gone"

	// StatusReasonInvalid means the requested create or update operation cannot be
	// completed due to invalid data provided as part of the request. The client may
	// need to alter the request. When set, the client may use the StatusDetails
	// message field as a summary of the issues encountered.
	// Details (optional):
	//   "kind" string - the kind attribute of the invalid resource
	//   "id"   string - the identifier of the invalid resource
	//   "causes"      - one or more StatusCause entries indicating the data in the
	//                   provided resource that was invalid.  The code, message, and
	//                   field attributes will be set.
	// Status code 422
	StatusReasonInvalid StatusReason = "Invalid"

	// StatusReasonServerTimeout means the server can be reached and understood the request,
	// but cannot complete the action in a reasonable time. The client should retry the request.
	// This is may be due to temporary server load or a transient communication issue with
	// another server. Status code 500 is used because the HTTP spec provides no suitable
	// server-requested client retry and the 5xx class represents actionable errors.
	// Details (optional):
	//   "kind" string - the kind attribute of the resource being acted on.
	//   "id"   string - the operation that is being attempted.
	//   "retryAfterSeconds" int32 - the number of seconds before the operation should be retried
	// Status code 500
	StatusReasonServerTimeout StatusReason = "ServerTimeout"

	// StatusReasonTimeout means that the request could not be completed within the given time.
	// Clients can get this response only when they specified a timeout param in the request,
	// or if the server cannot complete the operation within a reasonable amount of time.
	// The request might succeed with an increased value of timeout param. The client *should*
	// wait at least the number of seconds specified by the retryAfterSeconds field.
	// Details (optional):
	//   "retryAfterSeconds" int32 - the number of seconds before the operation should be retried
	// Status code 504
	StatusReasonTimeout StatusReason = "Timeout"

	// StatusReasonTooManyRequests means the server experienced too many requests within a
	// given window and that the client must wait to perform the action again. A client may
	// always retry the request that led to this error, although the client should wait at least
	// the number of seconds specified by the retryAfterSeconds field.
	// Details (optional):
	//   "retryAfterSeconds" int32 - the number of seconds before the operation should be retried
	// Status code 429
	StatusReasonTooManyRequests StatusReason = "TooManyRequests"

	// StatusReasonBadRequest means that the request itself was invalid, because the request
	// doesn't make any sense, for example deleting a read-only object.  This is different than
	// StatusReasonInvalid above which indicates that the API call could possibly succeed, but the
	// data was invalid.  API calls that return BadRequest can never succeed.
	// Status code 400
	StatusReasonBadRequest StatusReason = "BadRequest"

	// StatusReasonMethodNotAllowed means that the action the client attempted to perform on the
	// resource was not supported by the code - for instance, attempting to delete a resource that
	// can only be created. API calls that return MethodNotAllowed can never succeed.
	// Status code 405
	StatusReasonMethodNotAllowed StatusReason = "MethodNotAllowed"

	// StatusReasonNotAcceptable means that the accept types indicated by the client were not acceptable
	// to the server - for instance, attempting to receive protobuf for a resource that supports only json and yaml.
	// API calls that return NotAcceptable can never succeed.
	// Status code 406
	StatusReasonNotAcceptable StatusReason = "NotAcceptable"

	// StatusReasonRequestEntityTooLarge means that the request entity is too large.
	// Status code 413
	StatusReasonRequestEntityTooLarge StatusReason = "RequestEntityTooLarge"

	// StatusReasonUnsupportedMediaType means that the content type sent by the client is not acceptable
	// to the server - for instance, attempting to send protobuf for a resource that supports only json and yaml.
	// API calls that return UnsupportedMediaType can never succeed.
	// Status code 415
	StatusReasonUnsupportedMediaType StatusReason = "UnsupportedMediaType"

	// StatusReasonInternalError indicates that an internal error occurred, it is unexpected
	// and the outcome of the call is unknown.
	// Details (optional):
	//   "causes" - The original error
	// Status code 500
	StatusReasonInternalError StatusReason = "InternalError"

	// StatusReasonExpired indicates that the request is invalid because the content you are requesting
	// has expired and is no longer available. It is typically associated with watches that can't be
	// serviced.
	// Status code 410 (gone)
	StatusReasonExpired StatusReason = "Expired"

	// StatusReasonServiceUnavailable means that the request itself was valid,
	// but the requested service is unavailable at this time.
	// Retrying the request after some time might succeed.
	// Status code 503
	StatusReasonServiceUnavailable StatusReason = "ServiceUnavailable"
)

// CauseType is a machine readable value providing more detail about what
// occurred in a status response. An operation may have multiple causes for a
// status (whether Failure or Success).
type CauseType string

const (
	// CauseTypeFieldValueNotFound is used to report failure to find a requested value
	// (e.g. looking up an ID).
	CauseTypeFieldValueNotFound CauseType = "FieldValueNotFound"
	// CauseTypeFieldValueRequired is used to report required values that are not
	// provided (e.g. empty strings, null values, or empty arrays).
	CauseTypeFieldValueRequired CauseType = "FieldValueRequired"
	// CauseTypeFieldValueDuplicate is used to report collisions of values that must be
	// unique (e.g. unique IDs).
	CauseTypeFieldValueDuplicate CauseType = "FieldValueDuplicate"
	// CauseTypeFieldValueInvalid is used to report malformed values (e.g. failed regex
	// match).
	CauseTypeFieldValueInvalid CauseType = "FieldValueInvalid"
	// CauseTypeFieldValueNotSupported is used to report valid (as per formatting rules)
	// values that can not be handled (e.g. an enumerated string).
	CauseTypeFieldValueNotSupported CauseType = "FieldValueNotSupported"
	// CauseTypeUnexpectedServerResponse is used to report when the server responded to the client
	// without the expected return type. The presence of this cause indicates the error may be
	// due to an intervening proxy or the server software malfunctioning.
	CauseTypeUnexpectedServerResponse CauseType = "UnexpectedServerResponse"
	// FieldManagerConflict is used to report when another client claims to manage this field,
	// It should only be returned for a request using server-side apply.
	CauseTypeFieldManagerConflict CauseType = "FieldManagerConflict"
	// CauseTypeResourceVersionTooLarge is used to report that the requested resource version
	// is newer than the data observed by the API server, so the request cannot be served.
	CauseTypeResourceVersionTooLarge CauseType = "ResourceVersionTooLarge"
)

type Status struct {
	Status  string       `json:"status,omitempty"`
	Message string       `json:"message,omitempty"`
	Reason  StatusReason `json:"reason,omitempty"`
	Code    int32        `json:"code,omitempty"`
}

type StatusDetails struct {
	Causes []StatusCause `json:"causes,omitempty"`
}

type StatusCause struct {
	Type    CauseType `json:"type,omitempty"`
	Message string    `json:"message,omitempty"`
	Field   string    `json:"field,omitempty"`
}

type CreateOptions struct {
}

type UpdateOptions struct {
}

type DeleteOptions struct {
}

type ListOptions struct {
}

type GetOptions struct {
}
