package login

// Credential represents a password hash for a given Party
type Credential struct {
	PartyID string `json:"partyId" bson:"partyId"`
	Hash    string `json:"hash" bson:"hash"`
}
