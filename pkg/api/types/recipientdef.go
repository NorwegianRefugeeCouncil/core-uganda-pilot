package types

// RecipientDef is a special type of FormDefinition that represents a Case Recipient.
// A Case Recipient is... the recipient of a Case, such as a Household, Individual or other.
type RecipientDef struct {
	FormDefinition

	// Extends represents the Parent of the Recipient. Sometimes, we need to have multiple "profiles"
	// for case recipients. For example, Uganda Beneficiary Individual, Kenya Beneficiary Individual,
	// Colombia BeneficiaryIndividual are 3 different recipients that "extend" the same "Global Beneficiary Individual"
	// In this case, Uganda Beneficiary Individual would have Global Beneficiary Individual as Parent
	//
	// Also, in some cases, the parent might be a custom FormDefinition that is not a RecipientDef.
	// If the user has already defined a custom free-form, eg. Individual, and wants to add Case Management
	// to this currently existing form, then the Extends property would reference the Individual FormDefinition.
	Extends *FormRef `json:"extends,omitempty"`
}
