package teamorganizations

import (
	"github.com/nrc-no/core-kafka/pkg/organizations"
	"github.com/nrc-no/core-kafka/pkg/teams"
)

type Scope struct {
	ID           string
	Team         *teams.Team
	Organization *organizations.Organization
}
