package login

import (
	"context"
	"github.com/nrc-no/core/pkg/apps/iam"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) VerifyPassword(ctx context.Context, email, password string) (*iam.Individual, bool) {

	individualList, err := s.iam.Individuals().List(ctx, iam.IndividualListOptions{
		Attributes: map[string]string{
			iam.EMailAttribute.ID: email,
		},
	})
	if err != nil {
		return nil, false
	}

	if len(individualList.Items) == 0 {
		return nil, false
	}

	res := s.Collection.FindOne(ctx, bson.M{
		"partyId": individualList.Items[0].ID,
	})
	if res.Err() != nil {
		return nil, false
	}

	raw, err := res.DecodeBytes()
	if err != nil {
		return nil, false
	}

	hash, ok := raw.Lookup("hash").StringValueOK()
	if !ok {
		return nil, false
	}

	match := comparePasswords(hash, []byte(password))
	if !match {
		return nil, false
	}

	return individualList.Items[0], true

}

// SetPassword will set the Party password
func (s *Server) SetPassword(ctx context.Context, partyID string, password string) error {
	saltedHash, err := HashAndSalt(s.BCryptCost, []byte(password))
	if err != nil {
		return err
	}
	_, err = s.Collection.UpdateOne(ctx, bson.M{
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

// CreatePassword will create a new credential for the Party
func (s *Server) CreatePassword(ctx context.Context, partyID string, password string) error {
	saltedHash, err := HashAndSalt(s.BCryptCost, []byte(password))
	if err != nil {
		return err
	}

	var newCredential = Credential{
		PartyID: partyID,
		Hash:    saltedHash,
	}

	_, err = s.Collection.InsertOne(ctx, newCredential)
	if err != nil {
		return err
	}

	return nil
}

// HashAndSalt uses bcrypt algorithm to compute a salt + hash of a password
// Used to mitigate bruteforce attacks, rainbow tables, etc.
func HashAndSalt(cost int, pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, cost)
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
