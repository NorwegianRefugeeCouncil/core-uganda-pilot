package organizations

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/server/admin/templates"
	"net/http"
)

func restfulAddGet() restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handleAddGet()(res.ResponseWriter, req.Request)
	}
}

func handleAddGet() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		_ = templates.Template.ExecuteTemplate(w, "organization_add", map[string]interface{}{
			"IsNew": true,
		})
		return
	}
}
