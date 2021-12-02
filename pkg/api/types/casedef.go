package types

// Recipient represents the recipient of a CaseDef
// Introducing a struct here since it would be possible
// to have recipients that are not records in other forms in the future.
// For example, we might have recipients that are a number of people.
// (Extensibility)
type Recipient struct {
	// Form represents a recipient that is a record in another form
	Form *FormRef `json:"form,omitempty"`
}

// CaseDef represents a special type of FormDefinition that represents a Case.
// A Case is an instance of a service given to the Recipient.
type CaseDef struct {
	FormDefinition

	// Recipient represents the type of recipient that this case is for.
	// For example, if the CaseDef is Colombia Individual Intake, then
	// the Recipient could be a link to a ColombiaIndividualRecipient form.
	Recipient Recipient `json:"recipient,omitempty"`
}
