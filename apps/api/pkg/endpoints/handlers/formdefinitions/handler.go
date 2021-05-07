package formdefinitions

import (
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers"
	"github.com/nrc-no/core/apps/api/pkg/storage"
)

type Handler struct {
	storage storage.Interface
	scope   *handlers.RequestScope
}
