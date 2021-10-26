package store

import (
	"context"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/types"
	uuid "github.com/satori/go.uuid"
)

type Permission string

const (
	SuperAdmin = "super_admin"

	ReadDatabase   = "database_read"
	ManageDatabase = "database_manage"

	ReadFolder   = "folder_read"
	ManageFolder = "folder_manage"

	ReadForm   = "form_read"
	ManageForm = "form_manage"

	ReadRecords = "record_read"
	WriteRecord = "record_write"
)

type Organization struct {
	ID   string
	Name string
	Key  string
}

type OrganizationStore interface {
	Get(ctx context.Context, id string) (*types.Organization, error)
	Create(ctx context.Context, organiztion *types.Organization) (*types.Organization, error)
	List(ctx context.Context) ([]*types.Organization, error)
}

type organizationStore struct {
	db Factory
}

func NewOrganizationStore(db Factory) OrganizationStore {
	return &organizationStore{
		db: db,
	}
}

func mapFromOrg(org *types.Organization) *Organization {
	return &Organization{
		ID:   org.ID,
		Name: org.Name,
		Key:  org.Key,
	}
}

func mapToOrg(org *Organization) *types.Organization {
	return &types.Organization{
		ID:   org.ID,
		Name: org.Name,
		Key:  org.Key,
	}
}

func mapAllToOrg(orgs []*Organization) []*types.Organization {
	var result []*types.Organization
	for _, org := range orgs {
		result = append(result, mapToOrg(org))
	}
	if result == nil {
		result = []*types.Organization{}
	}
	return result
}

func (o organizationStore) Get(ctx context.Context, id string) (*types.Organization, error) {
	db, err := o.db.Get()
	if err != nil {
		return nil, err
	}

	var org Organization
	if err := db.WithContext(ctx).First(&org, "id = ?", id).Error; err != nil {
		return nil, meta.NewInternalServerError(err)
	}
	return mapToOrg(&org), nil
}

func (o organizationStore) Create(ctx context.Context, organiztion *types.Organization) (*types.Organization, error) {
	db, err := o.db.Get()
	if err != nil {
		return nil, err
	}

	var org = mapFromOrg(organiztion)
	org.ID = uuid.NewV4().String()
	if err := db.WithContext(ctx).Create(org).Error; err != nil {
		return nil, meta.NewInternalServerError(err)
	}
	return mapToOrg(org), nil
}

func (o organizationStore) List(ctx context.Context) ([]*types.Organization, error) {
	db, err := o.db.Get()
	if err != nil {
		return nil, err
	}

	var orgList []*Organization
	if err := db.WithContext(ctx).Find(&orgList).Error; err != nil {
		return nil, meta.NewInternalServerError(err)
	}
	return mapAllToOrg(orgList), nil
}

type User struct {
	ID                  string
	Email               string
	EmailVerified       bool
	PhoneNumber         string
	PhoneNumberVerified bool
}

type UserProfile struct {
	ID                  string
	UserID              string
	User                User
	Provider            string
	Email               string
	EmailVerified       bool
	PhoneNumber         string
	PhoneNumberVerified bool
}

type OrganizationRole struct {
	ID             string
	OrganizationID string
	Organization   Organization
	Name           string
}

type OrganizationRolePermission struct {
	ID                 string
	OrganizationID     string
	Organization       Organization
	OrganizationRoleID string
	OrganizationRole   OrganizationRole
	Permission         Permission
	Scope              string
}

type OrganizationUser struct {
	UserID         string
	OrganizationID string
}

type OrganizationDatabase struct {
	OrganizationID string
	DatabaseID     string
}
