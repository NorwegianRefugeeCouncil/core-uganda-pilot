package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/utils/sets"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type FolderStore interface {
	Get(ctx context.Context, folderID string) (*types.Folder, error)
	List(ctx context.Context) (*types.FolderList, error)
	Create(ctx context.Context, folder *types.Folder) (*types.Folder, error)
	Delete(ctx context.Context, id string) error
}

type Folder struct {
	ID         string
	DatabaseID string
	Database   *Database
	ParentID   string
	Parent     *Folder
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewFolderStore(db Factory) FolderStore {
	return &folderStore{db: db}
}

type folderStore struct {
	db Factory
}

var _ FolderStore = &folderStore{}

func (d *folderStore) Get(ctx context.Context, folderID string) (*types.Folder, error) {

	db, err := d.db.Get()
	if err != nil {
		return nil, err
	}

	var folder Folder
	if err := db.WithContext(ctx).First(&folder, folderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, meta.NewNotFound(meta.GroupResource{
				Group:    "nrc.no",
				Resource: "folders",
			}, folderID)
		} else {
			return nil, meta.NewInternalServerError(err)
		}
	}
	return mapFolderTo(&folder), nil
}

func (d *folderStore) List(ctx context.Context) (*types.FolderList, error) {

	db, err := d.db.Get()
	if err != nil {
		return nil, err
	}

	var folders []Folder
	if err := db.WithContext(ctx).Find(&folders).Error; err != nil {
		return nil, meta.NewInternalServerError(err)
	}
	var result = make([]*types.Folder, len(folders))
	for i, folder := range folders {
		result[i] = mapFolderTo(&folder)
	}
	if result == nil {
		result = []*types.Folder{}
	}
	return &types.FolderList{
		Items: result,
	}, nil
}

func (d *folderStore) Create(ctx context.Context, folder *types.Folder) (*types.Folder, error) {

	db, err := d.db.Get()
	if err != nil {
		return nil, err
	}

	folder.ID = uuid.NewV4().String()
	database := mapFolderFrom(folder)
	database.CreatedAt = time.Now()
	if err := db.WithContext(ctx).Create(database).Error; err != nil {
		return nil, meta.NewInternalServerError(err)
	}
	return mapFolderTo(database), nil
}

func (d *folderStore) Delete(ctx context.Context, id string) error {

	db, err := d.db.Get()
	if err != nil {
		return err
	}

	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var folderForms []*Form
		if err := tx.Find(&folderForms, "folder_id = ?", id).Error; err != nil {
			return meta.NewInternalServerError(err)
		}

		formIds := sets.NewString()
		fieldParams := sets.NewString()
		formIdParams := sets.NewString()
		rootFormIdParams := sets.NewString()
		for i := 0; i < len(folderForms); i++ {
			form := folderForms[i]
			formIds = formIds.Insert(form.ID)
			fieldParams = fieldParams.Insert("?")
			formIdParams = formIdParams.Insert("?")
			rootFormIdParams = rootFormIdParams.Insert("?")
		}

		fieldFormMatcher := fmt.Sprintf("root_form_id in (%s)", fieldParams.Join(","))
		if err := tx.Where(fieldFormMatcher, formIds.ListIntf()...).Delete(&Field{}).Error; err != nil {
			return meta.NewInternalServerError(err)
		}

		formMatcher := fmt.Sprintf("id in (%s) or root_id in (%s)",
			formIdParams.Join(","),
			rootFormIdParams.Join(","),
		)
		var formParams []interface{}
		for formId := range formIdParams {
			formParams = append(formParams, formId)
		}
		for rootFormIdParam := range rootFormIdParams {
			formParams = append(formParams, rootFormIdParam)
		}
		var formMatcherArgs []interface{}
		for formId := range formIds {
			formMatcherArgs = append(formMatcherArgs, formId)
		}
		for formId := range formIds {
			formMatcherArgs = append(formMatcherArgs, formId)
		}
		if err := tx.Where(formMatcher, formParams...).Delete(&Form{}).Error; err != nil {
			return meta.NewInternalServerError(err)
		}

		if err := tx.Delete(&Form{}, "folder_id = ?", id).Error; err != nil {
			return meta.NewInternalServerError(err)
		}

		if err := tx.Delete(&Folder{}, "id = ?", id).Error; err != nil {
			return meta.NewInternalServerError(err)
		}

		return nil

	})

	if err != nil {
		return err
	}

	return nil

}

func mapFolderTo(folder *Folder) *types.Folder {
	return &types.Folder{
		ID:         folder.ID,
		Name:       folder.Name,
		DatabaseID: folder.DatabaseID,
		ParentID:   folder.ParentID,
	}
}

func mapFolderFrom(folder *types.Folder) *Folder {
	return &Folder{
		ID:         folder.ID,
		DatabaseID: folder.DatabaseID,
		ParentID:   folder.ParentID,
		Name:       folder.Name,
	}
}
