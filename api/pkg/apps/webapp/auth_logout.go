package webapp

import (
	"net/http"
	"net/url"
	"strings"
)

func (s *Server) Logout(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	logoutURL, err := url.Parse(strings.Replace(s.oauth2Config.Endpoint.TokenURL, "/oauth2/token", "/oauth2/sessions/logout", -1))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := url.Values{}
	params.Set("id_token_hint", s.sessionManager.GetString(ctx, "id-token"))
	params.Set("post_logout_redirect_uri", s.baseURL)
	logoutURL.RawQuery = params.Encode()

	s.sessionManager.Clear(ctx)
	http.Redirect(w, req, logoutURL.String(), http.StatusTemporaryRedirect)
}
