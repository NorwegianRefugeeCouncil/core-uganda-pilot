package webapp

import (
	"fmt"
	"github.com/nrc-no/core/pkg/iam"
	"github.com/nrc-no/core/pkg/registrationctrl"
	"github.com/nrc-no/core/pkg/seeder"
	"net/http"
)

func (s *Server) GetRegistrationController(w http.ResponseWriter, req *http.Request, individual *iam.Individual) (*registrationctrl.RegistrationController, error) {
	var individualId string

	if individual.ID == "new" || individual.ID == "" {
		if !s.GetPathParam("id", w, req, &individualId) {
			return nil, fmt.Errorf("cannot find id in path")
		}
	}
	irh := NewIndividualRegistrationHandler(s, individual, req)

	return registrationctrl.NewRegistrationController(irh, seeder.UgandaRegistrationFlow), nil
}
