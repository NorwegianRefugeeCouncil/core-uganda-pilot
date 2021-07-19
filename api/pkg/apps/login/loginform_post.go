package login

import (
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"net/http"
)

func (s *Server) PostLoginForm(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	// Parsing the form values
	if err := req.ParseForm(); err != nil {
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

	// Verify password hash
	individual, isValid := s.VerifyPassword(ctx, email, password)
	if !isValid {
		s.Render(w, req, "login", map[string]interface{}{
			"Challenge": loginChallenge,
			"Error":     "Invalid credentials",
		})
		return
	}

	// Getting login request
	getLoginRequestParams := admin.
		NewGetLoginRequestParams().
		WithLoginChallenge(loginChallenge).
		WithHTTPClient(s.HydraHTTPClient)

	_, err := s.HydraAdmin.GetLoginRequest(
		getLoginRequestParams)
	if err != nil {
		s.Error(w, err)
		return
	}

	// Accept login request
	respLoginAccept, err := s.HydraAdmin.AcceptLoginRequest(
		admin.NewAcceptLoginRequestParams().
			WithContext(ctx).
			WithLoginChallenge(loginChallenge).
			WithHTTPClient(s.HydraHTTPClient).
			WithBody(&models.AcceptLoginRequest{
				Remember: rememberMe,
				Subject:  &individual.ID,
				// Subject:  &party.ID, TODO
			}))
	if err != nil {
		s.Error(w, err)
		return
	}

	http.Redirect(w, req, *respLoginAccept.Payload.RedirectTo, http.StatusFound)

}
