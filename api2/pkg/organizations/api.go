package organizations

import (
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
)

var PartyType = partytypes.PartyType{
	ID:        "09a7eef9-3f23-4c40-86f4-9b9440c56c6f",
	Name:      "Organization",
	IsBuiltIn: true,
}

type Organization struct {
	*parties.Party
}
