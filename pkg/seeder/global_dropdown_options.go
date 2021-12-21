package seeder

import "github.com/nrc-no/core/pkg/api/types"

// GLobal options
var (
	yesNoChoice = []*types.SelectOption{
		{Name: "Yes"},
		{Name: "No"},
	}
	wgShortSet = []*types.SelectOption{
		{Name: "Moderate Impairment"},
		{Name: "Severe Impairment"},
	}
	globalDisplacementStatuses = []*types.SelectOption{
		{Name: "Refugee"},
		{Name: "Internally Displaced Person (DP)"},
		{Name: "Host Community"},
		{Name: "Other"},
	}
	globalGenders = []*types.SelectOption{
		{Name: "Male"},
		{Name: "Female"},
		{Name: "Non-Binary"},
		{Name: "Other"},
	}
)
