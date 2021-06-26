package login

import (
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"net/http"
)

// GetLoginForm is called when Ory Hydra redirects the user to login
func (s *Server) GetLoginForm(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	qry := req.URL.Query()
	challenge := qry.Get("login_challenge")

	// Getting the login request
	resp, err := s.HydraAdmin.GetLoginRequest(
		admin.NewGetLoginRequestParams().
			WithLoginChallenge(challenge).
			WithContext(ctx),
	)

	if err != nil {
		s.Render(w, req, "login", map[string]interface{}{
			"Challenge": challenge,
			"Error":     err.Error(),
		})
		return
	}

	// When skip is true, user does not have to login again
	if resp.Payload.Skip != nil && *resp.Payload.Skip {

		// Accept login request
		respLoginAccept, err := s.HydraAdmin.AcceptLoginRequest(
			admin.NewAcceptLoginRequestParams().
				WithContext(ctx).
				WithLoginChallenge(challenge).
				WithBody(&models.AcceptLoginRequest{
					Subject: resp.GetPayload().Subject,
				}),
		)

		if err != nil {
			s.Render(w, req, "login", map[string]interface{}{
				"Challenge": challenge,
				"Error":     err.Error(),
			})
			return
		}

		// Redirect to the requested resource
		http.Redirect(w, req, *respLoginAccept.Payload.RedirectTo, http.StatusSeeOther)
		return
	}

	// Render the login page
	s.Render(w, req, "login", map[string]interface{}{
		"Challenge": challenge,
	})

}
