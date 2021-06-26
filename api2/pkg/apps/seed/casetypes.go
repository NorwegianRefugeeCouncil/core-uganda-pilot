package seed

import (
	"github.com/nrc-no/core-kafka/pkg/apps/cms"
)

var GenderViolence = cms.CaseType{
	ID:          "2ab2aba2-c18f-4235-9ccd-52161defca5b",
	Name:        "Gender Violence",
	PartyTypeID: "a842e7cb-3777-423a-9478-f1348be3b4a5",
	TeamID:      KampalaResponseTeam.ID,
}

var Childcare = cms.CaseType{
	ID:          "73f47b43-eaa3-4ece-af91-0a72ff4c742e",
	Name:        "Childcare",
	PartyTypeID: "a842e7cb-3777-423a-9478-f1348be3b4a5",
	TeamID:      NairobiICLATeam.ID,
}

var HousingRights = cms.CaseType{
	ID:          "9a5ee26f-8df3-447c-a4b0-ed7f36710d95",
	Name:        "Housing Rights",
	PartyTypeID: "d38a7085-7dff-4730-8be1-7c9d92a20cc3",
	TeamID:      KampalaICLATeam.ID,
}

var FinancialAssistInd = cms.CaseType{
	ID:          "4b37e5d0-56e7-48b3-8227-bed8ce72019a",
	Name:        "Financial Assistance",
	PartyTypeID: "d38a7085-7dff-4730-8be1-7c9d92a20cc3",
	TeamID:      KampalaResponseTeam.ID,
}

var FinancialAssistHH = cms.CaseType{
	ID:          "e8a9a733-c6c9-46aa-ad32-23ed57ec8c58",
	Name:        "Financial Assistance",
	PartyTypeID: "a842e7cb-3777-423a-9478-f1348be3b4a5",
	TeamID:      NairobiResponseTeam.ID,
}
