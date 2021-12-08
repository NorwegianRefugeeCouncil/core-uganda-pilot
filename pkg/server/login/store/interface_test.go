package store

import (
	"github.com/nrc-no/core/pkg/mocks"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/utils/pointers"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

var identity1 = Identity{
	ID:    "identity1",
	State: IdentityStateActive,
	Credentials: []*Credential{
		{
			ID:     "credential1",
			Kind:   OidcCredential,
			Issuer: pointers.String("provider1"),
			Identifiers: []*CredentialIdentifier{
				{
					ID:         "identifier1",
					Identifier: "identifier-one",
				},
			},
		},
	},
}
var identity2 = Identity{
	ID:    "identity2",
	State: IdentityStateActive,
	Credentials: []*Credential{
		{
			ID:     "credential2",
			Kind:   OidcCredential,
			Issuer: pointers.String("provider2"),
			Identifiers: []*CredentialIdentifier{
				{
					ID:         "identifier2a",
					Identifier: "identifier-two-a",
				}, {
					ID:         "identifier2b",
					Identifier: "identifier-two-b",
				},
			},
		},
	},
}
var identity3 = Identity{
	ID:    "identity3",
	State: IdentityStateActive,
	Credentials: []*Credential{
		{
			ID:     "credential3a",
			Kind:   OidcCredential,
			Issuer: pointers.String("provider3a"),
			Identifiers: []*CredentialIdentifier{
				{
					ID:         "identifier3a",
					Identifier: "identifier-three-a",
				},
			},
		}, {
			ID:     "credential3b",
			Kind:   OidcCredential,
			Issuer: pointers.String("provider3b"),
			Identifiers: []*CredentialIdentifier{
				{
					ID:         "identifier3b",
					Identifier: "identifier-three-b",
				},
			},
		},
	},
}

func TestFindOidcIdentifier(t *testing.T) {

	s, _ := getStore(t)

	id, err := s.FindOidcIdentifier("identifier-one", "provider1")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "identifier1", id.ID)
	assert.Equal(t, "identifier-one", id.Identifier)
	assert.Equal(t, "credential1", id.Credential.ID)
	assert.Equal(t, pointers.String("provider1"), id.Credential.Issuer)
	assert.Equal(t, "identity1", id.Credential.IdentityID)
	assert.Equal(t, "identity1", id.Credential.Identity.ID)
	assert.Equal(t, IdentityStateActive, id.Credential.Identity.State)

}

func TestFindOidcIdentifierWhenUserHasMany(t *testing.T) {
	s, _ := getStore(t)
	id, err := s.FindOidcIdentifier("identifier-two-b", "provider2")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "identifier2b", id.ID)
}

func TestFindOidcIdentifierWhenUserHasMultipleProviders(t *testing.T) {
	s, _ := getStore(t)
	id, err := s.FindOidcIdentifier("identifier-three-b", "provider3b")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "identifier3b", id.ID)
}

func getStore(t *testing.T) (loginStore, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open("test.db"))
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}
	prepDb(t, db)
	return loginStore{db: mocks.NewMockFactory(db)}, db
}

func prepDb(t *testing.T, db *gorm.DB) {
	if err := db.AutoMigrate(
		&Identity{},
		&Credential{},
		&CredentialIdentifier{},
		&store.IdentityProvider{},
		&store.Organization{}); !assert.NoError(t, err) {
		t.Fatal(err)
	}
	db.Delete(&Identity{}, "id = id")
	db.Delete(&Credential{}, "id = id")
	db.Delete(&CredentialIdentifier{}, "id = id")
	db.Delete(&store.IdentityProvider{}, "id = id")
	db.Delete(&store.Organization{}, "id = id")

	var identities = []Identity{
		identity1,
		identity2,
		identity3,
	}
	err := db.Create(identities).Error
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}
}
