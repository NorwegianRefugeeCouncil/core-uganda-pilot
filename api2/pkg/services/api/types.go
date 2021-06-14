package api

import (
	"github.com/nrc-no/core-kafka/pkg/expressions"
	i81n "github.com/nrc-no/core-kafka/pkg/i81n/api/v1"
	"github.com/nrc-no/core-kafka/pkg/parties/api"
)

type Vulnerability struct {
	ID                         string                `json:"id" bson:"id"`
	Name                       string                `json:"name" bson:"name"`
	SubjectType                api.SubjectType       `json:"subjectType" bson:"subjectType"`
	ValueType                  expressions.ValueType `json:"valueType" bson:"valueType"`
	LongFormulation            i81n.Translations     `json:"longFormulation" bson:"longFormulation"`
	ShortFormulation           i81n.Translations     `json:"shortFormulation" bson:"shortFormulation"`
	AttributesForDetermination []string              `json:"attributesForDetermination" bson:"attributesForDetermination"`
}

type VulnerabilityList struct {
	Items []*Vulnerability
}

type Service struct {
	Name                string                 `json:"name" bson:"name"`
	EligibilityCriteria []EligibilityCriterion `json:"eligibilityCriteria" bson:"eligibilityCriteria"`
}

type EligibilityCriterion struct {
	Vulnerability *string `json:"vulnerability" bson:"vulnerability"`
}

type ServiceSubjectRepresentation struct {
	SubjectType string                 `json:"subjectType" bson:"subjectType"`
	Subject     string                 `json:"subject" bson:"subject"`
	ServiceName string                 `json:"serviceName" bson:"serviceName"`
	Attributes  map[string]interface{} `json:"attributes" bson:"attributes"`
}

type ServiceSubjectVulnerabilities struct {
	SubjectType     string                 `json:"subjectType" bson:"subjectType"`
	Subject         string                 `json:"subject" bson:"subject"`
	Vulnerabilities map[string]interface{} `json:"vulnerabilities" bson:"vulnerabilities"`
}
