package writers

import (
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/sirupsen/logrus"
	"net/http"
)

type statusError interface {
	Status() metav1.Status
}

func ErrorNegotiated(err error, serializer runtime.Serializer, w http.ResponseWriter, req *http.Request) int {
	status := ErrorToApiStatus(err)
	code := int(status.Code)
	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return code
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := serializer.Encode(status, w); err != nil {
		logrus.Errorf("unable to encode error response: %v", err)
	}
	return code

}

func ErrorToApiStatus(err error) *metav1.Status {
	switch t := err.(type) {
	case statusError:
		status := t.Status()
		if len(status.Status) == 0 {
			status.Status = metav1.StatusFailure
		}
		switch status.Status {
		case metav1.StatusSuccess:
			if status.Code == 0 {
				status.Code = http.StatusOK
			}
		case metav1.StatusFailure:
			if status.Code == 0 {
				status.Code = http.StatusInternalServerError
			}
		default:
			if status.Code == 0 {
				status.Code = http.StatusInternalServerError
			}
		}
		status.Kind = "Status"
		status.APIVersion = "meta/v1"
		return &status
	default:
		status := http.StatusInternalServerError
		return &metav1.Status{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "meta/v1",
				Kind:       "Status",
			},
			Status:  metav1.StatusFailure,
			Message: err.Error(),
			Reason:  metav1.StatusReasonUnknown,
			Code:    int32(status),
		}
	}
}
