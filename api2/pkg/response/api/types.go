package api

import subjects "github.com/nrc-no/core-kafka/pkg/parties/api"

type ResponseResult struct {
	Request         ResponseRequest
	ReferredToTeams []string
}

type ResponseRequest struct {
	SubjectType string `json:"subjectType" bson:"subjectType"`
	Subject     string `json:"subject" bson:"subject"`
	Attributes  map[string]subjects.AttributeValue
}
