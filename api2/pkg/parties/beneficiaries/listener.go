package beneficiaries

import (
	"context"
	"encoding/json"
	"github.com/nrc-no/core-kafka/pkg/intake"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// Listener receives submissions and updates beneficiary info
type Listener struct {
	reader *kafka.Reader
	store  *Store
	syncFn func(submission *intake.Submission) error
}

func NewListener(store *Store, reader *kafka.Reader) *Listener {
	l := &Listener{
		reader: reader,
		store:  store,
	}
	l.syncFn = l.sync
	return l
}

func (l *Listener) sync(submission *intake.Submission) error {

	// ctx := context.Background()
	// l.log(ctx).Infof("received submission: %#v", submission)
	//
	// bySubject := map[string][]intake.AnswerToQuestion{}
	//
	// for _, answer := range submission.Answers {
	// 	subjectType := answer.SubjectType
	// 	if subjectType != api.BeneficiaryType {
	// 		continue
	// 	}
	// 	subject := answer.Subject
	// 	bySubject[subject] = append(bySubject[subject], answer)
	// }

	// for beneficiaryID, questions := range bySubject {
	// 	var attributes []*api.AttributeValue
	// 	for _, question := range questions {
	// 		attributes = append(attributes, &api.AttributeValue{
	// 			Attribute: api.Attribute{
	// 				Name:        question.Name,
	// 				ValueType:   question.ValueType,
	// 				PartyTypes: question.PartyTypes,
	// 				// TODO
	// 				// LongFormulation:              question.LongFormulation,
	// 				// ShortFormulation:             question.ShortFormulation,
	// 				IsPersonallyIdentifiableInfo: question.IsPersonallyIdentifiableData,
	// 			},
	// 			Value: question.Answer,
	// 		})
	// 	}
	// 	if err := l.store.Upsert(ctx, beneficiaryID, attributes); err != nil {
	// 		l.logError(ctx, err, "failed to upsert beneficiary")
	// 	}
	// }

	return nil
}

func (l *Listener) Run(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		default:
			message, err := l.reader.ReadMessage(ctx)
			if err != nil {
				l.logError(ctx, err, "failed to read incoming message")
				continue
			}

			var submission intake.Submission
			if err := json.Unmarshal(message.Value, &submission); err != nil {
				l.logError(ctx, err, "failed to unmarshal submission")
				continue
			}

			if err := l.syncFn(&submission); err != nil {
				l.logError(ctx, err, "failed to process submission")
				continue
			}
		}

	}

}

func (l *Listener) log(ctx context.Context) *logrus.Entry {
	return logrus.WithContext(ctx).WithField("Listener", "Listener")
}

func (l *Listener) logError(ctx context.Context, err error, message string, args ...interface{}) {
	l.log(ctx).WithError(err).Errorf(message, args...)
}
