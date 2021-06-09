package beneficiaries

import (
	"context"
	"encoding/json"
	"github.com/nrc-no/core-kafka/pkg/expressions"
	i81n "github.com/nrc-no/core-kafka/pkg/i81n/api/v1"
	"github.com/nrc-no/core-kafka/pkg/intake"
	"github.com/nrc-no/core-kafka/pkg/subjects/api"
	"github.com/nrc-no/core-kafka/pkg/testhelpers"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListener(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mongoClient, err := testhelpers.NewMongoClient(ctx)
	if !assert.NoError(t, err) {
		return
	}
	store := NewStore(mongoClient)
	reader, err := testhelpers.NewReader(ctx, "intake", nil)
	if !assert.NoError(t, err) {
		return
	}
	writer, err := testhelpers.NewWriter("")
	if !assert.NoError(t, err) {
		return
	}

	listener := NewListener(store, reader)
	listener.syncFn = func(submission *intake.Submission) error {
		defer cancel()
		return listener.sync(submission)
	}
	go func() {
		listener.Run(ctx)
	}()

	submissionID := uuid.NewV4().String()
	beneficiaryID := uuid.NewV4().String()
	submission := intake.Submission{
		Answers: []intake.AnswerToQuestion{
			{
				Answer:  "john",
				Subject: beneficiaryID,
				Question: intake.Question{
					Name: "firstName",
					ShortFormulation: []i81n.Translation{
						{
							Locale:      "en",
							Translation: "First Name",
						},
					},
					LongFormulation: []i81n.Translation{
						{
							Locale:      "en",
							Translation: "What is your First Name",
						},
					},
					ValueType:                    expressions.ValueType{},
					SubjectType:                  api.BeneficiaryType,
					IsPersonallyIdentifiableData: true,
				},
			},
		},
		ID: submissionID,
	}
	intakeJson, err := json.Marshal(submission)
	if !assert.NoError(t, err) {
		return
	}

	err = writer.WriteMessages(ctx, kafka.Message{
		Topic: "intake",
		Key:   []byte(submissionID),
		Value: intakeJson,
	})
	if !assert.NoError(t, err) {
		return
	}

	<-ctx.Done()

}
