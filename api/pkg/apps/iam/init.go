package iam

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) Init(ctx context.Context) error {

	for _, partyType := range []PartyType{
		IndividualPartyType,
		HouseholdPartyType,
		TeamPartyType,
		StaffPartyType,
	} {
		if err := s.PartyTypeStore.Create(ctx, &partyType); err != nil {
			if !mongo.IsDuplicateKeyError(err) {
				return err
			}
			if err := s.PartyTypeStore.Update(ctx, &partyType); err != nil {
				return err
			}
		}
	}

	for _, relationshipType := range []RelationshipType{
		HeadOfHouseholdRelationshipType,
		SpousalRelationshipType,
		SiblingRelationshipType,
		ParentalRelationshipType,
		MembershipRelationshipType,
	} {
		if err := s.RelationshipTypeStore.Create(ctx, &relationshipType); err != nil {
			if !mongo.IsDuplicateKeyError(err) {
				return err
			}
			if err := s.RelationshipTypeStore.Update(ctx, &relationshipType); err != nil {
				return err
			}
		}
	}

	for _, attribute := range []Attribute{
		FirstNameAttribute,
		LastNameAttribute,
		BirthDateAttribute,
		EMailAttribute,
		TeamNameAttribute,
	} {
		if err := s.AttributeStore.Create(ctx, &attribute); err != nil {
			if !mongo.IsDuplicateKeyError(err) {
				return err
			}
			if err := s.AttributeStore.Update(ctx, &attribute); err != nil {
				return err
			}
		}
	}

	return nil

}
