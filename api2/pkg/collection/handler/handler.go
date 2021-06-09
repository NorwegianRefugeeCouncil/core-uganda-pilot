package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	v1 "github.com/nrc-no/core-kafka/pkg/collection/api/v1"
	"github.com/nrc-no/core-kafka/pkg/collection/store"
	"io/ioutil"
	"net/http"
)

type CollectionHandler struct {
	router     *mux.Router
	topicStore *store.Topic
}

func NewCollectionHandler(topicStore *store.Topic) *CollectionHandler {
	router := mux.NewRouter()
	return &CollectionHandler{
		router:     router,
		topicStore: topicStore,
	}
}

type PostTopic struct {
}

func (c *CollectionHandler) PostTopic(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	topic := v1.TopicDescription{}
	if err := json.Unmarshal(bodyBytes, &topic); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.topicStore.CreateTopic(ctx, &topic); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(topic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)

}

func (c *CollectionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c.router.ServeHTTP(w, req)
}
