package login

import (
	"github.com/ory/hydra-client-go/client/admin"
	"net/http"
)

// GetLogoutForm is called when Ory Hydra redirects the user to login
func (s *Server) GetLogoutForm(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	qry := req.URL.Query()
	challenge := qry.Get("logout_challenge")

	// Getting the login request
	_, err := s.HydraAdmin.GetLogoutRequest(
		admin.NewGetLogoutRequestParams().
			WithLogoutChallenge(challenge).
			WithContext(ctx),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	acceptResp, err := s.HydraAdmin.AcceptLogoutRequest(
		admin.NewAcceptLogoutRequestParams().
			WithLogoutChallenge(challenge).
			WithContext(ctx))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to the requested resource
	http.Redirect(w, req, *acceptResp.Payload.RedirectTo, http.StatusSeeOther)
	return

}
