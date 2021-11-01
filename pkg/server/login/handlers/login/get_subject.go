package login

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"net/http"
)

func restfulGetSubject(hydraAdmin admin.ClientService) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		getSubject(hydraAdmin)(res.ResponseWriter, req.Request)
	}
}

func getSubject(hydraAdmin admin.ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		loginChallenge := req.URL.Query().Get("login_challenge")
		ctx := req.Context()

		loginRequest, err := hydraAdmin.GetLoginRequest(&admin.GetLoginRequestParams{
			Context:        ctx,
			LoginChallenge: loginChallenge,
		})

		if err != nil {
			renderSubjectLogin(w, loginChallenge, err)
			return
		}

		if loginRequest.Payload.Skip != nil && *loginRequest.Payload.Skip {
			resp, err := hydraAdmin.AcceptLoginRequest(&admin.AcceptLoginRequestParams{
				Body: &models.AcceptLoginRequest{
					Subject: loginRequest.Payload.Subject,
				},
				LoginChallenge: loginChallenge,
				Context:        ctx,
			})
			if err != nil {
				renderSubjectLogin(w, loginChallenge, err)
				return
			}
			w.Header().Set("Location", *resp.Payload.RedirectTo)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		renderSubjectLogin(w, loginChallenge, nil)

	}
}
