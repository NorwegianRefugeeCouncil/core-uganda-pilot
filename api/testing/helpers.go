package testing

import (
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/util/exceptions"
	"github.com/stretchr/testify/assert"
	"testing"
)

/**
  assert.Error(t, err)
  if !assert.IsType(t, &exceptions.StatusError{}, err) {
    return
  }
  status := err.(*exceptions.StatusError).Status()
  assert.Equal(t, metav1.StatusFailure, status.Status)
  assert.Equal(t, http.StatusUnprocessableEntity, int(status.Code))
  assert.Equal(t, metav1.StatusReasonInvalid, status.Reason)
  assert.Equal(t, "could not process entity", status.Message)
  assert.Equal(t, "FormDefinition", status.Details.Kind)
  assert.Equal(t, "core", status.Details.Group)
  if !assert.Len(t, status.Details.Causes, 1) {
    return
  }
  assert.Equal(t, "Required value: group is required", status.Details.Causes[0].Message)
  assert.Equal(t, metav1.CauseTypeFieldValueRequired, status.Details.Causes[0].Type)
  assert.Equal(t, "spec.group", status.Details.Causes[0].Field)
*/

func AssertStatusFailure(t *testing.T, status *metav1.Status) bool {
	return assert.Equal(t, metav1.StatusFailure, status.Status)
}

func AssertSuccess(t *testing.T, status *metav1.Status) bool {
	return assert.Equal(t, metav1.StatusSuccess, status.Status)
}

func AssertStatusCode(t *testing.T, status *metav1.Status, statusCode int) bool {
	return assert.Equal(t, statusCode, int(status.Code))
}

func AssertFailureReason(t *testing.T, status *metav1.Status, reason metav1.StatusReason) bool {
	return assert.Equal(t, reason, status.Reason)
}

func AssertStatusMessage(t *testing.T, status *metav1.Status, message string) bool {
	if !assert.NotNil(t, status) {
		return false
	}
	return assert.Equal(t, message, status.Message)
}

func AssertStatusDetailsKind(t *testing.T, status *metav1.Status, kind string) bool {
	if !assert.NotNil(t, status) {
		return false
	}
	if !assert.NotNil(t, status.Details) {
		return false
	}
	return assert.Equal(t, kind, status.Details.Kind)
}

func AssertStatusDetailsGroup(t *testing.T, status *metav1.Status, group string) bool {
	if !assert.NotNil(t, status) {
		return false
	}
	if !assert.NotNil(t, status.Details) {
		return false
	}
	return assert.Equal(t, group, status.Details.Group)
}

func AssertStatusCauseCount(t *testing.T, status *metav1.Status, count int) bool {
	if !assert.NotNil(t, status) {
		return false
	}
	if !assert.NotNil(t, status.Details) {
		return false
	}
	return assert.Len(t, status.Details.Causes, count)
}

func AssertIsErrStatus(t *testing.T, err error, out *metav1.Status) bool {
	if !assert.NotNil(t, err) {
		return false
	}
	if !assert.IsType(t, &exceptions.StatusError{}, err) {
		return false
	}
	if out != nil {
		coerced := err.(*exceptions.StatusError)
		*out = coerced.Status()
	}
	return true
}

func AssertCauseMessage(t *testing.T, cause metav1.StatusCause, message string) bool {
	return assert.Equal(t, message, cause.Message)
}

func AssertCauseType(t *testing.T, cause metav1.StatusCause, causeType metav1.CauseType) bool {
	return assert.Equal(t, causeType, cause.Type)
}

func AssertCauseField(t *testing.T, cause metav1.StatusCause, field string) bool {
	return assert.Equal(t, field, cause.Field)
}
