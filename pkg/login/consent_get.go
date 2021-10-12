package login

import (
	"fmt"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) GetConsent(w http.ResponseWriter, req *http.Request) {
	logrus.Trace("getting consent challenge")

	consentChallenge := req.URL.Query().Get("consent_challenge")
	if len(consentChallenge) == 0 {
		s.Error(w, fmt.Errorf("no consent challenge found in URL"))
		return
	}

	logrus.Tracef("getting consent request")

	consentGetResp, err := s.HydraAdmin.GetConsentRequest(admin.NewGetConsentRequestParams().
		WithConsentChallenge(consentChallenge).
		WithHTTPClient(s.HydraHTTPClient).
		WithContext(req.Context()))
	if err != nil {
		s.Error(w, fmt.Errorf("could not request consent: %w", err))
		return
	}

	consentRequest := consentGetResp.GetPayload()

	if consentRequest.Skip {
		logrus.Trace("consent skipped. Accepting consent request")

		consentAcceptResp, err := s.HydraAdmin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
			WithContext(req.Context()).
			WithConsentChallenge(consentChallenge).
			WithHTTPClient(s.HydraHTTPClient).
			WithBody(&models.AcceptConsentRequest{
				GrantAccessTokenAudience: consentRequest.RequestedAccessTokenAudience,
				GrantScope:               consentRequest.RequestedScope,
			}))
		if err != nil {
			s.Error(w, fmt.Errorf("could not accept consent: %w", err))
			return
		}

		logrus.Trace("accepted consent request")

		http.Redirect(w, req, *consentAcceptResp.GetPayload().RedirectTo, http.StatusFound)
		return
	}

	logrus.Trace("consent required. Displaying consent request to user")

	consentMessage := fmt.Sprintf("Application %s wants to access resources on your behalf and to:",
		consentRequest.Client.ClientName)

	s.Render(w, "consent", map[string]interface{}{
		"ConsentChallenge": consentChallenge,
		"ConsentMessage":   consentMessage,
		"RequestedScopes":  consentRequest.RequestedScope,
	})
}
