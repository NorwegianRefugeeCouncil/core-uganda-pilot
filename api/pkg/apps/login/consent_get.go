package login

import (
	"fmt"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"net/http"
)

func (s *Server) GetConsent(w http.ResponseWriter, req *http.Request) {

	consentChallenge := req.URL.Query().Get("consent_challenge")
	if len(consentChallenge) == 0 {
		s.Error(w, fmt.Errorf("no consent challenge found in URL"))
		return
	}

	consentGetResp, err := s.HydraAdmin.GetConsentRequest(admin.NewGetConsentRequestParams().
		WithConsentChallenge(consentChallenge).
		WithHTTPClient(s.HydraHTTPClient).
		WithContext(req.Context()))
	if err != nil {
		s.Error(w, fmt.Errorf("could not request consent"))
		return
	}

	if consentGetResp.Payload.Skip {
		consentAcceptResp, err := s.HydraAdmin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
			WithContext(req.Context()).
			WithConsentChallenge(consentChallenge).
			WithHTTPClient(s.HydraHTTPClient).
			WithBody(&models.AcceptConsentRequest{
				GrantAccessTokenAudience: consentGetResp.GetPayload().RequestedAccessTokenAudience,
				GrantScope:               consentGetResp.GetPayload().RequestedScope,
			}))
		if err != nil {
			s.Error(w, fmt.Errorf("could not accept consent"))
			return
		}
		http.Redirect(w, req, *consentAcceptResp.GetPayload().RedirectTo, http.StatusFound)
		return
	}

	consentMessage := fmt.Sprintf("Application %s wants to access resources on your behalf and to:",
		consentGetResp.GetPayload().Client.ClientName)

	s.Render(w, req, "consent", map[string]interface{}{
		"ConsentChallenge": consentChallenge,
		"ConsentMessage":   consentMessage,
		"RequestedScopes":  consentGetResp.GetPayload().RequestedScope,
	})

}
