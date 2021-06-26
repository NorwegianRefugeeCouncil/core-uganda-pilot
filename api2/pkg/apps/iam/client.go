package iam

import "github.com/nrc-no/core-kafka/pkg/rest"

type ClientSet struct {
	c *rest.Client
}

func NewClientSet(restConfig *rest.RESTConfig) *ClientSet {
	return &ClientSet{
		c: rest.NewClient(restConfig),
	}
}

func (c ClientSet) Parties() PartyClient {
	return &RESTPartyClient{
		c: c.c,
	}
}

func (c ClientSet) PartyTypes() PartyTypeClient {
	return &RESTPartyTypeClient{
		c: c.c,
	}
}

func (c ClientSet) Relationships() RelationshipClient {
	return &RESTRelationshipClient{
		c: c.c,
	}
}

func (c ClientSet) RelationshipTypes() RelationshipTypeClient {
	return &RESTRelationshipTypeClient{
		c: c.c,
	}
}

func (c ClientSet) Attributes() AttributeClient {
	return &RESTAttributeClient{
		c: c.c,
	}
}

func (c ClientSet) Teams() TeamClient {
	return &RESTTeamClient{
		c: c.c,
	}
}

func (c ClientSet) Organizations() OrganizationClient {
	return &RESTOrganizationClient{
		c: c.c,
	}
}

func (c ClientSet) Staff() StaffClient {
	return &RESTStaffClient{
		c: c.c,
	}
}

func (c ClientSet) Memberships() MembershipClient {
	return &RESTMembershipClient{
		c: c.c,
	}
}

func (c ClientSet) Individuals() IndividualClient {
	return &RESTIndividualClient{
		c: c.c,
	}
}

var _ Interface = &ClientSet{}
