package login

import (
	"fmt"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) PostConsent(w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		s.Error(w, fmt.Errorf("failed to parse form: %v", err))
		return
	}

	logrus.Trace("getting consent challenge")

	consentChallenge := req.Form.Get("consent_challenge")

	logrus.Tracef("getting consent request")

	consentGetResp, err := s.HydraAdmin.GetConsentRequest(admin.NewGetConsentRequestParams().
		WithContext(req.Context()).
		WithHTTPClient(s.HydraHTTPClient).
		WithConsentChallenge(consentChallenge))
	if err != nil {
		s.Error(w, fmt.Errorf("failed to get consent request: %v", err))
		return
	}

	logrus.Trace("getting consent subject")

	subject := consentGetResp.GetPayload().Subject

	logrus.Tracef("consent subject: %s", subject)

	individual, err := s.iam.Individuals().Get(req.Context(), subject)
	if err != nil {
		s.Error(w, fmt.Errorf("failed to get individual: %v", err))
		return
	}

	logrus.Tracef("found consent subject")

	fullName := individual.Get(iam.FullNameAttribute.ID)

	logrus.Tracef("accepting consent request")

	consentAcceptResp, err := s.HydraAdmin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
		WithContext(req.Context()).
		WithConsentChallenge(consentChallenge).
		WithHTTPClient(s.HydraHTTPClient).
		WithBody(&models.AcceptConsentRequest{
			GrantScope:               req.Form["grant_scope"],
			GrantAccessTokenAudience: consentGetResp.GetPayload().RequestedAccessTokenAudience,
			Remember:                 true,
			Session: &models.ConsentRequestSession{
				AccessToken: map[string]interface{}{
					"full_name": fullName,
				},
				IDToken: map[string]interface{}{
					"full_name": fullName,
				},
			},
		}))
	if err != nil {
		s.Error(w, fmt.Errorf("failed to accept consent request: %v", err))
		return
	}

	logrus.Tracef("consent request accepted. Redirecting to %s", *consentAcceptResp.GetPayload().RedirectTo)

	http.Redirect(w, req, *consentAcceptResp.GetPayload().RedirectTo, http.StatusFound)

}
