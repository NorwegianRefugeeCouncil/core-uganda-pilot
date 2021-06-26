package webapp

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"net/http"
)

func (s *Server) Login(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state := base64.StdEncoding.EncodeToString(b)
	s.sessionManager.Put(ctx, "state", state)

	provider, err := oidc.NewProvider(ctx, "http://localhost:4444/")
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

	redirectUrl := conf.AuthCodeURL(state)
	http.Redirect(w, req, redirectUrl, http.StatusTemporaryRedirect)

}
