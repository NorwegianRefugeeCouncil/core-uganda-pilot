package v1

type Beneficiary struct {
}

type BeneficiaryList struct {
}

type ConsentType string

const (
	ConsentURL ConsentType = "URL"
)

type BeneficiaryConsent struct {
	Type ConsentType `json:"type,omitempty"`
	URL  *string     `json:"url,omitempty"`
}

type BeneficiaryMajority string

const (
	Minor BeneficiaryMajority = "minor"
)

type BeneficiarySpec struct {
	Consent BeneficiaryConsent `json:"consent"`
}

type Household struct {
}

type HouseholdList struct {
}

type HouseholdSpec struct {
}
