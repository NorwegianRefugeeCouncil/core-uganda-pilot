package webapp

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"path"
	"regexp"
	"strings"
)

func (s *Server) serveJS(staticDir string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		file, ok := mux.Vars(req)["file"]
		if !ok || len(file) == 0 {
			err := fmt.Errorf("no file found in path")
			s.Error(w, err)
			return
		}
		if file, ok = sanitize(file); !ok {
			s.Error(w, errors.New("invalid filename"))
		}
		p := path.Join(staticDir, file)
		http.ServeFile(w, req, p)
	}
}

func sanitize(file string) (string, bool) {
	s, err := safeFileName(file)
	return s, err == nil
}

func safeFileName(str string) (string, error) {
	name := strings.ToLower(str)
	name = path.Clean(path.Base(name))
	name = strings.Trim(name, " ")
	separators, err := regexp.Compile(`[ &_=+:]`)
	if err != nil {
		return "", err
	}
	name = separators.ReplaceAllString(name, "-")
	legal, err := regexp.Compile(`[^[:alnum:]-.]`)
	if err != nil {
		return "", err
	}
	name = legal.ReplaceAllString(name, "")
	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}
	return name, nil
}
