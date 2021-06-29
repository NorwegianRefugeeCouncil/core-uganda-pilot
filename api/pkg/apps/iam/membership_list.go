package iam

import "net/http"

func (s *Server) ListMemberships(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var listOptions MembershipListOptions
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ret, err := s.MembershipStore.List(ctx, listOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.JSON(w, http.StatusOK, ret)
}
