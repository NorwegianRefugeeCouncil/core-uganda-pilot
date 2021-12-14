package store

import (
	"context"
	"errors"

	"github.com/nrc-no/core/pkg/api/types"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Why do I have to define this again for the migrations?
// Can we not import from types for migrations?
type Foo struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name       string    `gorm:"column:name;size:128;not null;"`
	OtherField int       `gorm:"column:other_field;not null"`
	Valid      bool      `gorm:"column:valid;not null;default:true"`
}

type FooStore interface {
	Get(ctx context.Context, fooId uuid.UUID) (*types.Foo, error)
	ListValid(ctx context.Context) ([]*types.Foo, error)
	Create(ctx context.Context, foo *types.Foo) (*types.Foo, error)
	Delete(ctx context.Context, fooId uuid.UUID) error
}

type fooStore struct {
	db Factory
}

func NewFooStore(db Factory) FooStore {
	return &fooStore{db: db}
}

func (d *fooStore) Get(ctx context.Context, fooId uuid.UUID) (*types.Foo, error) {
	_, db, l, done, err := actionContext(ctx, d.db, "foo", "get", zap.String("foo_id", fooId.String()))
	if err != nil {
		return nil, err
	}
	defer done()

	var f types.Foo
	if err := db.First(&f, fooId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.Error("Foo not found")
		} else {
			l.Error("Foo Get error")
		}
		return nil, err
	}

	return &f, nil
}

func (d *fooStore) ListValid(ctx context.Context) ([]*types.Foo, error) {
	_, db, l, done, err := actionContext(ctx, d.db, "foo", "list_valid")
	if err != nil {
		return nil, err
	}
	defer done()

	var f []*types.Foo
	if err := db.Where("valid IS TRUE").Find(&f).Error; err != nil {
		l.Error("Foo List error")
		return nil, err
	}

	return f, nil
}

func (d *fooStore) Create(ctx context.Context, foo *types.Foo) (*types.Foo, error) {
	_, db, l, done, err := actionContext(ctx, d.db, "foo", "create", zap.String("foo_id", foo.ID.String()))
	if err != nil {
		return nil, err
	}
	defer done()

	var f types.Foo
	if err := db.FirstOrCreate(&f, foo).Error; err != nil {
		l.Error("Foo Create error")
		return nil, err
	}

	return &f, nil
}

func (d *fooStore) Delete(ctx context.Context, fooId uuid.UUID) error {
	_, db, l, done, err := actionContext(ctx, d.db, "foo", "delete", zap.String("foo_id", fooId.String()))
	if err != nil {
		return err
	}
	defer done()

	if err := db.Delete(types.Foo{}, fooId).Error; err != nil {
		l.Error("Foo Delete error")
		return err
	}

	return nil
}
