package webapp

import (
	"github.com/nrc-no/core-kafka/pkg/individuals"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"net/http"
)

func (h *Handler) Login(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	qry := req.URL.Query()
	challenge := qry.Get("login_challenge")

	resp, err := h.hydraAdminClient.Admin.GetLoginRequest(
		admin.NewGetLoginRequestParams().
			WithLoginChallenge(challenge).
			WithContext(ctx),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.Payload.Skip != nil && *resp.Payload.Skip {
		respLoginAccept, err := h.hydraAdminClient.Admin.AcceptLoginRequest(
			admin.NewAcceptLoginRequestParams().
				WithContext(ctx).
				WithLoginChallenge(challenge).
				WithBody(&models.AcceptLoginRequest{
					Subject: resp.GetPayload().Subject,
				}),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, req, *respLoginAccept.Payload.RedirectTo, http.StatusSeeOther)
		return
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "auth_login", map[string]interface{}{
		"Challenge": challenge,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) PostLogin(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	values := req.Form

	party, err := h.partyStore.Find(ctx, parties.FindOptions{
		Attributes: map[string]string{
			individuals.EMailAttribute.ID: values.Get("email"),
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("wrong email or password"))
		return
	}

	isValid := h.credentialsClient.VerifyPassword(ctx, party.ID, values.Get("password"))

	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("wrong email or password"))
		return
	}

	loginChallenge := values.Get("login_challenge")

	loginRequest := admin.NewGetLoginRequestParams().WithLoginChallenge(loginChallenge)
	_, err = h.hydraAdminClient.Admin.GetLoginRequest(loginRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respLoginAccept, err := h.hydraAdminClient.Admin.AcceptLoginRequest(
		admin.NewAcceptLoginRequestParams().
			WithContext(ctx).
			WithLoginChallenge(loginChallenge).
			WithBody(&models.AcceptLoginRequest{
				Remember: values.Get("remember_me") == "true",
				Subject:  &party.ID,
			}))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, *respLoginAccept.Payload.RedirectTo, http.StatusFound)

}
