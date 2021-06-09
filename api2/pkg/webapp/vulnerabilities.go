package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/services/api"
	"github.com/nrc-no/core-kafka/pkg/services/vulnerability"
	"net/http"
	"net/url"
	"strings"
)

func (h *Handler) Vulnerabilities(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	cli := vulnerability.NewClient("http://localhost:9000")
	list, err := cli.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.template.ExecuteTemplate(w, "vulnerabilities", map[string]interface{}{
		"Vulnerabilities": list,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Vulnerability(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok {
		err := fmt.Errorf("no id in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cli := vulnerability.NewClient("http://localhost:9000")
	v, err := cli.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		if err := req.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var err error
		v, err = h.UpdateVulnerability(ctx, cli, v, req.Form)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", req.Host+"/vulnerabilities/"+v.ID)
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	err = h.template.ExecuteTemplate(w, "vulnerability", map[string]interface{}{
		"Vulnerability": v,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateVulnerability(
	ctx context.Context,
	cli *vulnerability.Client,
	vulnerability *api.Vulnerability,
	values url.Values) (*api.Vulnerability, error) {

	var attributes []string
	for key, values := range values {
		if strings.HasPrefix(key, "attribute-") {
			attributes = append(attributes, values[0])
		}
	}
	vulnerability.AttributesForDetermination = attributes

	out, err := cli.Update(ctx, vulnerability)
	if err != nil {
		return nil, err
	}
	return out, nil

}
