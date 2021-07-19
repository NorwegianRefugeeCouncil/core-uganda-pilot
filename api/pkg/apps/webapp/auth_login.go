package webapp

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

func (s *Server) Login(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		s.Error(w, err)
		return
	}

	state := base64.StdEncoding.EncodeToString(b)
	s.sessionManager.Put(ctx, "state", state)
	conf := s.oauth2Config

	redirectUrl := conf.AuthCodeURL(state)
	http.Redirect(w, req, redirectUrl, http.StatusTemporaryRedirect)

}
