package webapp

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (s *Server) Logout(w http.ResponseWriter, req *http.Request) {

	logoutURL, err := url.Parse(strings.Replace(s.publicOauth2Config.Endpoint.TokenURL, "/oauth2/token", "/oauth2/sessions/logout", -1))
	if err != nil {
		s.Error(w, err)
		return
	}

	session, err := s.sessionManager.Get(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	idTokenIntf, ok := session.Values["id-token"]
	if !ok {
		s.Error(w, fmt.Errorf("no id token in session"))
		return
	}

	idTokenStr, ok := idTokenIntf.(string)
	if !ok {
		s.Error(w, fmt.Errorf("wrong id token type"))
		return
	}

	params := url.Values{}
	params.Set("id_token_hint", idTokenStr)
	params.Set("post_logout_redirect_uri", s.baseURL)
	logoutURL.RawQuery = params.Encode()

	session.Values = map[interface{}]interface{}{}
	if err := session.Save(req, w); err != nil {
		s.Error(w, err)
		return
	}

	http.Redirect(w, req, logoutURL.String(), http.StatusTemporaryRedirect)
}
