package intake

import (
	"github.com/nrc-no/core-kafka/pkg/expressions"
	i81n "github.com/nrc-no/core-kafka/pkg/i81n/api/v1"
	"github.com/nrc-no/core-kafka/pkg/parties/api"
)

type Question struct {
	Name                         string
	ShortFormulation             i81n.Translations
	LongFormulation              i81n.Translations
	ValueType                    expressions.ValueType
	SubjectType                  api.SubjectType
	IsPersonallyIdentifiableData bool
}

type QuestionTranslation struct {
}

type AnswerToQuestion struct {
	Question
	Subject string
	Answer  interface{}
}

type Submission struct {
	Answers []AnswerToQuestion `json:"answers"`
	ID      string             `json:"id"`
}
