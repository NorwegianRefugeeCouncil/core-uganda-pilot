package store

import (
	"context"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/bla/sqlconvert"
	"github.com/nrc-no/core/pkg/bla/types"
	"github.com/nrc-no/core/pkg/pointers"
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
	RootID     *string
	Root       *Form
	ParentID   *string
	Parent     *Form
	FolderID   *string
	Folder     *Folder
	Name       string
	Code       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type FieldType string

const (
	FieldTypeUnknown   = "unknown"
	FieldTypeText      = "text"
	FieldTypeSubForm   = "subform"
	FieldTypeReference = "reference"
)

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
	Type                 FieldType
	CreatedAt            time.Time
	UpdatedAt            time.Time
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

	var forms []*Form
	var fields []*Field

	var requestedForm *Form
	if err := db.WithContext(ctx).Find(&requestedForm, "id = ?", formID).Error; err != nil {
		return nil, meta.NewInternalServerError(err)
	}

	var rootId = requestedForm.ID
	if requestedForm.RootID != nil {
		rootId = *requestedForm.RootID
	}

	err = db.WithContext(ctx).Find(&forms, "id = ? or root_id = ?", rootId, rootId).Error
	if err != nil {
		return nil, meta.NewInternalServerError(err)
	}

	var allParamValues []interface{}
	var allParams []string
	for _, form := range forms {
		allParamValues = append(allParamValues, form.ID)
		allParams = append(allParams, "?")
	}

	sqlQuery := fmt.Sprintf("form_id in (%s)", strings.Join(allParams, ","))
	if err := db.WithContext(ctx).Where(sqlQuery, allParamValues...).Find(&fields).Error; err != nil {
		return nil, meta.NewInternalServerError(err)
	}

	if len(forms) == 0 {
		return nil, meta.NewNotFound(meta.GroupResource{
			Group:    "nrc.no",
			Resource: "forms",
		}, formID)
	}

	fds, err := mapToFormDefinitions(forms, fields)
	if err != nil {
		return nil, err
	}

	if len(fds) != 1 {
		return nil, meta.NewInternalServerError(fmt.Errorf("unexpected number of form definitions"))
	}

	return fds[0], nil

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

	frms, flds, err := mapToFormFields(form)
	if err != nil {
		return nil, err
	}

	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(frms).Error; err != nil {
			return meta.NewInternalServerError(err)
		}
		if err := tx.Create(flds).Error; err != nil {
			return meta.NewInternalServerError(err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, meta.NewInternalServerError(err)
	}

	if err := sqlconvert.CreateForm(ctx, sqlDB, form); err != nil {
		return nil, meta.NewInternalServerError(err)
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
		field.Code = snake(field.Code, field.Name)
		if field.FieldType.SubForm != nil {
			field.FieldType.SubForm.Code = snake(field.FieldType.SubForm.Code, field.FieldType.SubForm.Name)
			newFieldCodes(field.FieldType.SubForm.Fields)
		}
	}
}

func snake(vals ...string) string {
	for _, val := range vals {
		if len(val) > 0 {
			return strcase.ToSnake(val)
		}
	}
	return ""
}

func mapToFormFields(fd *types.FormDefinition) ([]*Form, []*Field, error) {
	h, err := NewFormHierarchyFrom(fd)
	if err != nil {
		return nil, nil, err
	}
	flds := h.AllFields()
	frms := h.AllForms()
	return frms, flds, nil
}

func mapToFormDefinitions(forms []*Form, fields []*Field) ([]*types.FormDefinition, error) {

	hierarchies, err := buildHierarchies(forms, fields)
	if err != nil {
		return nil, err
	}

	var result []*types.FormDefinition
	for _, formHierarchy := range hierarchies {
		fd, err := formHierarchy.convertToFormDef()
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

	var rootIdForForm *string
	if root != nil {
		rootIdForForm = pointers.String(root.Form.ID)
	}
	var parentIdForForm *string
	if parent != nil {
		parentIdForForm = pointers.String(parent.Form.ID)
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

func (f *FormHierarchy) convertToFormDef() (*types.FormDefinition, error) {
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
		case FieldTypeText:
			fd.FieldType.Text = &types.FieldTypeText{}
		case FieldTypeSubForm:
			child, err := f.GetSubFormForField(field)
			if err != nil {
				return nil, err
			}
			childFd, err := child.convertToFormDef()
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
		case FieldTypeReference:
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

func buildHierarchies(forms []*Form, fields []*Field) (FormHierarchies, error) {

	var result FormHierarchies
	hierarchyMap := map[string]*FormHierarchy{}
	for _, form := range forms {
		hierarchyMap[form.ID] = &FormHierarchy{
			Form:     form,
			Children: []*FormHierarchy{},
		}
	}

	for _, hierarchy := range hierarchyMap {
		isRoot := hierarchy.Form.RootID == nil
		hierarchy.IsRoot = isRoot
		if !isRoot {
			parent := hierarchyMap[*hierarchy.Form.ParentID]
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
		if field.Type == FieldTypeSubForm {
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

func getFieldType(field *types.FieldDefinition) (FieldType, error) {
	if field.FieldType.SubForm != nil {
		return FieldTypeSubForm, nil
	}
	if field.FieldType.Text != nil {
		return FieldTypeText, nil
	}
	if field.FieldType.Reference != nil {
		return FieldTypeReference, nil
	}
	return FieldTypeUnknown, fmt.Errorf("could not determine field type")
}
