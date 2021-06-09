package intake

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/nrc-no/core-kafka/pkg/testhelpers"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestPostSubmission(t *testing.T) {

	ctx := context.Background()
	mongoClient, err := testhelpers.NewMongoClient(ctx)
	if !assert.NoError(t, err) {
		return
	}

	store := NewStore(mongoClient)

	conn, err := kafka.Dial("tcp", "localhost:9092")
	if !assert.NoError(t, err) {
		return
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if !assert.NoError(t, err) {
		return
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if !assert.NoError(t, err) {
		return
	}
	defer controllerConn.Close()

	topic := kafka.TopicConfig{
		Topic:             "submissions",
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	if err := controllerConn.CreateTopics(topic); !assert.NoError(t, err) {
		return
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Balancer: &kafka.LeastBytes{},
	}

	handler := NewHandler(store, writer)

	postSubmission := PostSubmission{}
	bodyBytes, err := json.Marshal(postSubmission)
	if !assert.NoError(t, err) {
		return
	}

	req := httptest.NewRequest("GET", "/", bytes.NewReader(bodyBytes))
	recorder := httptest.NewRecorder()

	handler.PostSubmission(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

}
