package endpoints

import (
	"github.com/emicklei/go-restful"
	"github.com/nrc-no/core/apps/api/pkgs2/models"
	"net/http"
	"path"
)

func Register(group, version, pluralName string, goRestfulContainer *restful.Container) {

	ws := new(restful.WebService)

	basePath := path.Join("/apis", group, version, pluralName)

	ws.Path(basePath)
	ws.Doc("API at /apis/core.nrc.no/v1/formdefinitions")
	ws.ApiVersion("core.nrc.no/v1")

	ws.Route(ws.Method("GET").Path("").To(listResource()).
		Doc("list FormDefinitions").
		Operation("listFormDefinitions").
		Returns(http.StatusOK, "OK", &models.FormDefinitionList{}).
		Writes(&models.FormDefinitionList{}))

	ws.Route(ws.Method("GET").Path("{id}").To(getResource()).
		Doc("read the specified FormDefinition").
		Operation("getFormDefinition").
		Returns(http.StatusOK, "OK", &models.FormDefinition{}).
		Writes(&models.FormDefinition{}))

	goRestfulContainer.Add(ws)

}

func listResource() func(request *restful.Request, response *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {

		result := &models.FormDefinitionList{
			TypeMeta: models.TypeMeta{
				APIVersion: "nrc.core.no/v1",
				Kind:       "FormDefinitionList",
			},
			ObjectMeta: models.ObjectMeta{},
			Items: []models.FormDefinition{
				{
					TypeMeta: models.TypeMeta{
						APIVersion: "nrc.core.no/v1",
						Kind:       "FormDefinition",
					},
					ObjectMeta: models.ObjectMeta{
						Labels:      nil,
						Annotations: nil,
					},
					Spec: models.FormDefinitionSpec{
						Group: "nrc.core.no/v1",
						Names: models.FormDefinitionNames{
							Singular: "customform",
							Plural:   "customforms",
							Kind:     "CustomForm",
						},
						Versions: []models.FormDefinitionVersion{
							{
								Name: "v1",
								Schema: models.FormModel{
									FormSchema: models.FormSchema{
										Root: models.FormElement{
											ID:       "",
											Children: []models.FormElement{},
											Type:     "bla",
										},
									},
								},
							},
						},
					},
				},
			},
		}

		response.WriteAsJson(result)

	}
}

func getResource() func(request *restful.Request, response *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {

	}
}
