package store

import (
	"context"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/pointers"
	"github.com/nrc-no/core/pkg/sets"
	"github.com/nrc-no/core/pkg/sqlconvert"
	"github.com/nrc-no/core/pkg/types"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type FormStore interface {
	Get(ctx context.Context, formID string) (*types.FormDefinition, error)
	List(ctx context.Context) (*types.FormDefinitionList, error)
	Create(ctx context.Context, form *types.FormDefinition) (*types.FormDefinition, error)
}

func NewFormStore(db Factory) FormStore {
	return &formStore{
		db: db,
	}
}

type Form struct {
	ID         string `gorm:"index:idx_form_id_database_id,unique"`
	DatabaseID string `gorm:"index:idx_form_id_database_id,unique"`
	Database   Database
	RootID     string
	Root       *Form
	ParentID   string
	Parent     *Form
	FolderID   *string
	Folder     *Folder
	Name       string
	Code       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (f *Form) IsRoot() bool {
	return f.ID == f.RootID
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

func (f Forms) RootFormIDs() sets.String {
	result := sets.NewString()
	for _, form := range f {
		if form.IsRoot() {
			result.Insert(form.ID)
		}
	}
	return result
}

type Field struct {
	ID                   string
	DatabaseID           string
	FormID               string
	RootFormID           string
	RootForm             *Form
	Name                 string
	Code                 string
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

	db, err := d.db.Get()
	if err != nil {
		return nil, err
	}

	return d.getFormDefinitionInternal(ctx, db, formID)

}

func (d *formStore) getFormDefinitionsInternal(ctx context.Context, db *gorm.DB, formOrSubFormIds sets.String) (*types.FormDefinitionList, error) {

	forms, err := findRelatedFormsInternal(ctx, db, formOrSubFormIds)
	if err != nil {
		return nil, err
	}

	if len(forms) == 0 {
		return types.NewFormDefinitionList(), nil
	}

	rootFormIDs := forms.RootForms().FormIDs().List()

	fields, err := findAllFieldsUnderRootFormIds(ctx, rootFormIDs, db)
	if err != nil {
		return nil, err
	}

	formDefinitions, err := mapToFormDefinitions(forms, fields)
	if err != nil {
		return nil, err
	}

	return types.NewFormDefinitionList(formDefinitions...), nil

}

func (d *formStore) getFormDefinitionInternal(ctx context.Context, db *gorm.DB, formOrSubFormId string) (*types.FormDefinition, error) {

	formDefinitions, err := d.getFormDefinitionsInternal(ctx, db, sets.NewString(formOrSubFormId))
	if err != nil {
		return nil, err
	}

	if formDefinitions.Empty() {
		return nil, meta.NewNotFound(meta.GroupResource{
			Group:    "nrc.no",
			Resource: "forms",
		}, formOrSubFormId)
	}

	return formDefinitions.GetAtIndex(0), nil
}

func findRelatedFormsInternal(ctx context.Context, db *gorm.DB, formIDs sets.String) (Forms, error) {

	if formIDs.Len() == 0 {
		return Forms{}, nil
	}

	// Will execute a query with a subquery like
	// SELECT * FROM forms WHERE root_id in (SELECT root_id FROM forms WHERE id IN (id1, id2, ...))

	subQuery := db.Model(&Form{}).
		Select("root_id").
		Where(fmt.Sprintf("id in (%s)", sqlPlaceholders(formIDs)), formIDs.ListIntf()...)

	var forms []*Form

	query := db.WithContext(ctx).
		Model(&Form{}).
		Where("root_id in (?)", subQuery).
		Find(&forms)

	if err := query.Error; err != nil {
		return nil, meta.NewInternalServerError(err)
	}

	return forms, nil
}

func findFieldForFormIds(ctx context.Context, formIDs []string, db *gorm.DB) (Fields, error) {
	var allParamValues []interface{}
	var allParams []string
	for _, formID := range formIDs {
		allParamValues = append(allParamValues, formID)
		allParams = append(allParams, "?")
	}

	var fields []*Field
	sqlQuery := fmt.Sprintf("form_id in (%s)", strings.Join(allParams, ","))
	if err := db.WithContext(ctx).Where(sqlQuery, allParamValues...).Find(&fields).Error; err != nil {
		return nil, meta.NewInternalServerError(err)
	}
	return fields, nil
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

	db, err := d.db.Get()
	if err != nil {
		return nil, err
	}

	var forms []*Form
	var fields []*Field

	if err := db.Find(&forms).Error; err != nil {
		return nil, meta.NewInternalServerError(err)
	}

	if err := db.Find(&fields).Error; err != nil {
		return nil, meta.NewInternalServerError(err)
	}

	result, err := mapToFormDefinitions(forms, fields)
	if err != nil {
		return nil, err
	}
	if result == nil {
		result = []*types.FormDefinition{}
	}

	return &types.FormDefinitionList{
		Items: result,
	}, nil

}

func (d *formStore) Create(ctx context.Context, form *types.FormDefinition) (*types.FormDefinition, error) {

	db, err := d.db.Get()
	if err != nil {
		return nil, err
	}

	newFormIDs(form)
	newFormCodes(form)

	forms, fields, err := mapToFormFields(form)
	if err != nil {
		return nil, err
	}

	referencedFormIDs := fields.OfType(types.FieldKindReference).FormIDs()

	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(forms).Error; err != nil {
			return meta.NewInternalServerError(err)
		}
		if err := tx.Create(fields).Error; err != nil {
			return meta.NewInternalServerError(err)
		}
		referencedForms, err := d.getFormDefinitionsInternal(ctx, tx, referencedFormIDs)
		if err != nil {
			return err
		}
		if err := sqlconvert.CreateForm(ctx, tx, form, referencedForms); err != nil {
			return meta.NewInternalServerError(err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

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
			field.FieldType.SubForm.ID = uuid.NewV4().String()
			newFieldIDs(field.FieldType.SubForm.Fields)
		}
	}
}

func newFormCodes(form *types.FormDefinition) {
	if len(form.Code) > 0 {
		form.Code = strcase.ToSnake(form.Code)
	} else {
		form.Code = strcase.ToSnake(form.Name)
	}
	newFieldCodes(form.Fields)
}

func newFieldCodes(fields []*types.FieldDefinition) {
	for _, field := range fields {
		field.Code = firstSnake(field.Code, field.Name)
		if field.FieldType.SubForm != nil {
			field.FieldType.SubForm.Code = firstSnake(field.FieldType.SubForm.Code, field.FieldType.SubForm.Name)
			newFieldCodes(field.FieldType.SubForm.Fields)
		}
	}
}

func firstSnake(vals ...string) string {
	for _, val := range vals {
		if len(val) > 0 {
			return strcase.ToSnake(val)
		}
	}
	return ""
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

	return result, nil
}

type FormHierarchy struct {
	Form        *Form
	Fields      []*Field
	Parent      *FormHierarchy
	ParentField *Field
	IsRoot      bool
	Children    []*FormHierarchy
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
		fd.Code,
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
	formCode string,
	fields []*types.FieldDefinition,
	parent *FormHierarchy,
	parentField *Field,
	root *FormHierarchy,
) (*FormHierarchy, error) {

	var rootIdForForm = formId
	if root != nil {
		rootIdForForm = root.Form.ID
	}
	var parentIdForForm = formId
	if parent != nil {
		parentIdForForm = parent.Form.ID
	}

	hierarchy := &FormHierarchy{
		Form: &Form{
			ID:         formId,
			DatabaseID: databaseId,
			FolderID:   folderId,
			Name:       formName,
			Code:       formCode,
			ParentID:   parentIdForForm,
			RootID:     rootIdForForm,
		},
		Parent:      parent,
		ParentField: parentField,
		IsRoot:      false,
		Children:    []*FormHierarchy{},
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
		f := &Field{
			ID:         field.ID,
			DatabaseID: databaseId,
			FormID:     formId,
			RootFormID: rootIdForField,
			Name:       field.Name,
			Code:       field.Code,
			Type:       ft,
		}
		hierarchy.Fields = append(hierarchy.Fields, f)

		if field.FieldType.SubForm != nil {

			f.SubFormID = pointers.String(field.FieldType.SubForm.ID)

			child, err := newFormHierarchyFrom(
				databaseId,
				folderId,
				field.FieldType.SubForm.ID,
				field.FieldType.SubForm.Name,
				field.FieldType.SubForm.Code,
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
		fd := &types.FieldDefinition{
			ID:        field.ID,
			Name:      field.Name,
			Code:      field.Code,
			Required:  false,
			FieldType: types.FieldType{},
		}
		switch field.Type {
		case types.FieldKindText:
			fd.FieldType.Text = &types.FieldTypeText{}
		case types.FieldKindMultilineText:
			fd.FieldType.MultilineText = &types.FieldTypeMultilineText{}
		case types.FieldKindDate:
			fd.FieldType.Date = &types.FieldTypeDate{}
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
				ID:     childFd.ID,
				Fields: childFd.Fields,
				Name:   childFd.Name,
				Code:   childFd.Code,
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
	var folderId string = ""
	if f.Form.FolderID != nil {
		folderId = *f.Form.FolderID
	}
	formDef := &types.FormDefinition{
		ID:         f.Form.ID,
		DatabaseID: f.Form.DatabaseID,
		FolderID:   folderId,
		Name:       f.Form.Name,
		Code:       f.Form.Code,
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
		isRoot := hierarchy.Form.RootID == hierarchy.Form.ID
		hierarchy.IsRoot = isRoot
		if !isRoot {
			parent := hierarchyMap[hierarchy.Form.ParentID]
			hierarchy.Parent = parent
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
			hierarchyMap[*field.SubFormID].ParentField = field
		}
	}

	return result, nil
}

func (h *FormHierarchy) GetFormName() string {
	walk := h
	result := ""
	for walk != nil {
		result = h.Form.Name + "_" + result
		walk = walk.Parent
	}
	return result
}

func (h *FormHierarchy) GetSubFormByName(subFormName string) (*FormHierarchy, error) {
	for _, child := range h.Children {
		if child.Form.Name == subFormName {
			return child, nil
		}
	}
	return nil, fmt.Errorf("could not find child with name " + subFormName)
}

func (h *FormHierarchy) GetSubFormForField(field *Field) (*FormHierarchy, error) {
	for _, child := range h.Children {
		if child.ParentField == field {
			return child, nil
		}
	}
	return nil, fmt.Errorf("could not find child for field " + field.ID)
}

func (h *FormHierarchy) GetSubFormByID(subFormID string) (*FormHierarchy, error) {
	for _, child := range h.Children {
		if child.Form.ID == subFormID {
			return child, nil
		}
	}
	return nil, fmt.Errorf("could not find child with id " + subFormID)
}

func (h *FormHierarchy) AllForms() []*Form {
	result := []*Form{h.Form}
	for _, child := range h.Children {
		result = append(result, child.AllForms()...)
	}
	return result
}

func (h *FormHierarchy) AllFields() []*Field {
	var result []*Field
	result = append(result, h.Fields...)
	for _, child := range h.Children {
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
	return types.FieldKindUnknown, fmt.Errorf("could not determine field type")
}

func sqlPlaceholders(val sets.String) string {
	return strings.Join(val.MapToSlice(func(val string) string {
		return "?"
	}), ",")
}
