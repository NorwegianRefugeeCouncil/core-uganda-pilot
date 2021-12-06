package store

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/utils/pointers"
	"strings"
	"time"

	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/convert"
	"github.com/nrc-no/core/pkg/utils/sets"
	"github.com/nrc-no/core/pkg/utils/slices"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FormStore interface {
	Get(ctx context.Context, formID string) (*types.FormDefinition, error)
	List(ctx context.Context) (*types.FormDefinitionList, error)
	Create(ctx context.Context, form *types.FormDefinition) (*types.FormDefinition, error)
	Delete(ctx context.Context, id string) error
}

func NewFormStore(db Factory) FormStore {
	return &formStore{
		db: db,
	}
}

type Form struct {
	ID          string `gorm:"index:idx_form_id_database_id,unique"`
	DatabaseID  string `gorm:"index:idx_form_id_database_id,unique"`
	Database    Database
	RootOwnerID string
	RootOwner   *Form
	OwnerID     string
	Owner       *Form
	FolderID    *string
	Folder      *Folder
	Name        string
	Code        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Field struct {
	ID                   string
	DatabaseID           string
	FormID               string
	RootFormID           string
	RootForm             *Form
	Name                 string
	Options              []Option
	Description          string
	Code                 string
	Key                  bool
	Required             bool
	SubFormID            *string
	SubForm              *Form
	ReferencedDatabaseID *string
	ReferencedDatabase   *Database
	ReferencedFormID     *string
	ReferencedForm       *Form
	Type                 types.FieldKind
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type Option struct {
	Value   string `gorm:"primarykey"`
	FieldID string `gorm:"primarykey"`
}

func (f *Form) IsRoot() bool {
	return f.ID == f.RootOwnerID
}

type Forms []*Form

func (f Forms) FormIDs() sets.String {
	result := sets.NewString()
	for _, form := range f {
		result.Insert(form.ID)
	}
	return result
}

func (f Forms) RootForms() Forms {
	var result Forms
	for _, form := range f {
		if form.IsRoot() {
			result = append(result, form)
		}
	}
	return result
}

type Fields []*Field

func (f Fields) OfType(fieldType types.FieldKind) Fields {
	var result Fields
	for _, field := range f {
		if field.Type == fieldType {
			result = append(result, field)
		}
	}
	return result
}

func (f Fields) FormIDs() sets.String {
	result := sets.NewString()
	for _, field := range f {
		result.Insert(field.FormID)
	}
	return result
}

type formStore struct {
	db Factory
}

var _ FormStore = &formStore{}

func (d *formStore) Get(ctx context.Context, formID string) (*types.FormDefinition, error) {
	ctx, db, l, done, err := actionContext(ctx, d.db, "form", "get", zap.String("form_id", formID))
	if err != nil {
		return nil, err
	}
	defer done()

	return d.getFormDefinitionInternal(ctx, db, formID, l)

}

func (d *formStore) getFormDefinitionsInternal(ctx context.Context, db *gorm.DB, formOrSubFormIds sets.String, l *zap.Logger) (*types.FormDefinitionList, error) {

	l.Debug("getting form definitions")
	forms, err := findRelatedFormsInternal(ctx, db, formOrSubFormIds, l)
	if err != nil {
		l.Error("failed to get form definitions", zap.Error(err))
		return nil, err
	}

	if len(forms) == 0 {
		l.Debug("no form definitions found")
		return types.NewFormDefinitionList(), nil
	}

	rootFormIDs := forms.RootForms().FormIDs().List()

	l.Debug("finding all fields form form")
	fields, err := findAllFieldsUnderRootFormIds(ctx, rootFormIDs, db)
	if err != nil {
		l.Error("failed to find fields for form", zap.Error(err))
		return nil, err
	}

	l.Debug("mapping form definitions")
	formDefinitions, err := mapToFormDefinitions(forms, fields)
	if err != nil {
		l.Error("failed to map form definitions", zap.Error(err))
		return nil, err
	}

	l.Debug("successfully listed form definitions")
	return types.NewFormDefinitionList(formDefinitions...), nil

}

func (d *formStore) getFormDefinitionInternal(ctx context.Context, db *gorm.DB, formOrSubFormId string, l *zap.Logger) (*types.FormDefinition, error) {

	l.Debug("getting form definition")

	formDefinitions, err := d.getFormDefinitionsInternal(ctx, db, sets.NewString(formOrSubFormId), l)
	if err != nil {
		l.Error("failed to get form definition", zap.Error(err))
		return nil, err
	}

	if formDefinitions.IsEmpty() {
		l.Debug("form definition not found")
		return nil, meta.NewNotFound(meta.GroupResource{
			Group:    "nrc.no",
			Resource: "forms",
		}, formOrSubFormId)
	}

	l.Debug("found form definition")

	return formDefinitions.GetAtIndex(0), nil
}

func findRelatedFormsInternal(ctx context.Context, db *gorm.DB, formIDs sets.String, l *zap.Logger) (Forms, error) {

	l = l.With(zap.Strings("form_ids", formIDs.List()))

	if formIDs.Len() == 0 {
		return Forms{}, nil
	}

	// Will execute a query with a subquery like
	// SELECT * FROM forms WHERE root_owner_id in (SELECT root_owner_id FROM forms WHERE id IN (id1, id2, ...))

	subQuery := db.Model(&Form{}).
		Select("root_owner_id").
		Where(fmt.Sprintf("id in (%s)", sqlPlaceholders(formIDs)), formIDs.ListIntf()...)

	var forms []*Form

	l.Debug("finding related forms")
	query := db.WithContext(ctx).
		Model(&Form{}).
		Where("root_owner_id in (?)", subQuery).
		Find(&forms)

	if err := query.Error; err != nil {
		l.Error("failed to find related forms", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	l.Debug("successfully found related forms")

	return forms, nil
}

func findAllFieldsUnderRootFormIds(ctx context.Context, rootFormIDs []string, db *gorm.DB) ([]*Field, error) {
	formIdSet := sets.NewString(rootFormIDs...)

	var fields []*Field

	query := db.WithContext(ctx).
		Where(fmt.Sprintf("root_form_id in (%s)", sqlPlaceholders(formIdSet)), formIdSet.ListIntf()...).
		Find(&fields)

	if err := query.Error; err != nil {
		return nil, meta.NewInternalServerError(err)
	}

	return fields, nil
}

func (d *formStore) List(ctx context.Context) (*types.FormDefinitionList, error) {
	ctx, db, l, done, err := actionContext(ctx, d.db, "form", "list")
	if err != nil {
		return nil, err
	}
	defer done()

	var forms []*Form
	var fields []*Field

	db = db.WithContext(ctx)

	l.Debug("finding forms")
	if err := db.Find(&forms).Error; err != nil {
		l.Error("failed to find forms", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	l.Debug("finding fields")
	if err := db.Preload("Options").Find(&fields).Error; err != nil {
		l.Error("failed to find fields", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	l.Debug("mapping forms", zap.Int("field_count", len(fields)), zap.Int("form_count", len(forms)))
	result, err := mapToFormDefinitions(forms, fields)
	if err != nil {
		l.Error("failed to map forms", zap.Error(err))
		return nil, err
	}

	l.Debug("successfully listed forms", zap.Int("count", len(result)))
	return &types.FormDefinitionList{
		Items: result,
	}, nil

}

func (d *formStore) Delete(ctx context.Context, id string) error {

	ctx, db, l, done, err := actionContext(ctx, d.db, "form", "delete")
	if err != nil {
		return err
	}
	defer done()

	l.Debug("getting form definition")
	formDef, err := d.getFormDefinitionInternal(ctx, db, id, l)
	if err != nil {
		l.Error("failed to get form definition", zap.Error(err))
		return err
	}

	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		formIdSet := formDef.GetAllFormsAndSubFormIDs()
		if formIdSet.IsEmpty() {
			return nil
		}
		var formIdsIntf []interface{}

		var params []string
		for _, formId := range formIdSet.List() {
			formIdsIntf = append(formIdsIntf, formId)
			params = append(params, "?")
		}

		l.Debug("deleting database tables")
		formIds := slices.ReversedStrings(formIdSet.List())
		for _, formId := range formIds {
			if err := convert.DeleteTableIfExists(tx, formDef.DatabaseID, formId); err != nil {
				l.Error("failed to delete database table",
					zap.Error(err),
					zap.String("database_id", formDef.DatabaseID),
					zap.String("form_id", formId))
				return err
			}
		}

		l.Debug("deleting fields", zap.Strings("form_ids", formIds))
		fieldsWhereClause := fmt.Sprintf("root_form_id in (%s)", strings.Join(params, ","))
		if err := tx.Where(fieldsWhereClause, formIdsIntf...).Delete(&Field{}).Error; err != nil {
			l.Error("failed to delete fields", zap.Error(err))
			return err
		}

		l.Debug("deleting forms", zap.Strings("form_ids", formIds))
		formsWhereClause := fmt.Sprintf("id in (%s)", strings.Join(params, ","))
		if err := tx.Where(formsWhereClause, formIdsIntf...).Delete(&Form{}).Error; err != nil {
			l.Error("failed to delete fields", zap.Error(err))
			return err
		}

		return nil

	})

	l.Debug("transaction ended")
	return err

}

func (d *formStore) Create(ctx context.Context, form *types.FormDefinition) (*types.FormDefinition, error) {
	ctx, db, l, done, err := actionContext(ctx, d.db, "form", "create")
	if err != nil {
		return nil, err
	}
	defer done()

	newFormIDs(form)

	l.Debug("mapping form definition")
	forms, fields, err := mapToFormFields(form)
	if err != nil {
		l.Error("failed to map form definition", zap.Error(err))
		return nil, err
	}

	referencedFormIDs := sets.NewString()
	for _, field := range fields.OfType(types.FieldKindReference) {
		referencedFormIDs.Insert(*field.ReferencedFormID)
	}

	l.Debug("starting transaction")
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		l.Debug("storing form")
		if err := tx.Create(forms).Error; err != nil {
			l.Error("failed to store form", zap.Error(err))
			return meta.NewInternalServerError(err)
		}

		l.Debug("storing fields")
		if err := tx.Create(fields).Error; err != nil {
			l.Error("failed to store fields", zap.Error(err))
			return meta.NewInternalServerError(err)
		}

		l.Debug("getting referenced forms")
		referencedForms, err := d.getFormDefinitionsInternal(ctx, tx, referencedFormIDs, l)
		if err != nil {
			l.Error("failed to get referenced form")
			return err
		}

		l.Debug("creating form database schema")
		if err := convert.CreateForm(ctx, tx, form, referencedForms); err != nil {
			l.Error("failed to create form database schema")
			return meta.NewInternalServerError(err)
		}

		return nil
	})
	l.Debug("transaction ended")

	if err != nil {
		l.Error("failed to execute transaction", zap.Error(err))
		return nil, err
	}

	l.Debug("successfully created form")
	return d.Get(ctx, form.ID)
}

func newFormIDs(form *types.FormDefinition) {
	form.ID = uuid.NewV4().String()
	newFieldIDs(form.Fields)
}

func newFieldIDs(fields []*types.FieldDefinition) {
	for _, field := range fields {
		field.ID = uuid.NewV4().String()
		if field.FieldType.SubForm != nil {
			newFieldIDs(field.FieldType.SubForm.Fields)
		}
	}
}

func mapToFormFields(fd *types.FormDefinition) (Forms, Fields, error) {
	h, err := NewFormHierarchyFrom(fd)
	if err != nil {
		return nil, nil, err
	}
	flds := h.AllFields()
	frms := h.AllForms()
	return frms, flds, nil
}

func mapToFormDefinitions(forms []*Form, fields []*Field) ([]*types.FormDefinition, error) {

	hierarchies, err := Build(forms, fields)
	if err != nil {
		return nil, err
	}

	var result []*types.FormDefinition
	for _, formHierarchy := range hierarchies {
		fd, err := formHierarchy.ToFormDefinition()
		if err != nil {
			return nil, err
		}
		result = append(result, fd)
	}

	if result == nil {
		result = []*types.FormDefinition{}
	}
	return result, nil
}

type FormHierarchy struct {
	Form       *Form
	Fields     []*Field
	Owner      *FormHierarchy
	OwnerField *Field
	IsRoot     bool
	Children   []*FormHierarchy
}

func NewFormHierarchyFrom(fd *types.FormDefinition) (*FormHierarchy, error) {

	var folderId *string = nil
	if len(fd.FolderID) > 0 {
		folderId = pointers.String(fd.FolderID)
	}

	return newFormHierarchyFrom(
		fd.DatabaseID,
		folderId,
		fd.ID,
		fd.Name,
		fd.Fields,
		nil,
		nil,
		nil)
}

func newFormHierarchyFrom(
	databaseId string,
	folderId *string,
	formId string,
	formName string,
	fields []*types.FieldDefinition,
	owner *FormHierarchy,
	parentField *Field,
	root *FormHierarchy,
) (*FormHierarchy, error) {

	var rootIdForForm = formId
	if root != nil {
		rootIdForForm = root.Form.ID
	}
	var ownerIdForForm = formId
	if owner != nil {
		ownerIdForForm = owner.Form.ID
	}

	hierarchy := &FormHierarchy{
		Form: &Form{
			ID:          formId,
			DatabaseID:  databaseId,
			FolderID:    folderId,
			Name:        formName,
			OwnerID:     ownerIdForForm,
			RootOwnerID: rootIdForForm,
		},
		Owner:      owner,
		OwnerField: parentField,
		IsRoot:     false,
		Children:   []*FormHierarchy{},
	}

	rootIdForField := formId
	if root != nil {
		rootIdForField = root.Form.ID
	}

	rootForField := hierarchy
	if root != nil {
		rootForField = root
	}

	for _, field := range fields {
		ft, err := getFieldType(field)
		if err != nil {
			return nil, err
		}
		var opts []Option
		for _, option := range field.Options {
			opts = append(opts, Option{Value: option})
		}
		f := &Field{
			ID:          field.ID,
			DatabaseID:  databaseId,
			FormID:      formId,
			RootFormID:  rootIdForField,
			Name:        field.Name,
			Required:    field.Required,
			Description: field.Description,
			Code:        field.Code,
			Key:         field.Key,
			Options:     opts,
			Type:        ft,
		}
		hierarchy.Fields = append(hierarchy.Fields, f)

		if field.FieldType.SubForm != nil {

			f.SubFormID = pointers.String(field.ID)

			child, err := newFormHierarchyFrom(
				databaseId,
				folderId,
				field.ID,
				field.Name,
				field.FieldType.SubForm.Fields,
				hierarchy,
				f,
				rootForField,
			)
			if err != nil {
				return nil, err
			}
			hierarchy.Children = append(hierarchy.Children, child)
		} else if field.FieldType.Reference != nil {
			f.ReferencedDatabaseID = &field.FieldType.Reference.DatabaseID
			f.ReferencedFormID = &field.FieldType.Reference.FormID
		}
	}

	return hierarchy, nil

}

func (f *FormHierarchy) ToFormDefinition() (*types.FormDefinition, error) {
	var flds []*types.FieldDefinition
	for _, field := range f.Fields {
		var options []string
		for _, option := range field.Options {
			options = append(options, option.Value)
		}
		fd := &types.FieldDefinition{
			ID:          field.ID,
			Name:        field.Name,
			Code:        field.Code,
			Description: field.Description,
			Key:         field.Key,
			Required:    field.Required,
			FieldType:   types.FieldType{},
			Options:     options,
		}
		switch field.Type {
		case types.FieldKindText:
			fd.FieldType.Text = &types.FieldTypeText{}
		case types.FieldKindMultilineText:
			fd.FieldType.MultilineText = &types.FieldTypeMultilineText{}
		case types.FieldKindDate:
			fd.FieldType.Date = &types.FieldTypeDate{}
		case types.FieldKindMonth:
			fd.FieldType.Month = &types.FieldTypeMonth{}
		case types.FieldKindQuantity:
			fd.FieldType.Quantity = &types.FieldTypeQuantity{}
		case types.FieldKindSingleSelect:
			fd.FieldType.SingleSelect = &types.FieldTypeSingleSelect{}
		case types.FieldKindSubForm:
			child, err := f.GetSubFormForField(field)
			if err != nil {
				return nil, err
			}
			childFd, err := child.ToFormDefinition()
			if err != nil {
				return nil, err
			}
			if childFd.Fields == nil {
				childFd.Fields = []*types.FieldDefinition{}
			}
			fd.FieldType.SubForm = &types.FieldTypeSubForm{
				Fields: childFd.Fields,
			}
		case types.FieldKindReference:
			fd.FieldType.Reference = &types.FieldTypeReference{
				DatabaseID: *field.ReferencedDatabaseID,
				FormID:     *field.ReferencedFormID,
			}
		default:
			return nil, fmt.Errorf("cannot map field type %v", fd.FieldType)
		}
		flds = append(flds, fd)
	}
	var folderId = ""
	if f.Form.FolderID != nil {
		folderId = *f.Form.FolderID
	}
	if len(flds) == 0 {
		flds = []*types.FieldDefinition{}
	}
	formDef := &types.FormDefinition{
		ID:         f.Form.ID,
		DatabaseID: f.Form.DatabaseID,
		FolderID:   folderId,
		Name:       f.Form.Name,
		Fields:     flds,
	}
	return formDef, nil
}

type FormHierarchies []*FormHierarchy

func Build(forms []*Form, fields []*Field) (FormHierarchies, error) {

	var result FormHierarchies
	hierarchyMap := map[string]*FormHierarchy{}
	for _, form := range forms {
		hierarchyMap[form.ID] = &FormHierarchy{
			Form:     form,
			Children: []*FormHierarchy{},
		}
	}

	for _, hierarchy := range hierarchyMap {
		isRoot := hierarchy.Form.RootOwnerID == hierarchy.Form.ID
		hierarchy.IsRoot = isRoot
		if !isRoot {
			parent := hierarchyMap[hierarchy.Form.OwnerID]
			hierarchy.Owner = parent
			parent.Children = append(parent.Children, hierarchy)
		} else {
			result = append(result, hierarchy)
		}
	}

	for _, field := range fields {
		hierarchy := hierarchyMap[field.FormID]
		hierarchy.Fields = append(hierarchy.Fields, field)
	}

	for _, field := range fields {
		if field.Type == types.FieldKindSubForm {
			hierarchyMap[*field.SubFormID].OwnerField = field
		}
	}

	return result, nil
}

func (f *FormHierarchy) GetFormName() string {
	walk := f
	result := ""
	for walk != nil {
		result = f.Form.Name + "_" + result
		walk = walk.Owner
	}
	return result
}

func (f *FormHierarchy) GetSubFormByName(subFormName string) (*FormHierarchy, error) {
	for _, child := range f.Children {
		if child.Form.Name == subFormName {
			return child, nil
		}
	}
	return nil, fmt.Errorf("could not find child with name " + subFormName)
}

func (f *FormHierarchy) GetSubFormForField(field *Field) (*FormHierarchy, error) {
	for _, child := range f.Children {
		if child.OwnerField == field {
			return child, nil
		}
	}
	return nil, fmt.Errorf("could not find child for field " + field.ID)
}

func (f *FormHierarchy) GetSubFormByID(subFormID string) (*FormHierarchy, error) {
	for _, child := range f.Children {
		if child.Form.ID == subFormID {
			return child, nil
		}
	}
	return nil, fmt.Errorf("could not find child with id " + subFormID)
}

func (f *FormHierarchy) AllForms() []*Form {
	result := []*Form{f.Form}
	for _, child := range f.Children {
		result = append(result, child.AllForms()...)
	}
	return result
}

func (f *FormHierarchy) AllFields() []*Field {
	var result []*Field
	result = append(result, f.Fields...)
	for _, child := range f.Children {
		result = append(result, child.AllFields()...)
	}
	return result
}

func getFieldType(field *types.FieldDefinition) (types.FieldKind, error) {
	if field.FieldType.SubForm != nil {
		return types.FieldKindSubForm, nil
	}
	if field.FieldType.Text != nil {
		return types.FieldKindText, nil
	}
	if field.FieldType.Reference != nil {
		return types.FieldKindReference, nil
	}
	if field.FieldType.MultilineText != nil {
		return types.FieldKindMultilineText, nil
	}
	if field.FieldType.Date != nil {
		return types.FieldKindDate, nil
	}
	if field.FieldType.Quantity != nil {
		return types.FieldKindQuantity, nil
	}
	if field.FieldType.Month != nil {
		return types.FieldKindMonth, nil
	}
	if field.FieldType.SingleSelect != nil {
		return types.FieldKindSingleSelect, nil
	}
	return types.FieldKindUnknown, fmt.Errorf("could not determine field type")
}

func sqlPlaceholders(val sets.String) string {
	return strings.Join(val.MapToSlice(func(val string) string {
		return "?"
	}), ",")
}
