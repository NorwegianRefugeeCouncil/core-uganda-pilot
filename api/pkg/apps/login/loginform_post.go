package login

import (
	"fmt"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) PostLoginForm(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	logrus.Trace("parsing login form values")

	// Parsing the form values
	if err := req.ParseForm(); err != nil {
		logrus.WithError(err).Error()
		s.Render(w, req, "login", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	values := req.Form
	loginChallenge := values.Get("login_challenge")
	rememberMe := values.Get("remember_me") == "true"
	email := values.Get("email")
	password := values.Get("password")

	defer func() {
		// Clearing values
		values.Set("email", "")
		values.Set("password", "")
		email = ""
		password = ""
	}()

	logrus.Trace("verifying username/password")

	// Verify password hash
	individual, isValid := s.VerifyPassword(ctx, email, password)
	if !isValid {
		logrus.Trace("invalid username/password")
		s.Render(w, req, "login", map[string]interface{}{
			"Challenge": loginChallenge,
			"Error":     "Invalid credentials",
		})
		return
	}

	logrus.Trace("getting login request")

	// Getting login request
	getLoginRequestParams := admin.
		NewGetLoginRequestParams().
		WithLoginChallenge(loginChallenge).
		WithHTTPClient(s.HydraHTTPClient)

	_, err := s.HydraAdmin.GetLoginRequest(
		getLoginRequestParams)
	if err != nil {
		s.Error(w, fmt.Errorf("failed to get login request: %v", err))
		return
	}

	logrus.Trace("accepting login request")

	// Accept login request
	respLoginAccept, err := s.HydraAdmin.AcceptLoginRequest(
		admin.NewAcceptLoginRequestParams().
			WithContext(ctx).
			WithLoginChallenge(loginChallenge).
			WithHTTPClient(s.HydraHTTPClient).
			WithBody(&models.AcceptLoginRequest{
				Remember: rememberMe,
				Subject:  &individual.ID,
			}))
	if err != nil {
		s.Error(w, fmt.Errorf("failed to accept login request: %v", err))
		return
	}

	logrus.Tracef("redirecting user to %s", *respLoginAccept.Payload.RedirectTo)

	http.Redirect(w, req, *respLoginAccept.Payload.RedirectTo, http.StatusFound)

}
