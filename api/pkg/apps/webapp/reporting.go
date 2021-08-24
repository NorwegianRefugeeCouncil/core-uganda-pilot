package webapp

import (
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"net/http"
)

func (s *Server) Reporting(w http.ResponseWriter, req *http.Request) {

	// get list of individuals
	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	iList, err := iamClient.Individuals().List(req.Context(), iam.IndividualListOptions{})
	if err != nil {
		s.Error(w, err)
		return
	}

	// get list of cases and case types
	cmsClient, err := s.CMSClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	cList, err := cmsClient.Cases().List(req.Context(), cms.CaseListOptions{})
	if err != nil {
		s.Error(w, err)
		return
	}

	ctList, err := cmsClient.CaseTypes().List(req.Context(), cms.CaseTypeListOptions{})
	if err != nil {
		s.Error(w, err)
		return
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "reporting", map[string]interface{}{
		"Individuals":                   iList,
		"Cases":                         cList,
		"CaseTypes":                     ctList,
		"GenderAttributeID":             iam.GenderAttribute.ID,
		"DisplacementStatusAttributeID": iam.DisplacementStatusAttribute.ID,
		"BirthDateAttributeID":          iam.BirthDateAttribute.ID,
		"BeneficiaryPartyTypeID":        iam.BeneficiaryPartyType.ID,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
