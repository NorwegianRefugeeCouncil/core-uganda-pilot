package login

import (
	"fmt"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

// GetLoginForm is called when Ory Hydra redirects the user to login
func (s *Server) GetLoginForm(w http.ResponseWriter, req *http.Request) {

	logrus.Trace("getting login challenge")

	ctx := req.Context()
	qry := req.URL.Query()
	challenge := qry.Get("login_challenge")

	logrus.Trace("getting login request")

	// Getting the login request
	resp, err := s.HydraAdmin.GetLoginRequest(
		admin.NewGetLoginRequestParams().
			WithLoginChallenge(challenge).
			WithHTTPClient(s.HydraHTTPClient).
			WithContext(ctx),
	)
	if err != nil {
		logrus.WithError(err).Error("failed to get login request")
		s.Render(w, req, "login", map[string]interface{}{
			"Challenge": challenge,
			"Error":     fmt.Errorf("failed to get login request: %v", err),
		})
		return
	}

	logrus.Trace("found login request")

	// When skip is true, user does not have to login again
	if resp.Payload.Skip != nil && *resp.Payload.Skip {

		logrus.Trace("skipping manual login request")

		// Accept login request
		respLoginAccept, err := s.HydraAdmin.AcceptLoginRequest(
			admin.NewAcceptLoginRequestParams().
				WithContext(ctx).
				WithLoginChallenge(challenge).
				WithHTTPClient(s.HydraHTTPClient).
				WithBody(&models.AcceptLoginRequest{
					Subject: resp.GetPayload().Subject,
				}),
		)
		if err != nil {
			logrus.WithError(err).Error("failed to accept login request")
			s.Render(w, req, "login", map[string]interface{}{
				"Challenge": challenge,
				"Error":     fmt.Errorf("failed to accept login request: %v", err),
			})
			return
		}

		redirectTo := *respLoginAccept.Payload.RedirectTo

		// Redirect to the requested resource
		logrus.Tracef("accepted login request. Redirecting to %s", redirectTo)
		http.Redirect(w, req, redirectTo, http.StatusSeeOther)
		return
	}

	logrus.Tracef("not skipping manual login request. Rendering login form")

	// Render the login page
	s.Render(w, req, "login", map[string]interface{}{
		"Challenge": challenge,
	})

}
