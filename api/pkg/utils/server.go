package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/validation"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
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
	if status, ok := err.(*validation.Status); ok {
		code := status.Code
		if code == 0 {
			code = http.StatusInternalServerError
		}
		JSONResponse(w, code, status)
		return
	} else {
		status := &validation.Status{
			Code: http.StatusInternalServerError,
		}
		JSONResponse(w, status.Code, status)
		return
	}
}

func BindJSON(req *http.Request, into interface{}) error {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %v", err)
	}

	if err := json.Unmarshal(bodyBytes, &into); err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
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
