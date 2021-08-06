package webapp

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/iam"
	"golang.org/x/sync/errgroup"
	"net/http"
	"sync"
)

func (s *Server) Teams(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	t, err := iamClient.Teams().List(ctx, iam.TeamListOptions{})
	if err != nil {
		s.Error(w, err)
		return
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "teams", map[string]interface{}{
		"Teams": t,
	}); err != nil {
		s.Error(w, err)
		return
	}

}

func (s *Server) Team(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id found in path")
		s.Error(w, err)
		return
	}

	t, err := iamClient.Teams().Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	m, err := iamClient.Memberships().List(ctx, iam.MembershipListOptions{
		TeamID: id,
	})
	if err != nil {
		s.Error(w, err)
		return
	}

	n, err := iamClient.Nationalities().List(ctx, iam.NationalityListOptions{
		TeamID: id,
	})
	if err != nil {
		s.Error(w, err)
		return
	}
	if len(n.Items) < 1 {
		err := fmt.Errorf("failed to find nationality relationship for team")
		s.Error(w, err)
		return
	}

	c, err := iamClient.Countries().Get(ctx, n.Items[0].CountryID)
	if err != nil {
		s.Error(w, err)
		return
	}

	var members []*iam.Party
	lock := sync.Mutex{}

	g, ctx := errgroup.WithContext(ctx)

	for _, item := range m.Items {
		i := item
		g.Go(func() error {
			individual, err := iamClient.Parties().Get(ctx, i.IndividualID)
			if err != nil {
				return err
			}
			lock.Lock()
			defer lock.Unlock()
			members = append(members, individual)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		s.Error(w, err)
		return
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "team", map[string]interface{}{
		"Team":               t,
		"Members":            members,
		"LastNameAttribute":  iam.LastNameAttribute,
		"FirstNameAttribute": iam.FirstNameAttribute,
		"Constants":          s.Constants,
		"Country":            c,
	}); err != nil {
		s.Error(w, err)
		return
	}

}

func (s *Server) AddIndividualToTeam(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	i := req.URL.Query().Get("individualId")
	t := req.URL.Query().Get("teamId")

	m, err := iamClient.Memberships().List(ctx, iam.MembershipListOptions{
		IndividualID: i,
		TeamID:       t,
	})
	if err != nil {
		s.Error(w, err)
		return
	}

	if len(m.Items) != 0 {
		return
	}

	_, err = iamClient.Memberships().Create(ctx, &iam.Membership{
		ID:           "",
		TeamID:       t,
		IndividualID: i,
	})
	if err != nil {
		s.Error(w, err)
		return
	}

	w.Header().Set("Location", "/teams/"+t)
}
