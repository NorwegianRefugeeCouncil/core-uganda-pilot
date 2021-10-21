package seeder

import (
	"github.com/nrc-no/core/pkg/cms"
	"github.com/nrc-no/core/pkg/form"
	"github.com/nrc-no/core/pkg/iam"
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

func team(id, name string) iam.Team {
	t := iam.Team{
		ID:   id,
		Name: name,
	}
	teams = append(teams, t)
	return t
}

func country(id, name string) iam.Country {
	t := iam.Country{
		ID:   id,
		Name: name,
	}
	countries = append(countries, t)
	return t
}

func individual(id string, fullName string, displayName string, birthDate string, email string, displacementStatus string, gender string, consent string, consentProof string, anonymous string, minor string, protectionConcerns string, physicalImpairment string, physicalImpairmentIntensity string, sensoryImpairment string, sensoryImpairmentIntensity string, mentalImpairment string, mentalImpairmentIntensity string, nationality string, spokenLanguages string, preferredLanguage string, physicalAddress string, primaryPhoneNumber string, secondaryPhoneNumber string, preferredMeansOfContact string, requireAnInterpreter string) iam.Individual {
	var i = iam.Individual{
		Party: &iam.Party{
			ID: id,
			PartyTypeIDs: []string{
				iam.IndividualPartyType.ID,
			},
			Attributes: map[string][]string{
				iam.FullNameAttribute.Alias:                    {fullName},
				iam.DisplayNameAttribute.Alias:                 {displayName},
				iam.EMailAttribute.Alias:                       {email + "@email.com"},
				iam.BirthDateAttribute.Alias:                   {birthDate},
				iam.DisplacementStatusAttribute.Alias:          {displacementStatus},
				iam.GenderAttribute.Alias:                      {gender},
				iam.ConsentToNrcDataUseAttribute.Alias:         {consent},
				iam.ConsentToNrcDataUseProofAttribute.Alias:    {consentProof},
				iam.AnonymousAttribute.Alias:                   {anonymous},
				iam.MinorAttribute.Alias:                       {minor},
				iam.ProtectionConcernsAttribute.Alias:          {protectionConcerns},
				iam.PhysicalImpairmentAttribute.Alias:          {physicalImpairment},
				iam.PhysicalImpairmentIntensityAttribute.Alias: {physicalImpairmentIntensity},
				iam.SensoryImpairmentAttribute.Alias:           {sensoryImpairment},
				iam.SensoryImpairmentIntensityAttribute.Alias:  {sensoryImpairmentIntensity},
				iam.MentalImpairmentAttribute.Alias:            {mentalImpairment},
				iam.MentalImpairmentIntensityAttribute.Alias:   {mentalImpairmentIntensity},
				iam.UGNationalityAttribute.Alias:               {nationality},
				iam.UGSpokenLanguagesAttribute.Alias:           {spokenLanguages},
				iam.UGPreferredLanguageAttribute.Alias:         {preferredLanguage},
				iam.UGPhysicalAddressAttribute.Alias:           {physicalAddress},
				iam.PrimaryPhoneNumberAttribute.Alias:          {primaryPhoneNumber},
				iam.SecondaryPhoneNumberAttribute.Alias:        {secondaryPhoneNumber},
				iam.UGPreferredMeansOfContactAttribute.Alias:   {preferredMeansOfContact},
				iam.UGRequireAnInterpreterAttribute.Alias:      {requireAnInterpreter},
			},
		},
	}
	individuals = append(individuals, i)
	return i
}

func ugandaIndividual(
	individual iam.Individual,
	identificationDate string,
	identificationLocation string,
	identificationSource string,
	admin2 string,
	admin3 string,
	admin4 string,
	admin5 string,
) iam.Individual {
	individual.Attributes.Add(iam.UGIdentificationDateAttribute.Alias, identificationDate)
	individual.Attributes.Add(iam.UGIdentificationLocationAttribute.Alias, identificationLocation)
	individual.Attributes.Add(iam.UGIdentificationSourceAttribute.Alias, identificationSource)
	individual.Attributes.Add(iam.UGAdmin2Attribute.Alias, admin2)
	individual.Attributes.Add(iam.UGAdmin3Attribute.Alias, admin3)
	individual.Attributes.Add(iam.UGAdmin4Attribute.Alias, admin4)
	individual.Attributes.Add(iam.UGAdmin5Attribute.Alias, admin5)
	return individual
}

func staff(individual iam.Individual) iam.Individual {
	individual.AddPartyType(iam.StaffPartyType.ID)
	return individual
}

func beneficiary(individual iam.Individual) iam.Individual {
	individual.AddPartyType(iam.BeneficiaryPartyType.ID)
	return individual
}

func membership(id string, individual iam.Individual, team iam.Team) iam.Membership {
	m := iam.Membership{
		ID:           id,
		TeamID:       team.ID,
		IndividualID: individual.ID,
	}
	memberships = append(memberships, m)
	return m
}

func nationality(id string, team iam.Team, country iam.Country) iam.Nationality {
	m := iam.Nationality{
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

func identificationDocumentType(id, name string) iam.IdentificationDocumentType {
	idt := iam.IdentificationDocumentType{
		ID:   id,
		Name: name,
	}
	identificationDocumentTypes = append(identificationDocumentTypes, idt)
	return idt
}

func identificationDocument(id, partyID, documentNumber, identificationDocumentTypeId string) iam.IdentificationDocument {
	newID := iam.IdentificationDocument{
		ID:                           id,
		PartyID:                      partyID,
		DocumentNumber:               documentNumber,
		IdentificationDocumentTypeID: identificationDocumentTypeId,
	}
	identificationDocuments = append(identificationDocuments, newID)
	return newID
}
