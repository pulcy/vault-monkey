package client

import "fmt"

// Values of Status.Status
const (
	StatusSuccess = "Success"
	StatusFailure = "Failure"
)

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
)

type (

	// StatusReason is an enumeration of possible failure causes.  Each StatusReason
	// must map to a single HTTP status code, but multiple reasons may map
	// to the same HTTP status code.
	StatusReason string

	// Status is a return value for calls that don't return other objects.
	Status struct {
		TypeMeta `json:",inline"`
		// Standard list metadata.
		// More info: http://releases.k8s.io/release-1.3/docs/devel/api-conventions.md#types-kinds
		ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

		// Status of the operation.
		// One of: "Success" or "Failure".
		// More info: http://releases.k8s.io/release-1.3/docs/devel/api-conventions.md#spec-and-status
		Status string `json:"status,omitempty" protobuf:"bytes,2,opt,name=status"`
		// A human-readable description of the status of this operation.
		Message string `json:"message,omitempty" protobuf:"bytes,3,opt,name=message"`
		// A machine-readable description of why this operation is in the
		// "Failure" status. If this value is empty there
		// is no information available. A Reason clarifies an HTTP status
		// code but does not override it.
		Reason StatusReason `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason,casttype=StatusReason"`
		// Extended data associated with the reason.  Each reason may define its
		// own extended details. This field is optional and the data returned
		// is not guaranteed to conform to any schema except that defined by
		// the reason type.
		Details *StatusDetails `json:"details,omitempty" protobuf:"bytes,5,opt,name=details"`
		// Suggested HTTP return code for this status, 0 if not set.
		Code int32 `json:"code,omitempty" protobuf:"varint,6,opt,name=code"`
	}

	// StatusDetails is a set of additional properties that MAY be set by the
	// server to provide additional information about a response. The Reason
	// field of a Status object defines what attributes will be set. Clients
	// must ignore fields that do not match the defined type of each attribute,
	// and should assume that any attribute may be empty, invalid, or under
	// defined.
	StatusDetails struct {
		// The name attribute of the resource associated with the status StatusReason
		// (when there is a single name which can be described).
		Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
		// The group attribute of the resource associated with the status StatusReason.
		Group string `json:"group,omitempty" protobuf:"bytes,2,opt,name=group"`
		// The kind attribute of the resource associated with the status StatusReason.
		// On some operations may differ from the requested resource Kind.
		// More info: http://releases.k8s.io/release-1.3/docs/devel/api-conventions.md#types-kinds
		Kind string `json:"kind,omitempty" protobuf:"bytes,3,opt,name=kind"`
		// The Causes array includes more details associated with the StatusReason
		// failure. Not all StatusReasons may provide detailed causes.
		Causes []StatusCause `json:"causes,omitempty" protobuf:"bytes,4,rep,name=causes"`
		// If specified, the time in seconds before the operation should be retried.
		RetryAfterSeconds int32 `json:"retryAfterSeconds,omitempty" protobuf:"varint,5,opt,name=retryAfterSeconds"`
	}

	// StatusCause provides more information about an api.Status failure, including
	// cases when multiple errors are encountered.
	StatusCause struct {
		// A machine-readable description of the cause of the error. If this value is
		// empty there is no information available.
		Type CauseType `json:"reason,omitempty" protobuf:"bytes,1,opt,name=reason,casttype=CauseType"`
		// A human-readable description of the cause of the error.  This field may be
		// presented as-is to a reader.
		Message string `json:"message,omitempty" protobuf:"bytes,2,opt,name=message"`
		// The field of the resource that has caused this error, as named by its JSON
		// serialization. May include dot and postfix notation for nested attributes.
		// Arrays are zero-indexed.  Fields may appear more than once in an array of
		// causes due to fields having multiple errors.
		// Optional.
		//
		// Examples:
		//   "name" - the field "name" on the current resource
		//   "items[0].name" - the field "name" on the first array entry in "items"
		Field string `json:"field,omitempty" protobuf:"bytes,3,opt,name=field"`
	}

	// CauseType is a machine readable value providing more detail about what
	// occurred in a status response. An operation may have multiple causes for a
	// status (whether Failure or Success).
	CauseType string
)

func (s *Status) Error() string {
	return fmt.Sprintf("%d: %s", s.Code, s.Message)
}
