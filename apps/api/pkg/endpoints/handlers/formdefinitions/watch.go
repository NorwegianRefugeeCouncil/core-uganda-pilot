package formdefinitions

import (
	"github.com/gorilla/websocket"
	"github.com/nrc-no/core/apps/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/watch"
	"net/http"
)

// Watch formDefinitions
func (h *Handler) Watch(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	u := websocket.Upgrader{}
	conn, err := u.Upgrade(w, req, nil)
	if err != nil {
		h.scope.Error(err, w, req)
		return
	}
	defer conn.Close()

	err = h.storage.Watch(
		ctx,
		&v1.FormDefinition{},
		func(eventType string, obj runtime.Object) {
			err := conn.WriteJSON(&watch.Event{
				Type:   eventType,
				Object: obj,
			})
			if err != nil {
				h.scope.Error(err, w, req)
				return
			}
		})

	if err != nil {
		h.scope.Error(err, w, req)
		return
	}

}
