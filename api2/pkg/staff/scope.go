package staff

import (
	"github.com/nrc-no/core-kafka/pkg/organizations"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
)

// Scope represents the operating scope of an Individual that has
// a Staff relationship with an Organization
type Scope struct {
	ID           string
	Organization *organizations.Organization
	Individual   *parties.Party
}
