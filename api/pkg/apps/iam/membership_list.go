package iam

import "net/http"

func (s *Server) listMemberships(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var listOptions MembershipListOptions
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ret, err := s.membershipStore.list(ctx, listOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.json(w, http.StatusOK, ret)
}
