package seed

import (
	"github.com/nrc-no/core-kafka/pkg/apps/cms"
)

var DomesticAbuse = cms.Case{
	ID:          "763bcc10-444a-4853-b3d6-13a1224a77ec",
	CaseTypeID:  "2ab2aba2-c18f-4235-9ccd-52161defca5b",
	PartyID:     "ab7a1620-f34e-4811-8534-853167ed7944",
	TeamID:      KampalaResponseTeam.ID,
	Description: "Domestic abuse",
	Done:        false,
}

var MonthlyAllowance = cms.Case{
	ID:          "47499762-c189-4a74-9156-7969f899073b",
	CaseTypeID:  "4b37e5d0-56e7-48b3-8227-bed8ce72019a",
	PartyID:     "40b30fb0-c392-4798-9400-bda3e5837867",
	TeamID:      KampalaResponseTeam.ID,
	Description: "Monthly allowance",
	Done:        false,
}

var ChildCare = cms.Case{
	ID:          "8fb5f755-85eb-4d91-97a9-fdf86c01df25",
	CaseTypeID:  "73f47b43-eaa3-4ece-af91-0a72ff4c742e",
	PartyID:     "40b30fb0-c392-4798-9400-bda3e5837867",
	TeamID:      KampalaResponseTeam.ID,
	Description: "Monthly stipend for Bo Diddley's child",
	Done:        true,
}
