package formdefinitions

import (
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers"
	"github.com/nrc-no/core/apps/api/pkg/registry/rest"
	"github.com/nrc-no/core/apps/api/pkg/storage"
)

type Handler struct {
	storage storage.Interface
	getter  rest.Getter
	scope   *handlers.RequestScope
}
