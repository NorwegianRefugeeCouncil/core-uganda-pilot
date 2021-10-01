package iam

import (
	"fmt"
	"net/http"
)

func (s *Server) getIndividual(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\n\nGET INDIVIDUAL\n\n\n")
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	b, err := s.individualStore.get(ctx, id)
	fmt.Printf("\n\n\nGET INDIVIDUAL\n\n\n", b)
	if err != nil {
		fmt.Printf("\n\n\nGET INDIVIDUAL\n ERROR \n\n", err)
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, b)
}
