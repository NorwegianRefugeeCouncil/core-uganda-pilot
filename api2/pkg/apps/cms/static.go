package cms

import "github.com/nrc-no/core-kafka/pkg/teams"

var GenderViolence = CaseType{
	ID:          "2ab2aba2-c18f-4235-9ccd-52161defca5b",
	Name:        "Gender Violence",
	PartyTypeID: "a842e7cb-3777-423a-9478-f1348be3b4a5",
	TeamID:      teams.KampalaResponseTeam.ID,
}

var Childcare = CaseType{
	ID:          "73f47b43-eaa3-4ece-af91-0a72ff4c742e",
	Name:        "Childcare",
	PartyTypeID: "a842e7cb-3777-423a-9478-f1348be3b4a5",
	TeamID:      teams.NairobiICLATeam.ID,
}

var HousingRights = CaseType{
	ID:          "9a5ee26f-8df3-447c-a4b0-ed7f36710d95",
	Name:        "Housing Rights",
	PartyTypeID: "d38a7085-7dff-4730-8be1-7c9d92a20cc3",
	TeamID:      teams.KampalaICLATeam.ID,
}

var FinancialAssistInd = CaseType{
	ID:          "4b37e5d0-56e7-48b3-8227-bed8ce72019a",
	Name:        "Financial Assistance",
	PartyTypeID: "d38a7085-7dff-4730-8be1-7c9d92a20cc3",
	TeamID:      teams.KampalaResponseTeam.ID,
}

var FinancialAssistHH = CaseType{
	ID:          "e8a9a733-c6c9-46aa-ad32-23ed57ec8c58",
	Name:        "Financial Assistance",
	PartyTypeID: "a842e7cb-3777-423a-9478-f1348be3b4a5",
	TeamID:      teams.NairobiResponseTeam.ID,
}
