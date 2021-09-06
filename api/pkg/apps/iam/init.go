package iam

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) Init(ctx context.Context) error {

	if err := s.initPartyType(ctx); err != nil {
		return err
	}

	if err := s.initRelationshipType(ctx); err != nil {
		return err
	}

	if err := s.initAttribute(ctx); err != nil {
		return err
	}

	return nil

}

func (s *Server) initPartyType(ctx context.Context) error {
	for _, partyType := range []PartyType{
		IndividualPartyType,
		HouseholdPartyType,
		TeamPartyType,
		StaffPartyType,
	} {
		if err := s.partyTypeStore.Create(ctx, &partyType); err != nil {
			if !mongo.IsDuplicateKeyError(err) {
				return err
			}
			if err := s.partyTypeStore.Update(ctx, &partyType); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Server) initRelationshipType(ctx context.Context) error {
	for _, relationshipType := range []RelationshipType{
		HeadOfHouseholdRelationshipType,
		SpousalRelationshipType,
		SiblingRelationshipType,
		ParentalRelationshipType,
		MembershipRelationshipType,
	} {
		if err := s.relationshipTypeStore.Create(ctx, &relationshipType); err != nil {
			if !mongo.IsDuplicateKeyError(err) {
				return err
			}
			if err := s.relationshipTypeStore.Update(ctx, &relationshipType); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Server) initAttribute(ctx context.Context) error {
	for _, attribute := range []PartyAttributeDefinition{
		// "Built-in" Attributes
		FullNameAttribute,
		DisplayNameAttribute,
		BirthDateAttribute,
		EMailAttribute,
		TeamNameAttribute,

		// Individual Attributes
		DisplacementStatusAttribute,
		GenderAttribute,
		ConsentToNrcDataUseAttribute,
		ConsentToNrcDataUseProofAttribute,
		AnonymousAttribute,
		MinorAttribute,
		ProtectionConcernsAttribute,
		PhysicalImpairmentAttribute,
		PhysicalImpairmentIntensityAttribute,
		SensoryImpairmentAttribute,
		SensoryImpairmentIntensityAttribute,
		MentalImpairmentAttribute,
		MentalImpairmentIntensityAttribute,
		NationalityAttribute,
		SpokenLanguagesAttribute,
		PreferredLanguageAttribute,
		PhysicalAddressAttribute,
		PrimaryPhoneNumberAttribute,
		SecondaryPhoneNumberAttribute,
		PreferredMeansOfContactAttribute,
		RequireAnInterpreterAttribute,
		// -- End of Individual Attributes

		// Customisation for Uganda Demo
		IdentificationDateAttribute,
		IdentificationLocationAttribute,
		IdentificationSourceAttribute,
		Admin2Attribute,
		Admin3Attribute,
		Admin4Attribute,
		Admin5Attribute,
	} {
		if err := s.partyAttributeDefinitionStore.create(ctx, &attribute); err != nil {
			if !mongo.IsDuplicateKeyError(err) {
				return err
			}
			if err := s.partyAttributeDefinitionStore.update(ctx, &attribute); err != nil {
				return err
			}
		}
	}
	return nil
}
