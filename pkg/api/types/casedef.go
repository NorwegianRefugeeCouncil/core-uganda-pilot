package types

// CaseDef represents a special type of FormDefinition that represents a Case.
// A Case is an instance of a service given to the Recipient.
type CaseDef struct {
	FormDefinition

	// Recipient represents the type of recipient that this case is for.
	// For example, if the CaseDef is Colombia Individual Intake, then
	// the Recipient could be a link to a ColombiaIndividualRecipient form.
	Recipient FormRef `json:"recipient,omitempty"`
}
