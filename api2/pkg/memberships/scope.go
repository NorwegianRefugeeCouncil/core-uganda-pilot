package memberships

import (
	"github.com/nrc-no/core-kafka/pkg/staff"
	"github.com/nrc-no/core-kafka/pkg/teamorganizations"
)

type MemberScope struct {
	ID         string `json:"id"`
	TeamScope  *teamorganizations.Scope
	StaffScope *staff.Scope
}
