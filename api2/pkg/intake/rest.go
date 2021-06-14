package intake

import (
	"context"
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type Handler struct {
	store  *Store
	writer *kafka.Writer
}

func NewHandler(
	store *Store,
	writer *kafka.Writer,
) *Handler {
	h := &Handler{
		store:  store,
		writer: writer,
	}
	return h
}

func (h *Handler) log(ctx context.Context) *logrus.Entry {
	return logrus.WithContext(ctx).WithField("handler", "IntakeHandler")
}

func (h *Handler) logError(ctx context.Context, err error, message string, args ...interface{}) {
	h.log(ctx).WithError(err).Errorf(message, args...)
}

type PostSubmission struct {
	Answers []AnswerToQuestion `json:"answers" bson:"answers"`
}

func (h *Handler) PostSubmission(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var p PostSubmission
	if err := json.Unmarshal(bodyBytes, &p); err != nil {
		h.logError(ctx, err, "failed to unmarshal request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var submission = Submission{
		ID:      uuid.NewV4().String(),
		Answers: p.Answers,
	}

	submissionJson, err := json.Marshal(submission)
	if err != nil {
		h.logError(ctx, err, "failed to marshal submission to json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.store.StoreSubmission(ctx, &submission); err != nil {
		h.logError(ctx, err, "failed to store submission")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Infof("publishing submission %s", submission.ID)

	if err := h.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(submission.ID),
		Topic: "submissions",
		Value: submissionJson,
	}); err != nil {
		h.logError(ctx, err, "failed to publish submission")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
