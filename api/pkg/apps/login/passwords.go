package login

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/sirupsen/logrus"
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
		logrus.WithError(err).WithField("email", email).Warnf("failed to list individuals")
		return nil, false
	}

	if len(individualList.Items) == 0 {
		logrus.Trace("no individuals found with email %s", email)
		return nil, false
	}

	credentialsCollection, err := s.credentialsCollectionFn()
	if err != nil {
		logrus.WithError(err).Errorf("failed to get credentials collection")
		return nil, false
	}

	partyID := individualList.Items[0].ID
	res := credentialsCollection.FindOne(ctx, bson.M{
		"partyId": partyID,
	})
	if res.Err() != nil {
		logrus.WithError(res.Err()).WithField("PartyID", partyID).Warnf("failed to find party")
		return nil, false
	}

	raw, err := res.DecodeBytes()
	if err != nil {
		logrus.WithError(err).Warnf("failed to decode bytes")
		return nil, false
	}

	hash, ok := raw.Lookup("hash").StringValueOK()
	if !ok {
		logrus.WithError(fmt.Errorf("could not convert mongo hash to string"))
		return nil, false
	}

	match := comparePasswords(hash, []byte(password))
	if !match {
		logrus.
			WithField("PartyID", partyID).
			WithField("Email", email).
			Tracef("wrong password provided")

		return nil, false
	}

	return individualList.Items[0], true

}

// SetPassword will set the Party password
func (s *Server) SetPassword(ctx context.Context, partyID string, password string) error {

	credentialsCollection, err := s.credentialsCollectionFn()
	if err != nil {
		logrus.WithError(err).Errorf("failed to get credentials collection")
		return err
	}

	saltedHash, err := HashAndSalt(s.BCryptCost, []byte(password))
	if err != nil {
		return fmt.Errorf("failed to hash and salt: %v", err)
	}

	_, err = credentialsCollection.UpdateOne(ctx, bson.M{
		"partyId": partyID,
	}, bson.M{
		"$set": bson.M{
			"partyId": partyID,
			"hash":    saltedHash,
		},
	}, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("failed to upsert password: %v", err)
	}
	return nil
}

// CreatePassword will create a new credential for the Party
func (s *Server) CreatePassword(ctx context.Context, partyID string, password string) error {

	credentialsCollection, err := s.credentialsCollectionFn()
	if err != nil {
		logrus.WithError(err).Errorf("failed to get credentials collection")
		return err
	}

	saltedHash, err := HashAndSalt(s.BCryptCost, []byte(password))
	if err != nil {
		return fmt.Errorf("failed to hash and salt: %v", err)
	}

	var newCredential = Credential{
		PartyID: partyID,
		Hash:    saltedHash,
	}

	_, err = credentialsCollection.InsertOne(ctx, newCredential)
	if err != nil {
		return fmt.Errorf("failed to insert credential: %v", err)
	}

	return nil
}

// HashAndSalt uses bcrypt algorithm to compute a salt + hash of a password
// Used to mitigate bruteforce attacks, rainbow tables, etc.
func HashAndSalt(cost int, pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, cost)
	if err != nil {
		return "", fmt.Errorf("failed to generate hash from password: %v", err)
	}
	return string(hash), nil
}

// comparePasswords checks that the given password hash matches the given password
func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	return err == nil
}
