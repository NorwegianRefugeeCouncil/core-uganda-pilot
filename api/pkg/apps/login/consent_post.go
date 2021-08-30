package login

import (
	"github.com/nrc-no/core/pkg/apps/iam"
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
		WithHTTPClient(s.HydraHTTPClient).
		WithConsentChallenge(consentChallenge))
	if err != nil {
		s.Error(w, err)
		return
	}

	subject := consentGetResp.GetPayload().Subject
	individual, err := s.iam.Individuals().Get(req.Context(), subject)
	if err != nil {
		s.Error(w, err)
		return
	}

	fullName := individual.Get(iam.FullNameAttribute.ID)

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
		s.Error(w, err)
		return
	}

	http.Redirect(w, req, *consentAcceptResp.GetPayload().RedirectTo, http.StatusFound)

}
