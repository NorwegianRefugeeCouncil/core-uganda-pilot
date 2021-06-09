package intake

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/intake"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntakeListener(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "submissions",
	})
	if err := reader.SetOffset(0); !assert.NoError(t, err) {
		return
	}

	listener := NewIntakeListener(reader)
	listener.syncFn = func(submission *intake.Submission) error {
		t.Logf("received submission: %#v", submission)
		cancel()
		return nil
	}

	listener.Run(ctx)

}
