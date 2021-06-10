package response

import (
	"context"
	"encoding/json"
	"github.com/nrc-no/core-kafka/pkg/parties/api"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// BeneficiaryListener Listens to subject situation change and enqueues
// a need for response
type BeneficiaryListener struct {
	reader *kafka.Reader
	syncFn func(beneficiary *api.Beneficiary) error
}

func (l *BeneficiaryListener) Sync(beneficiary *api.Beneficiary) error {

	return nil
}

func (l *BeneficiaryListener) Run(ctx context.Context) {
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

			var beneficiary api.Beneficiary
			if err := json.Unmarshal(message.Value, &beneficiary); err != nil {
				l.logError(ctx, err, "failed to unmarshal beneficiary")
				continue
			}

			if err := l.syncFn(&beneficiary); err != nil {
				l.logError(ctx, err, "failed to process beneficiary")
				continue
			}
		}
	}
}

func (l *BeneficiaryListener) log(ctx context.Context) *logrus.Entry {
	return logrus.WithContext(ctx).WithField("Listener", "Listener")
}

func (l *BeneficiaryListener) logError(ctx context.Context, err error, message string, args ...interface{}) {
	l.log(ctx).WithError(err).Errorf(message, args...)
}
