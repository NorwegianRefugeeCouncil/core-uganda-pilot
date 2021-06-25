package auth

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Credential represents a password hash for a given Party
type Credential struct {
	PartyID string `json:"partyId" bson:"partyId"`
	Hash    string `json:"hash" bson:"hash"`
}

// CredentialsClient allows to set and verify passwords
type CredentialsClient struct {
	collection *mongo.Collection
}

func NewCredentialsClient(databaseName string, mongoClient *mongo.Client) *CredentialsClient {
	return &CredentialsClient{
		collection: mongoClient.Database(databaseName).Collection("credentials"),
	}
}

// VerifyPassword will check that the given password matches the stored password hash
func (c *CredentialsClient) VerifyPassword(ctx context.Context, partyID string, password string) bool {

	result := c.collection.FindOne(ctx, bson.M{
		"partyId": partyID,
	})
	if result.Err() != nil {
		return false
	}
	var cred Credential
	if err := result.Decode(&cred); err != nil {
		return false
	}

	if !comparePasswords(cred.Hash, []byte(password)) {
		return false
	} else {
		return true
	}
}

// SetPassword will set the Party password
func (c *CredentialsClient) SetPassword(ctx context.Context, partyID string, password string) error {
	saltedHash, err := hashAndSalt([]byte(password))
	if err != nil {
		return err
	}
	_, err = c.collection.UpdateOne(ctx, bson.M{
		"partyId": partyID,
	}, bson.M{
		"$set": bson.M{
			"partyId": partyID,
			"hash":    saltedHash,
		},
	}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

// hashAndSalt uses bcrypt algorithm to compute a salt + hash of a password
// Used to mitigate bruteforce attacks, rainbow tables, etc.
// The 15 parameter is the optimal bcrypt "cost" for speed and difficulty of bruteforce
// Hashing a password with this cost takes about 1.5 seconds
func hashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, 15)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// comparePasswords checks that the given password hash matches the given password
func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}
	return true
}
