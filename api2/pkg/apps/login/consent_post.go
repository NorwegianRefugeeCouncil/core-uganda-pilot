package login

import (
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"net/http"
)

func (s *Server) PostConsent(w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		s.Error(w, err)
		return
	}

	consentChallenge := req.Form.Get("consent_challenge")

	consentGetResp, err := s.HydraAdmin.GetConsentRequest(admin.NewGetConsentRequestParams().
		WithContext(req.Context()).
		WithConsentChallenge(consentChallenge))
	if err != nil {
		s.Error(w, err)
		return
	}

	consentAcceptResp, err := s.HydraAdmin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
		WithContext(req.Context()).
		WithConsentChallenge(consentChallenge).
		WithBody(&models.AcceptConsentRequest{
			GrantScope:               req.Form["grant_scope"],
			GrantAccessTokenAudience: consentGetResp.GetPayload().RequestedAccessTokenAudience,
		}))
	if err != nil {
		s.Error(w, err)
		return
	}

	http.Redirect(w, req, *consentAcceptResp.GetPayload().RedirectTo, http.StatusFound)

}
