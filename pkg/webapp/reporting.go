package webapp

import (
	"github.com/nrc-no/core/pkg/cms"
	iam2 "github.com/nrc-no/core/pkg/iam"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func (s *Server) Reporting(w http.ResponseWriter, req *http.Request) {
	var iList *iam2.IndividualList
	var cList *cms.CaseList
	var ctList *cms.CaseTypeList

	g, _ := errgroup.WithContext(req.Context())

	// initialise clients
	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	cmsClient, err := s.CMSClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	// get individuals
	g.Go(func() error {
		iList, err = iamClient.Individuals().List(req.Context(), iam2.IndividualListOptions{})
		if err != nil {
			return err
		}
		return nil
	})

	// get list of cases
	g.Go(func() error {
		cList, err = cmsClient.Cases().List(req.Context(), cms.CaseListOptions{})
		if err != nil {
			return err
		}
		return nil
	})

	// get list of case types
	g.Go(func() error {
		ctList, err = cmsClient.CaseTypes().List(req.Context(), cms.CaseTypeListOptions{})
		if err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		s.Error(w, err)
		return
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "reporting", map[string]interface{}{
		"Individuals":                   iList,
		"Cases":                         cList,
		"CaseTypes":                     ctList,
		"GenderAttributeID":             iam2.GenderAttribute.ID,
		"DisplacementStatusAttributeID": iam2.DisplacementStatusAttribute.ID,
		"BirthDateAttributeID":          iam2.BirthDateAttribute.ID,
		"BeneficiaryPartyTypeID":        iam2.BeneficiaryPartyType.ID,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
