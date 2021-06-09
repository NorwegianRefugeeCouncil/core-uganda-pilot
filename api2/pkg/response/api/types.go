package api

import subjects "github.com/nrc-no/core-kafka/pkg/subjects/api"

type ResponseResult struct {
	Request         ResponseRequest
	ReferredToTeams []string
}

type ResponseRequest struct {
	SubjectType string `json:"subjectType"`
	Subject     string `json:"subject"`
	Attributes  map[string]subjects.AttributeValue
}
