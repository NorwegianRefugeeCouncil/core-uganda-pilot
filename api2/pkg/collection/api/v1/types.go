package v1

import v1 "github.com/nrc-no/core-kafka/pkg/i81n/api/v1"

// Question represents a question to ask a beneficiary
type Question struct {
	Formulation string
	Key         string
	AnswerType  string
}

// AnsweredQuestion represents an answer to a question
type AnsweredQuestion struct {
	Question
	Answer interface{}
}

// GlobalBeneficiaryQuestions represents the questions we want to ask to all
// beneficiaries regardless of their country
type GlobalBeneficiaryQuestions struct {
	Questions []Question
}

// GlobalBeneficiaryAnswers represents the answers for the GlobalBeneficiaryQuestions
type GlobalBeneficiaryAnswers struct {
	AnsweredQuestions []AnsweredQuestion
}

// CountrySpecificBeneficiaryQuestions represent the questions we need to ask
// beneficiaries depending on the country
type CountrySpecificBeneficiaryQuestions struct {
	Country   string
	Questions []Question
}

type ContextualQuestions struct {
	Country   string
	Questions []Question
}

func NewGlobalQuestions(questions []Question) ContextualQuestions {
	return ContextualQuestions{
		Country:   "",
		Questions: questions,
	}
}

func NewCountrySpecificQuestions(country string, questions []Question) ContextualQuestions {
	return ContextualQuestions{
		Country:   country,
		Questions: questions,
	}
}

// CountrySpecificBeneficiaryAnswers represents the answers to CountrySpecificBeneficiaryQuestions
type CountrySpecificBeneficiaryAnswers struct {
	Country           string
	AnsweredQuestions []AnsweredQuestion
}

// AggregatedBeneficiaryQuestions represents the aggregation of GlobalBeneficiaryQuestions
// and CountrySpecificBeneficiaryQuestions
type AggregatedBeneficiaryQuestions struct {
	GlobalQuestions          GlobalBeneficiaryQuestions
	CountrySpecificQuestions CountrySpecificBeneficiaryQuestions
}

// AggregatedBeneficiaryAnswers represents the aggregation of GlobalBeneficiaryAnswers
// and CountrySpecificBeneficiaryAnswers
type AggregatedBeneficiaryAnswers struct {
	GlobalAnswers          GlobalBeneficiaryAnswers
	CountrySpecificAnswers CountrySpecificBeneficiaryAnswers
}

type Topic string

type TopicDescription struct {
	Topic       Topic            `json:"topic" bson:"topic"`
	DataShape   DataShape        `json:"dataShape" bson:"dataShape"`
	Formulation TopicFormulation `json:"formulation" bson:"formulation"`
	Relevance   []TopicRelevance `json:"relevance" bson:"relevance"`
}

type DataShape struct {
	Type string `json:"type" bson:"type"`
}

type TopicRelevance struct {
	OrgContexts []string `json:"orgContexts" bson:"orgContexts"`
	GeoContexts []string `json:"geoContexts" bson:"geoContexts"`
	SvcContexts []string `json:"svcContexts" bson:"svcContexts"`
}

type Subject struct {
	Type string `json:"type" bson:"type"`
	ID   string `json:"id" bson:"id"`
}

type TopicFormulation struct {
	Long  v1.Translations `json:"long" bson:"long"`
	Short v1.Translations `json:"short" bson:"short"`
}

type Observation struct {
	ObservedValues map[string]interface{} `json:"observedValues" bson:"observedValues"`
	Topic          TopicDescription       `json:"topic" bson:"topic"`
	Subject        Subject                `json:"subject" bson:"subject"`
}
