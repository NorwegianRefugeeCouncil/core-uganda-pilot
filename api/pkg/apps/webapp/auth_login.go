package webapp

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

func (s *Server) Login(w http.ResponseWriter, req *http.Request) {

	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		s.Error(w, err)
		return
	}

	state := base64.StdEncoding.EncodeToString(b)

	session, err := s.sessionManager.Get(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	session.Values["state"] = state
	if err := session.Save(req, w); err != nil {
		s.Error(w, err)
		return
	}

	redirectUrl := s.publicOauth2Config.AuthCodeURL(state)
	http.Redirect(w, req, redirectUrl, http.StatusTemporaryRedirect)

}
