package iam

import "github.com/nrc-no/core/pkg/rest"

type ClientSet struct {
	c *rest.Client
}

var _ Interface = &ClientSet{}

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

func (c ClientSet) PartyAttributeDefinitions() PartyAttributeDefinitionClient {
	return &RESTAttributeClient{
		c: c.c,
	}
}

func (c ClientSet) Teams() TeamClient {
	return &RESTTeamClient{
		c: c.c,
	}
}

func (c ClientSet) Countries() CountryClient {
	return &RESTCountryClient{
		c: c.c,
	}
}

func (c ClientSet) Memberships() MembershipClient {
	return &RESTMembershipClient{
		c: c.c,
	}
}

func (c ClientSet) Nationalities() NationalityClient {
	return &RESTNationalityClient{
		c: c.c,
	}
}

func (c ClientSet) Individuals() IndividualClient {
	return &RESTIndividualClient{
		c: c.c,
	}
}

func (c ClientSet) IdentificationDocumentTypes() IdentificationDocumentTypeClient {
	return &RESTIdentificationDocumentTypeClient{
		c: c.c,
	}
}

func (c ClientSet) IdentificationDocuments() IdentificationDocumentClient {
	return &RESTIdentificationDocumentClient{
		c: c.c,
	}
}
