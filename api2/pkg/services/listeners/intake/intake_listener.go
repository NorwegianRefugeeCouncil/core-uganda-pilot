package intake

import (
	"context"
	"encoding/json"
	"github.com/nrc-no/core-kafka/pkg/intake"
	"github.com/nrc-no/core-kafka/pkg/parties/api"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// IntakeListener receives submissions and forwards the attributes to the relevant services
type IntakeListener struct {
	reader *kafka.Reader
	syncFn func(submission *intake.Submission) error
}

func NewIntakeListener(reader *kafka.Reader) *IntakeListener {
	l := &IntakeListener{
		reader: reader,
	}
	l.syncFn = l.sync
	return l
}

func (l *IntakeListener) sync(submission *intake.Submission) error {
	ctx := context.Background()
	l.log(ctx).Infof("received submission: %#v", submission)

	var byTypeAndSubject = map[api.SubjectType]map[string][]intake.AnswerToQuestion{}

	for _, answer := range submission.Answers {
		subjectType := answer.SubjectType
		if _, ok := byTypeAndSubject[subjectType]; !ok {
			byTypeAndSubject[subjectType] = map[string][]intake.AnswerToQuestion{}
		}
		subject := answer.Subject
		if _, ok := byTypeAndSubject[subjectType][subject]; !ok {
			byTypeAndSubject[subjectType][subject] = []intake.AnswerToQuestion{}
		}
		byTypeAndSubject[subjectType][subject] = append(byTypeAndSubject[subjectType][subject], answer)
	}

	return nil
}

func (l *IntakeListener) Run(ctx context.Context) {

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

func (l *IntakeListener) log(ctx context.Context) *logrus.Entry {
	return logrus.WithContext(ctx).WithField("Listener", "IntakeListener")
}

func (l *IntakeListener) logError(ctx context.Context, err error, message string, args ...interface{}) {
	l.log(ctx).WithError(err).Errorf(message, args...)
}
