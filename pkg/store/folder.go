package store

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/utils/sets"
	"github.com/nrc-no/core/pkg/validation"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

// FolderStore is the store for types.Folder
type FolderStore interface {
	// Get a folder
	Get(ctx context.Context, folderID string) (*types.Folder, error)
	// List folders
	List(ctx context.Context) (*types.FolderList, error)
	// Create a folder
	Create(ctx context.Context, folder *types.Folder) (*types.Folder, error)
	// Delete a folder
	Delete(ctx context.Context, id string) error
}

// Folder is the store model for types.Folder
type Folder struct {
	ID         string
	DatabaseID string
	Database   Database
	ParentID   *string
	Parent     *Folder `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewFolderStore creates a FolderStore
func NewFolderStore(db Factory) FolderStore {
	return &folderStore{
		db: db,
	}
}

// folderStore is the implementation of FolderStore
type folderStore struct {
	db Factory
}

// Make sure that folderStore implements FolderStore
var _ FolderStore = &folderStore{}

// Get implements FolderStore.Get
func (d *folderStore) Get(ctx context.Context, folderID string) (*types.Folder, error) {
	ctx, db, l, done, err := actionContext(ctx, d.db, "folder", "get", zap.String("folder_id", folderID))
	if err != nil {
		return nil, err
	}
	defer done()

	var folder Folder
	if err := db.First(&folder, "id = ?", folderID).Error; err != nil {
		if IsNotFoundErr(err) {
			l.Error("folder not found", zap.Error(err))
			return nil, meta.NewNotFound(types.FolderGR, folderID)
		}
		l.Error("failed to get folder", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	return mapFolderTo(&folder), nil
}

// List implements FolderStore.List
func (d *folderStore) List(ctx context.Context) (*types.FolderList, error) {
	ctx, db, l, done, err := actionContext(ctx, d.db, "folder", "list")
	if err != nil {
		return nil, err
	}
	defer done()

	l.Debug("listing folders")
	var folders []Folder
	if err := db.WithContext(ctx).Find(&folders).Error; err != nil {
		l.Error("failed to list folders", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	l.Debug("mapping folders")
	var result = make([]*types.Folder, len(folders))
	for i, folder := range folders {
		result[i] = mapFolderTo(&folder)
	}
	if result == nil {
		result = []*types.Folder{}
	}

	l.Debug("successfully listed folders", zap.Int("count", len(result)))
	return &types.FolderList{
		Items: result,
	}, nil
}

// Create implements FolderStore.Create
func (d *folderStore) Create(ctx context.Context, folder *types.Folder) (*types.Folder, error) {
	ctx, db, l, done, err := actionContext(ctx, d.db, "folder", "create", zap.String("folder_id", folder.ID))
	if err != nil {
		return nil, err
	}
	defer done()

	storedFolder := mapFolderFrom(folder)
	storedFolder.CreatedAt = time.Now()

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&Database{}, "id = ?", folder.DatabaseID).Error; err != nil {
			if IsNotFoundErr(err) {
				err := meta.NewInvalid(types.FolderGR, "", validation.ErrorList{validation.NotFound(validation.NewPath("databaseId"), folder.DatabaseID)})
				l.Error("database not found", zap.Error(err))
				return err
			}
			l.Error("failed to lookup database", zap.Error(err))
			return meta.NewInternalServerError(err)
		}
		if len(folder.ParentID) != 0 {
			if err := tx.First(&Folder{}, "id = ?", folder.ParentID).Error; err != nil {
				if IsNotFoundErr(err) {
					err := meta.NewInvalid(types.FolderGR, "", validation.ErrorList{validation.NotFound(validation.NewPath("parentId"), folder.ParentID)})
					l.Error("parent folder not found", zap.Error(err))
					return err
				}
				l.Error("failed to lookup parent folder", zap.Error(err))
				return meta.NewInternalServerError(err)
			}
		}
		if err := tx.Create(storedFolder).Error; err != nil {
			l.Error("failed to create folder", zap.Error(err))
			if IsUniqueConstraintErr(err) {
				err := meta.NewAlreadyExists(types.FolderGR, folder.ID)
				l.Error("folder already exists", zap.Error(err))
				return err
			}
			return meta.NewInternalServerError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return mapFolderTo(storedFolder), nil
}

// Delete implements FolderStore.Delete
func (d *folderStore) Delete(ctx context.Context, id string) error {
	ctx, db, l, done, err := actionContext(ctx, d.db, "folder", "delete", zap.String("folder_id", id))
	if err != nil {
		return err
	}
	defer done()

	l.Debug("starting transaction")
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		l.Debug("finding root forms in folder")
		var folderForms []*Form
		if err := tx.Find(&folderForms, "folder_id = ?", id).Error; err != nil {
			l.Error("failed to find forms in folder", zap.Error(err))
			return meta.NewInternalServerError(err)
		}

		l.Debug("found root forms in folder", zap.Int("count", len(folderForms)))

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

		if len(formIds) != 0 {

			l.Debug("finding all fields for forms in folder")
			fieldFormMatcher := fmt.Sprintf("root_form_id in (%s)", fieldParams.Join(","))
			qry := tx.Where(fieldFormMatcher, formIds.ListIntf()...).Delete(&Field{})
			if err := qry.Error; err != nil {
				l.Error("failed to find all fields for forms in folder", zap.Error(err))
				return meta.NewInternalServerError(err)
			}
			l.Debug("deleted fields", zap.Int64("deleted_count", qry.RowsAffected))

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

			l.Debug("deleting all forms in folder (by form_id or root_id)",
				zap.Strings("form_ids", formIdParams.List()),
				zap.Strings("root_form_ids", rootFormIdParams.List()))

			qry = tx.Where(formMatcher, formParams...).Delete(&Form{})
			if err := qry.Error; err != nil {
				l.Error("failed to delete all forms and subforms in folder", zap.Error(err))
				return meta.NewInternalServerError(err)
			}
			l.Debug("deleted forms", zap.Int64("deleted_count", qry.RowsAffected))

			l.Debug("deleting all forms in folder (by folder id)")
			qry = tx.Delete(&Form{}, "folder_id = ?", id)
			if err := qry.Error; err != nil {
				return meta.NewInternalServerError(err)
			}
			l.Debug("deleted forms", zap.Int64("deleted_count", qry.RowsAffected))

		}

		l.Debug("deleting folder")
		deleteQry := tx.Delete(&Folder{}, "id = ?", id)
		if err := deleteQry.Error; err != nil {
			l.Error("failed to delete folder", zap.Error(err))
			return meta.NewInternalServerError(err)
		}
		if deleteQry.RowsAffected == 0 {
			return meta.NewNotFound(types.FolderGR, id)
		}

		return nil

	})

	if err != nil {
		return err
	}

	return nil

}

// mapFolderTo maps Folder to a types.Folder
func mapFolderTo(folder *Folder) *types.Folder {
	var parentId = ""
	if folder.ParentID != nil {
		parentId = *folder.ParentID
	}
	return &types.Folder{
		ID:         folder.ID,
		Name:       folder.Name,
		DatabaseID: folder.DatabaseID,
		ParentID:   parentId,
	}
}

// mapFolderFrom maps a types.Folder to a Folder for storage
func mapFolderFrom(folder *types.Folder) *Folder {
	var parentId *string = nil
	if len(folder.ParentID) != 0 {
		parentId = &folder.ParentID
	}
	return &Folder{
		ID:         folder.ID,
		DatabaseID: folder.DatabaseID,
		ParentID:   parentId,
		Name:       folder.Name,
	}
}
