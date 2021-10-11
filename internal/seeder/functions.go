package seeder

import (
	"github.com/nrc-no/core/internal/form"
	"github.com/nrc-no/core/pkg/cms"
	iam2 "github.com/nrc-no/core/pkg/iam"
)

func caseType(id, name, partyTypeID, teamID string, form form.Form, intakeCaseType bool) cms.CaseType {
	ct := cms.CaseType{
		ID:             id,
		Name:           name,
		PartyTypeID:    partyTypeID,
		TeamID:         teamID,
		Form:           form,
		IntakeCaseType: intakeCaseType,
	}
	caseTypes = append(caseTypes, ct)
	return ct
}

func team(id, name string) iam2.Team {
	t := iam2.Team{
		ID:   id,
		Name: name,
	}
	teams = append(teams, t)
	return t
}

func country(id, name string) iam2.Country {
	t := iam2.Country{
		ID:   id,
		Name: name,
	}
	countries = append(countries, t)
	return t
}

func individual(id string, fullName string, displayName string, birthDate string, email string, displacementStatus string, gender string, consent string, consentProof string, anonymous string, minor string, protectionConcerns string, physicalImpairment string, physicalImpairmentIntensity string, sensoryImpairment string, sensoryImpairmentIntensity string, mentalImpairment string, mentalImpairmentIntensity string, nationality string, spokenLanguages string, preferredLanguage string, physicalAddress string, primaryPhoneNumber string, secondaryPhoneNumber string, preferredMeansOfContact string, requireAnInterpreter string) iam2.Individual {
	var i = iam2.Individual{
		Party: &iam2.Party{
			ID: id,
			PartyTypeIDs: []string{
				iam2.IndividualPartyType.ID,
			},
			Attributes: map[string][]string{
				iam2.FullNameAttribute.ID:                    {fullName},
				iam2.DisplayNameAttribute.ID:                 {displayName},
				iam2.EMailAttribute.ID:                       {email + "@email.com"},
				iam2.BirthDateAttribute.ID:                   {birthDate},
				iam2.DisplacementStatusAttribute.ID:          {displacementStatus},
				iam2.GenderAttribute.ID:                      {gender},
				iam2.ConsentToNrcDataUseAttribute.ID:         {consent},
				iam2.ConsentToNrcDataUseProofAttribute.ID:    {consentProof},
				iam2.AnonymousAttribute.ID:                   {anonymous},
				iam2.MinorAttribute.ID:                       {minor},
				iam2.ProtectionConcernsAttribute.ID:          {protectionConcerns},
				iam2.PhysicalImpairmentAttribute.ID:          {physicalImpairment},
				iam2.PhysicalImpairmentIntensityAttribute.ID: {physicalImpairmentIntensity},
				iam2.SensoryImpairmentAttribute.ID:           {sensoryImpairment},
				iam2.SensoryImpairmentIntensityAttribute.ID:  {sensoryImpairmentIntensity},
				iam2.MentalImpairmentAttribute.ID:            {mentalImpairment},
				iam2.MentalImpairmentIntensityAttribute.ID:   {mentalImpairmentIntensity},
				iam2.UGNationalityAttribute.ID:               {nationality},
				iam2.UGSpokenLanguagesAttribute.ID:           {spokenLanguages},
				iam2.UGPreferredLanguageAttribute.ID:         {preferredLanguage},
				iam2.UGPhysicalAddressAttribute.ID:           {physicalAddress},
				iam2.PrimaryPhoneNumberAttribute.ID:          {primaryPhoneNumber},
				iam2.SecondaryPhoneNumberAttribute.ID:        {secondaryPhoneNumber},
				iam2.UGPreferredMeansOfContactAttribute.ID:   {preferredMeansOfContact},
				iam2.UGRequireAnInterpreterAttribute.ID:      {requireAnInterpreter},
			},
		},
	}
	individuals = append(individuals, i)
	return i
}

func ugandaIndividual(
	individual iam2.Individual,
	identificationDate string,
	identificationLocation string,
	identificationSource string,
	admin2 string,
	admin3 string,
	admin4 string,
	admin5 string,
) iam2.Individual {
	individual.Attributes.Add(iam2.UGIdentificationDateAttribute.ID, identificationDate)
	individual.Attributes.Add(iam2.UGIdentificationLocationAttribute.ID, identificationLocation)
	individual.Attributes.Add(iam2.UGIdentificationSourceAttribute.ID, identificationSource)
	individual.Attributes.Add(iam2.UGAdmin2Attribute.ID, admin2)
	individual.Attributes.Add(iam2.UGAdmin3Attribute.ID, admin3)
	individual.Attributes.Add(iam2.UGAdmin4Attribute.ID, admin4)
	individual.Attributes.Add(iam2.UGAdmin5Attribute.ID, admin5)
	return individual
}

func staff(individual iam2.Individual) iam2.Individual {
	individual.AddPartyType(iam2.StaffPartyType.ID)
	return individual
}

func beneficiary(individual iam2.Individual) iam2.Individual {
	individual.AddPartyType(iam2.BeneficiaryPartyType.ID)
	return individual
}

func membership(id string, individual iam2.Individual, team iam2.Team) iam2.Membership {
	m := iam2.Membership{
		ID:           id,
		TeamID:       team.ID,
		IndividualID: individual.ID,
	}
	memberships = append(memberships, m)
	return m
}

func nationality(id string, team iam2.Team, country iam2.Country) iam2.Nationality {
	m := iam2.Nationality{
		ID:        id,
		CountryID: country.ID,
		TeamID:    team.ID,
	}
	nationalities = append(nationalities, m)
	return m
}

func kase(id, createdByID, partyID, teamID string, caseType cms.CaseType, done, intakeCase bool, formData map[string][]string) cms.Case {
	k := cms.Case{
		ID:         id,
		CaseTypeID: caseType.ID,
		CreatorID:  createdByID,
		PartyID:    partyID,
		TeamID:     teamID,
		Done:       done,
		Form:       caseType.Form,
		FormData:   formData,
		IntakeCase: intakeCase,
	}
	cases = append(cases, k)
	return k
}

func identificationDocumentType(id, name string) iam2.IdentificationDocumentType {
	idt := iam2.IdentificationDocumentType{
		ID:   id,
		Name: name,
	}
	identificationDocumentTypes = append(identificationDocumentTypes, idt)
	return idt
}

func identificationDocument(id, partyId, documentNumber, identificationDocumentTypeId string) iam2.IdentificationDocument {
	newId := iam2.IdentificationDocument{
		ID:                           id,
		PartyID:                      partyId,
		DocumentNumber:               documentNumber,
		IdentificationDocumentTypeID: identificationDocumentTypeId,
	}
	identificationDocuments = append(identificationDocuments, newId)
	return newId
}
