package webapp

import (
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"net/http"
)

func (s *Server) Callback(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	provider, err := oidc.NewProvider(ctx, "http://localhost:4444")
	if err != nil {
		http.Error(w, "code not found", http.StatusInternalServerError)
		return
	}

	conf := oauth2.Config{
		ClientID:     s.OauthClientID,
		ClientSecret: s.OauthClientSecret,
		RedirectURL:  "http://localhost:9000/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	sessionState := s.sessionManager.PopString(ctx, "state")
	queryState := req.URL.Query().Get("state")

	if sessionState != queryState {
		http.Error(w, "state mismatch", http.StatusInternalServerError)
		return
	}

	code := req.URL.Query().Get("code")
	if len(code) == 0 {
		http.Error(w, "code not found", http.StatusInternalServerError)
		return
	}

	token, err := conf.Exchange(ctx, code)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to exchange code: %v", err).Error(), http.StatusInternalServerError)
		return
	}

	s.sessionManager.Put(ctx, "access-token", token.AccessToken)
	s.sessionManager.Put(ctx, "refresh-token", token.RefreshToken)

	http.Redirect(w, req, "/", http.StatusSeeOther)

}
