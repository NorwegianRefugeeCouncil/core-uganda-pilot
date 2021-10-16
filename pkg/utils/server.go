package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	responseBytes, err := json.Marshal(data)
	if err != nil {
		logrus.WithError(err).Errorf("failed to marshal response")
		ErrorResponse(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(responseBytes)
	if err != nil {
		return
	}
}

func ErrorResponse(w http.ResponseWriter, err error) {
	logrus.WithError(err).Errorf("server error")

	status := ErrorToAPIStatus(err)
	code := int(status.Code)

	if status.Details != nil && status.Details.RetryAfterSeconds > 0 {
		delay := strconv.Itoa(int(status.Details.RetryAfterSeconds))
		w.Header().Set("Retry-After", delay)
	}

	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return
	}

	JSONResponse(w, code, status)
}

func BindJSON(req *http.Request, into interface{}) error {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return meta.NewBadRequest(fmt.Errorf("failed to read request body: %w", err).Error())
	}
	if err := json.Unmarshal(bodyBytes, &into); err != nil {
		return meta.NewBadRequest(fmt.Errorf("failed to unmarshal json body: %w", err).Error())
	}
	return nil
}

func GetPathParam(param string, w http.ResponseWriter, req *http.Request, into *string) bool {
	id, ok := mux.Vars(req)[param]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("path parameter '%s' not found in path", param)
		ErrorResponse(w, err)
		return false
	}
	*into = id
	return true
}

// ErrorToAPIStatus converts an error to an metav1.Status object.
func ErrorToAPIStatus(err error) *meta.Status {
	switch t := err.(type) {
	case meta.APIStatus:
		status := t.Status()
		if len(status.Status) == 0 {
			status.Status = meta.StatusFailure
		}
		switch status.Status {
		case meta.StatusSuccess:
			if status.Code == 0 {
				status.Code = http.StatusOK
			}
		case meta.StatusFailure:
			if status.Code == 0 {
				status.Code = http.StatusInternalServerError
			}
		default:
			//TODO log error
			if status.Code == 0 {
				status.Code = http.StatusInternalServerError
			}
		}
		return &status
	default:
		status := http.StatusInternalServerError
		// Log errors that were not converted to an error status
		// by REST storage - these typically indicate programmer
		// error by not using pkg/api/errors, or unexpected failure
		// cases.
		// TODO log error
		return &meta.Status{
			Status:  meta.StatusFailure,
			Code:    int32(status),
			Reason:  meta.StatusReasonUnknown,
			Message: err.Error(),
		}
	}
}
