package exceptions

import (
	"fmt"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/util/validation/field"
	"net/http"
	"reflect"
)

type StatusError struct {
	ErrStatus metav1.Status
}

type APIStatus interface {
	Status() metav1.Status
}

var _ error = &StatusError{}

func (e *StatusError) Error() string {
	return e.ErrStatus.Message
}

func (e *StatusError) Status() metav1.Status {
	return e.ErrStatus
}

func NewInvalid(gk schema.GroupKind, uid string, errs field.ErrorList) *StatusError {
	causes := make([]metav1.StatusCause, 0, len(errs))
	for i := range errs {
		err := errs[i]
		causes = append(causes, metav1.StatusCause{
			Type:    metav1.CauseType(err.Type),
			Message: err.ErrorBody(),
			Field:   err.Field,
		})
	}
	return &StatusError{
		ErrStatus: metav1.Status{
			Status: metav1.StatusFailure,
			Code:   http.StatusUnprocessableEntity,
			Reason: metav1.StatusReasonInvalid,
			Details: &metav1.StatusDetails{
				UID:    uid,
				Group:  gk.Group,
				Kind:   gk.Kind,
				Causes: causes,
			},
			Message: "could not process entity",
		},
	}
}

func NewConflict(qualifiedResource schema.GroupResource, name string, err error) *StatusError {
	return &StatusError{metav1.Status{
		Status: metav1.StatusFailure,
		Code:   http.StatusConflict,
		Reason: metav1.StatusReasonConflict,
		Details: &metav1.StatusDetails{
			Group: qualifiedResource.Group,
			Kind:  qualifiedResource.Resource,
			//Name:  name,
		},
		Message: fmt.Sprintf("Operation cannot be fulfilled on %s %q: %v", qualifiedResource.String(), name, err),
	}}
}

// NewNotFound returns a new error which indicates that the resource of the kind and the name was not found.
func NewNotFound(qualifiedResource schema.GroupResource, name string) *StatusError {
	return &StatusError{metav1.Status{
		Status: metav1.StatusFailure,
		Code:   http.StatusNotFound,
		Reason: metav1.StatusReasonNotFound,
		Details: &metav1.StatusDetails{
			Group: qualifiedResource.Group,
			Kind:  qualifiedResource.Resource,
			// Name:  name,
		},
		Message: fmt.Sprintf("%s %q not found", qualifiedResource.String(), name),
	}}
}

// UnexpectedObjectError can be returned by FromObject if it's passed a non-status object.
type UnexpectedObjectError struct {
	Object runtime.Object
}

// Error returns an error message describing 'u'.
func (u *UnexpectedObjectError) Error() string {
	return fmt.Sprintf("unexpected object: %v", u.Object)
}

func FromObject(obj runtime.Object) error {
	switch t := obj.(type) {
	case *metav1.Status:
		return &StatusError{ErrStatus: *t}
	case runtime.Unstructured:
		var status metav1.Status
		obj := t.UnstructuredContent()
		if !reflect.DeepEqual(obj["kind"], "Status") {
			break
		}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(t.UnstructuredContent(), &status); err != nil {
			return err
		}
		if status.APIVersion != "v1" && status.APIVersion != "meta.k8s.io/v1" {
			break
		}
		return &StatusError{ErrStatus: status}
	}
	return &UnexpectedObjectError{obj}
}
