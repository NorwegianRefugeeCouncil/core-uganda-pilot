package store

import (
	"context"
	"errors"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/utils/pointers"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Interface interface {
	GetAuthRequest(ctx context.Context, id string) (*AuthRequest, error)
	CreateAuthRequest(ctx context.Context, authRequest *AuthRequest) (*AuthRequest, error)
	UpdateAuthRequest(ctx context.Context, authRequest *AuthRequest) (*AuthRequest, error)
	GetIdentity(ctx context.Context, id string) (*Identity, error)
	FindOidcIdentifier(identifier string, identityProviderId string) (*CredentialIdentifier, error)
	CreateOidcIdentity(issuer string, identifier string, initialAccessToken string, initialRefreshToken string, initialIdToken string) (*Identity, error)
	CreateOidcIdentityProfile(profile IdentityProfile) error
}

type loginStore struct {
	db store.Factory
}

func NewStore(db store.Factory) Interface {
	return &loginStore{db: db}
}

func (l *loginStore) GetAuthRequest(ctx context.Context, id string) (*AuthRequest, error) {
	db, err := l.db.Get()
	if err != nil {
		return nil, err
	}
	var authRequest AuthRequest
	authRequest.ID = uuid.NewV4().String()
	if err := db.WithContext(ctx).First(&authRequest, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &authRequest, nil
}

func (l *loginStore) UpdateAuthRequest(ctx context.Context, authRequest *AuthRequest) (*AuthRequest, error) {
	db, err := l.db.Get()
	if err != nil {
		return nil, err
	}
	if err := db.WithContext(ctx).Updates(authRequest).Error; err != nil {
		return nil, err
	}
	return authRequest, nil
}

func (l *loginStore) CreateAuthRequest(ctx context.Context, authRequest *AuthRequest) (*AuthRequest, error) {
	db, err := l.db.Get()
	if err != nil {
		return nil, err
	}
	authRequest.ID = uuid.NewV4().String()
	if err := db.WithContext(ctx).Create(authRequest).Error; err != nil {
		return nil, err
	}
	return authRequest, nil
}

func (l *loginStore) GetIdentity(ctx context.Context, id string) (*Identity, error) {

	db, err := l.db.Get()
	if err != nil {
		return nil, err
	}

	var identity *Identity
	qry := db.Model(&Identity{}).Preload("Credentials.Identifiers").First(&identity, "id = ?", id)
	if qry.Error != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, identityNotFoundError(id)
		} else {
			return nil, err
		}
	}

	for _, credential := range identity.Credentials {
		credential.Identity = identity
		for _, identifier := range credential.Identifiers {
			identifier.Credential = credential
		}
	}

	return identity, nil

}

func (l *loginStore) CreateOidcIdentity(
	issuer string,
	identifier string,
	initialAccessToken string,
	initialRefreshToken string,
	initialIdToken string,
) (*Identity, error) {
	identity := &Identity{
		ID:    uuid.NewV4().String(),
		State: IdentityStateActive,
		Credentials: []*Credential{
			{
				ID:     uuid.NewV4().String(),
				Kind:   OidcCredential,
				Issuer: pointers.String(issuer),
				Identifiers: []*CredentialIdentifier{
					{
						ID:         uuid.NewV4().String(),
						Identifier: identifier,
					},
				},
				InitialAccessToken:  pointers.String(initialAccessToken),
				InitialRefreshToken: pointers.String(initialRefreshToken),
				InitialIdToken:      pointers.String(initialIdToken),
			},
		},
	}
	db, err := l.db.Get()
	if err != nil {
		return nil, err
	}

	if err := db.Create(&identity).Error; err != nil {
		return nil, err
	}

	for _, credential := range identity.Credentials {
		credential.Identity = identity
		for _, identifier := range credential.Identifiers {
			identifier.Credential = credential
		}
	}

	return identity, nil

}

func (l *loginStore) CreateOidcIdentityProfile(
	profile IdentityProfile,
) error {

	db, err := l.db.Get()
	if err != nil {
		return err
	}

	if err := db.Create(&profile).Error; err != nil {
		return err
	}

	return nil
}

func (l *loginStore) FindOidcIdentifier(
	identifier string,
	issuer string,
) (*CredentialIdentifier, error) {

	db, err := l.db.Get()
	if err != nil {
		return nil, err
	}

	type Result struct {
		IdentityID   string
		IdentifierID string
	}

	var result Result
	var qry = db.Model(&Identity{}).Select(
		"credential_identifiers.id as identifier_id",
		"identities.id as identity_id",
	).
		Joins("left join credentials on credentials.identity_id = identities.id").
		Joins("left join credential_identifiers on credential_identifiers.credential_id = credentials.id").
		Where("credentials.kind = ? and credential_identifiers.identifier = ? and credentials.issuer = ?",
			OidcCredential,
			identifier,
			issuer,
		).
		First(&result)

	if qry.Error != nil {
		if errors.Is(qry.Error, gorm.ErrRecordNotFound) {
			return nil, meta.NewNotFound(meta.GroupResource{Resource: "credentialidentifiers", Group: "auth.nrc.no"}, "")
		} else {
			return nil, qry.Error
		}
	}

	var ident *Identity
	err = db.
		Model(&Identity{}).
		Preload("Credentials.Identifiers").
		First(&ident, "id = ?", result.IdentityID).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, identityNotFoundError(result.IdentityID)
		} else {
			return nil, err
		}
	}

	for _, credential := range ident.Credentials {
		credential.Identity = ident
		for _, identifier := range credential.Identifiers {
			identifier.Credential = credential
		}
	}

	var identif *CredentialIdentifier
	var found = false
	for _, credential := range ident.Credentials {
		if found {
			break
		}
		if credential.Kind != OidcCredential {
			continue
		}
		if credential.Issuer == nil || *credential.Issuer != issuer {
			continue
		}
		for _, credentialIdentifier := range credential.Identifiers {
			if credentialIdentifier.Identifier != identifier {
				continue
			}
			if credentialIdentifier.ID == result.IdentifierID {
				identif = credentialIdentifier
				found = true
				break
			}
		}
	}
	if identif == nil {
		return nil, meta.NewNotFound(meta.GroupResource{Resource: "credentialidentifiers", Group: "auth.nrc.no"}, result.IdentifierID)
	}

	return identif, nil

}

func identityNotFoundError(id string) error {
	return meta.NewNotFound(meta.GroupResource{Resource: "identities", Group: "auth.nrc.no"}, id)
}
