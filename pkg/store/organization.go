package store

import (
	"context"
	"errors"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	uuid "github.com/satori/go.uuid"
)

type Permission string

type Organization struct {
	ID          string
	Name        string
	EmailDomain string
	Key         string
}

type OrganizationStore interface {
	Get(ctx context.Context, id string) (*types.Organization, error)
	Create(ctx context.Context, organiztion *types.Organization) (*types.Organization, error)
	Update(ctx context.Context, organization *types.Organization) (*types.Organization, error)
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

func (o organizationStore) Update(ctx context.Context, organization *types.Organization) (*types.Organization, error) {
	db, err := o.db.Get()
	if err != nil {
		return nil, err
	}

	qry := db.WithContext(ctx).
		Where("id = ?", organization.ID).
		Model(&Organization{}).
		Updates(map[string]interface{}{
			"name": organization.Name,
		})
	if qry.Error != nil {
		return nil, meta.NewInternalServerError(qry.Error)
	}

	if qry.RowsAffected == 0 {
		return nil, meta.NewInternalServerError(errors.New("record not found"))
	}

	return o.Get(ctx, organization.ID)
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
	if orgList == nil {
		orgList = []*Organization{}
	}
	return mapAllToOrg(orgList), nil
}
