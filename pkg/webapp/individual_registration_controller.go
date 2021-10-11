package webapp

import (
	"fmt"
	"github.com/nrc-no/core/internal/registrationctrl"
	"github.com/nrc-no/core/internal/seeder"
	"github.com/nrc-no/core/pkg/iam"
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
