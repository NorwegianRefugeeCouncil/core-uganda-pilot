package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/utils/sets"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
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
	ctx, db, l, done, err := actionContext(ctx, d.db, "folder", "get", zap.String("folder_id", folderID))
	if err != nil {
		return nil, err
	}
	defer done()

	l.Debug("getting folder")
	var folder Folder
	if err := db.WithContext(ctx).First(&folder, folderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.Error("folder not found", zap.Error(err))
			return nil, meta.NewNotFound(meta.GroupResource{
				Group:    "nrc.no",
				Resource: "folders",
			}, folderID)
		} else {
			l.Error("failed to get folder", zap.Error(err))
			return nil, meta.NewInternalServerError(err)
		}
	}

	l.Debug("successfully got folder")
	return mapFolderTo(&folder), nil
}

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

func (d *folderStore) Create(ctx context.Context, folder *types.Folder) (*types.Folder, error) {
	ctx, db, l, done, err := actionContext(ctx, d.db, "folder", "create")
	if err != nil {
		return nil, err
	}
	defer done()

	folder.ID = uuid.NewV4().String()
	storedFolder := mapFolderFrom(folder)
	storedFolder.CreatedAt = time.Now()

	l.Debug("creating folder")
	if err := db.WithContext(ctx).Create(storedFolder).Error; err != nil {
		l.Error("failed to create folder", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	l.Debug("successfully created folder", zap.String("folder_id", folder.ID))
	return mapFolderTo(storedFolder), nil
}

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

		l.Debug("deleting folder")
		if err := tx.Delete(&Folder{}, "id = ?", id).Error; err != nil {
			l.Error("failed to delete folder", zap.Error(err))
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
