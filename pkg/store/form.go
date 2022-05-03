package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/sqlmanager"
	"github.com/nrc-no/core/pkg/utils/pointers"
	"golang.org/x/sync/errgroup"

	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/convert"
	"github.com/nrc-no/core/pkg/utils/sets"
	"github.com/nrc-no/core/pkg/utils/slices"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// FormStore is the store implementation for types.FormDefinition
//
// The responsibility of the FormStore is twofold.
//
// When creating new forms, the FormStore must be able to "flatten" the
// structure of the FormDefinition into entities that are storable in SQL Tables.
// It must also be able to re-hydrate the entities  into the tree-like structure of a FormDefinition.
// The "flatten-" and "hydrate-" method do that, respectively.
//
// The next responsibility of the FormStore is the create the underlying SQL Tables
// that will store the records for the created FormDefinitions.
type FormStore interface {
	// Get a single Form Definition
	Get(ctx context.Context, formID string) (*types.FormDefinition, error)
	// List FormDefinitions
	List(ctx context.Context) (*types.FormDefinitionList, error)
	// Create a FormDefinition
	Create(ctx context.Context, form *types.FormDefinition) (*types.FormDefinition, error)
	// Delete a FormDefinition
	Delete(ctx context.Context, id string) error
}

// NewFormStore returns a new FormStore
func NewFormStore(db Factory) FormStore {
	return &formStore{
		db: db,
	}
}

// Form represents the data structure used to store a types.FormDefinition
type Form struct {

	// ID stores the types.FormDefinition ID
	ID string `gorm:"index:idx_form_id_database_id,unique"`

	// DatabaseID stores the types.FormDefinition DatabaseID
	DatabaseID string `gorm:"index:idx_form_id_database_id,unique"`

	// Database is used so that gorm creates a (not nullable) Foreign Key on the Database table
	Database Database

	// RootOwnerID is present because when we want to fetch the root form and all SubForms
	// for a given FormDefinition, we can simply query WHERE root_owner_id = <RootFormID>.
	// This will retrieve all the forms, subforms and nested subforms that are part
	// of the same FormDefinition
	RootOwnerID string

	// RootOwner is the gorm hack to create a (not nullable) Foreign Key on the Form table
	RootOwner *Form

	// OwnerID represents a SubForm owner OR, when the form is the root form (FormDefinition)
	// then the OwnerID == ID. It's nice like this because we can query all root FormDefinitions
	// with a query like "SELECT * FROM forms where id == owner_id"
	OwnerID string

	// Owner is the same gorm hack to create a (not nullable) foreign key on the Form table
	Owner *Form

	// FolderID stores the types.FormDefinition FolderID
	FolderID *string

	// Folder is the gorm hack to create a (nullable) foreign key on the Folder table
	Folder *Folder

	// Name stores the types.FormDefinition Name
	Name string

	// Type stores the types.FormDefinition Type
	Type string

	// CreatedAt represents when this form was created
	CreatedAt time.Time

	// UpdatedAt represents when this form was updated
	UpdatedAt time.Time
}

type Forms []*Form

// Field represents the data structure used to store a types.FieldDefinition
type Field struct {

	// ID stores the types.FieldDefinition ID
	ID string

	// DatabaseID stores the types.FieldDefinition DatabaseID
	DatabaseID string

	// FormID stores the types.FieldDefinition FormID
	FormID string

	// RootFormID stores the root form definition ID for this field.
	// See comment on Form.RootOwnerID
	RootFormID string

	// RootForm is the gorm hack to create a foreign key
	RootForm *Form

	// Name stores the types.FieldDefinition Name
	Name string

	// Description stores the types.FieldDefinition Description
	Description string

	// Code stores the types.FieldDefinition Code
	Code string

	// Key stores the types.FieldDefinition Key
	Key bool

	// Required stores the types.FieldDefinition Required
	Required bool

	// SubFormID stores the SubForm ID. This should be equal to the Field.ID.
	// Merely keeping a duplicate value so we can have a foreign key
	SubFormID *string

	// SubForm creates a gorm ForeignKey
	SubForm *Form

	// ReferencedDatabaseID stores the referenced DatabaseID when the field is a Reference field
	ReferencedDatabaseID *string

	// ReferencedDatabase allows gorm to create a foreign key
	ReferencedDatabase *Database

	// ReferencedFormID stores the referenced FormID when the field is a Reference field
	ReferencedFormID *string

	// ReferencedForm allows gorm to create a foreign key
	ReferencedForm *Form

	// Type stores the field type
	Type types.FieldKind

	// CreatedAt stores the time the Field was created
	CreatedAt time.Time

	// UpdatedAt stores the time the field was updated
	UpdatedAt time.Time
}

type Fields []*Field

// Option represents the data structure used to store options for a types.FieldTypeMultiSelect or types.FieldTypeSingleSelect
type Option struct {

	// ID stores the types.SelectOption ID
	ID string `gorm:"primarykey"`

	// DatabaseID stores the Form's DatabaseID
	DatabaseID string

	// RootFormID stores the root FormID. See Form.RootOwnerID
	RootFormID string

	// FormID stores the FormID of the single/multi select field
	FormID string

	// FieldID stores the single/multi select Field.ID
	FieldID string

	// Name stores the types.SelectOption Name
	Name string
}

type Options []*Option

// formStore is the implementation of FormStore
type formStore struct {
	db Factory
}

// Make sure that formStore implements FormStore
var _ FormStore = &formStore{}

// Get implements FormStore.Get
func (d *formStore) Get(ctx context.Context, formID string) (*types.FormDefinition, error) {
	ctx, db, _, done, err := actionContext(ctx, d.db, "form", "get", zap.String("form_id", formID))
	if err != nil {
		return nil, err
	}
	defer done()
	return d.getFormDefinitionInternal(ctx, db, formID)
}

// List implements FormStore.List
func (d *formStore) List(ctx context.Context) (*types.FormDefinitionList, error) {
	ctx, db, _, done, err := actionContext(ctx, d.db, "form", "list")
	if err != nil {
		return nil, err
	}
	defer done()

	db = db.WithContext(ctx)

	// Finding the flattened structure of the forms
	flatForms, err := findFlattenedForms(ctx, db, nil)
	if err != nil {
		return nil, err
	}

	// Hydrating the flat structure into a tree structure
	formDefs, err := flatForms.hydrateForms()
	if err != nil {
		return nil, err
	}

	return &types.FormDefinitionList{
		Items: formDefs,
	}, nil

}

// Delete implements FormStore.Delete
func (d *formStore) Delete(ctx context.Context, id string) error {

	ctx, db, l, done, err := actionContext(ctx, d.db, "form", "delete")
	if err != nil {
		return err
	}
	defer done()

	formDef, err := d.getFormDefinitionInternal(ctx, db, id)
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

// Create implements FormStore.Create
func (d *formStore) Create(ctx context.Context, form *types.FormDefinition) (*types.FormDefinition, error) {
	ctx, db, l, done, err := actionContext(ctx, d.db, "form", "create")
	if err != nil {
		return nil, err
	}
	defer done()

	// Flattens the types.FormDefinition into a flat storage structure
	flatForm, err := flattenForm(form)
	if err != nil {
		l.Error("failed to flatten form: %v", zap.Error(err))
		return nil, err
	}

	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// Storing the forms
		if err := tx.Create(flatForm.Forms).Error; err != nil {
			l.Error("failed to store form", zap.Error(err))
			return meta.NewInternalServerError(err)
		}

		// Storing the fields
		if err := tx.Create(flatForm.Fields).Error; err != nil {
			l.Error("failed to store fields", zap.Error(err))
			return meta.NewInternalServerError(err)
		}

		// Storing the options
		if len(flatForm.Options) != 0 {
			if err := tx.Create(flatForm.Options).Error; err != nil {
				l.Error("failed to store fields", zap.Error(err))
				return meta.NewInternalServerError(err)
			}
		}

		// Creating the actual SQL Tables to contain records for the form
		m, err := sqlmanager.New().PutForms(&types.FormDefinitionList{
			Items: []*types.FormDefinition{form},
		})
		if err != nil {
			l.Error("failed to create form migrations", zap.Error(err))
			return err
		}

		// Executing the SQL statements
		for _, ddl := range m.GetStatements() {
			if err := tx.Exec(ddl.Query, ddl.Args...).Error; err != nil {
				l.Error("failed to execute form migration statement", zap.String("statement", ddl.Query), zap.Error(err))
				return err
			}
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

// getFormDefinitionsInternal is an internal method to retrieve multiple form definitions by form or subForm id
func (d *formStore) getFormDefinitionsInternal(
	ctx context.Context,
	db *gorm.DB,
	formOrSubFormIds sets.String,
) (*types.FormDefinitionList, error) {

	// Do not execute any query if the formOrSubFormIds is empty
	if formOrSubFormIds.Len() == 0 {
		return &types.FormDefinitionList{
			Items: []*types.FormDefinition{},
		}, nil
	}

	l := logging.NewLogger(ctx)

	// this is the query to retrieve the root form ids
	subQuery := db.Model(&Form{}).
		Select("root_owner_id").
		Where(fmt.Sprintf("id in (%s)",
			sqlPlaceholders(formOrSubFormIds)),
			formOrSubFormIds.ListIntf()...)

	// find the flattened structure for the given form definition query
	flattenedForms, err := findFlattenedForms(ctx, db, subQuery)
	if err != nil {
		l.Error("failed to get form definitions", zap.Error(err))
		return nil, err
	}

	// return empty result if we don't have any forms returned
	if len(flattenedForms.Forms) == 0 {
		l.Debug("no form definitions found")
		return types.NewFormDefinitionList(), nil
	}

	// hydrate back the form definitions
	formDefinitions, err := flattenedForms.hydrateForms()
	if err != nil {
		l.Error("failed to unflatten forms", zap.Error(err))
		return nil, err
	}

	return &types.FormDefinitionList{
		Items: formDefinitions,
	}, nil

}

// getFormDefinitionInternal returns a single form definition given the form or sub form id
func (d *formStore) getFormDefinitionInternal(ctx context.Context, db *gorm.DB, formOrSubFormId string) (*types.FormDefinition, error) {

	l := logging.NewLogger(ctx)
	l.Debug("getting form definition")

	formDefinitions, err := d.getFormDefinitionsInternal(ctx, db, sets.NewString(formOrSubFormId))
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

// findFlattenedForms will retrieve the flattened FormDefinition structure
// the passed formIdsCondition is optional, all is used to restrict which FormDefinitions
// are we looking for
func findFlattenedForms(
	ctx context.Context,
	db *gorm.DB,
	formIdsCondition *gorm.DB,
) (FlatForms, error) {

	l := logging.NewLogger(ctx)

	var forms []*Form
	var fields []*Field
	var options []*Option

	// We are querying for forms, fields and options in parallel
	g, ctx := errgroup.WithContext(ctx)

	// Retrieve the Forms
	g.Go(func() error {
		qry := db.WithContext(ctx).Model(&Form{})
		if formIdsCondition != nil {
			qry = qry.Where("root_owner_id in (?)", formIdsCondition)
		}
		if err := qry.Find(&forms).Error; err != nil {
			l.Error("failed to find forms", zap.Error(err))
			return meta.NewInternalServerError(err)
		}
		return nil
	})
	// Retrieve the Fields
	g.Go(func() error {
		qry := db.WithContext(ctx).Model(&Field{})
		if formIdsCondition != nil {
			qry = qry.Where("root_form_id in (?)", formIdsCondition)
		}
		if err := qry.Find(&fields).Error; err != nil {
			l.Error("failed to find fields", zap.Error(err))
			return meta.NewInternalServerError(err)
		}
		return nil
	})
	// Retrieve the Options
	g.Go(func() error {
		qry := db.WithContext(ctx).Model(&Option{})
		if formIdsCondition != nil {
			qry = qry.Where("root_form_id in (?)", formIdsCondition)
		}
		if err := qry.Find(&options).Error; err != nil {
			l.Error("failed to find options", zap.Error(err))
			return meta.NewInternalServerError(err)
		}
		return nil
	})

	// Waiting for the parallel execution to finish
	if err := g.Wait(); err != nil {
		return FlatForms{}, err
	}

	return FlatForms{
		Fields:  fields,
		Forms:   forms,
		Options: options,
	}, nil
}

func sqlPlaceholders(val sets.String) string {
	return strings.Join(val.MapToSlice(func(val string) string {
		return "?"
	}), ",")
}

type FlatForms struct {
	Fields  Fields
	Forms   Forms
	Options Options
}

// Add concatenates two FlatForms together
func (f FlatForms) Add(flatForm FlatForms) FlatForms {
	f.Fields = append(f.Fields, flatForm.Fields...)
	f.Forms = append(f.Forms, flatForm.Forms...)
	f.Options = append(f.Options, flatForm.Options...)
	return f
}

// hydrateForms will rehydrate a flattened types.FormDefinition
func (f FlatForms) hydrateForms() ([]*types.FormDefinition, error) {
	forms := make([]*types.FormDefinition, 0)
	for _, form := range f.Forms {
		if form.ID != form.RootOwnerID {
			continue
		}
		formDef, err := f.hydrateForm(form)
		if err != nil {
			return nil, err
		}
		forms = append(forms, formDef)
	}
	return forms, nil
}

// hydrateForm rehydrates a single types.FormDefinition
func (f FlatForms) hydrateForm(form *Form) (*types.FormDefinition, error) {
	folderId := ""
	if form.FolderID != nil {
		folderId = *form.FolderID
	}
	result := &types.FormDefinition{
		ID:         form.ID,
		DatabaseID: form.DatabaseID,
		FolderID:   folderId,
		Name:       form.Name,
		Type:       types.FormType(form.Type),
	}
	fields, err := f.hydrateFormFields(form.ID)
	if err != nil {
		return nil, err
	}
	result.Fields = fields
	return result, nil
}

// hydrateFormFields rehydrates the fields for the given formId
func (f FlatForms) hydrateFormFields(formId string) ([]*types.FieldDefinition, error) {
	var result []*types.FieldDefinition
	for _, field := range f.Fields {
		if field.FormID != formId {
			continue
		}
		fieldDef, err := f.hydrateField(field)
		if err != nil {
			return nil, err
		}
		result = append(result, fieldDef)
	}
	if result == nil {
		result = []*types.FieldDefinition{}
	}
	return result, nil
}

func hydrateFieldDefaults(field *Field) *types.FieldDefinition {
	return &types.FieldDefinition{
		ID:          field.ID,
		Code:        field.Code,
		Name:        field.Name,
		Description: field.Description,
		Key:         field.Key,
		Required:    field.Required,
	}
}

// hydrateField rehydrates a single types.FieldDefinition
func (f FlatForms) hydrateField(field *Field) (*types.FieldDefinition, error) {
	switch field.Type {
	case types.FieldKindText:
		return f.hydrateTextField(field)
	case types.FieldKindSubForm:
		return f.hydrateSubFormField(field)
	case types.FieldKindReference:
		return f.hydrateReferenceField(field)
	case types.FieldKindMultilineText:
		return f.hydrateMultilineTextField(field)
	case types.FieldKindDate:
		return f.hydrateDateField(field)
	case types.FieldKindQuantity:
		return f.hydrateQuantityField(field)
	case types.FieldKindMonth:
		return f.hydrateMonthField(field)
	case types.FieldKindWeek:
		return f.hydrateWeekField(field)
	case types.FieldKindSingleSelect:
		return f.hydrateSingleSelectField(field)
	case types.FieldKindMultiSelect:
		return f.hydrateMultiSelectField(field)
	case types.FieldKindCheckbox:
		return f.hydrateCheckboxField(field)
	}
	return nil, errors.New("unprocessable field type")
}

// hydrateSubFormField hydrates a types.FieldTypeSubForm
func (f FlatForms) hydrateSubFormField(field *Field) (*types.FieldDefinition, error) {
	fieldDef := hydrateFieldDefaults(field)
	fieldDef.FieldType = types.FieldType{
		SubForm: &types.FieldTypeSubForm{},
	}
	subFormFields, err := f.hydrateFormFields(field.ID)
	if err != nil {
		return nil, err
	}
	fieldDef.FieldType.SubForm.Fields = subFormFields
	return fieldDef, nil
}

// hydrateTextField hydrates a types.FieldTypeText
func (f FlatForms) hydrateTextField(field *Field) (*types.FieldDefinition, error) {
	fieldDef := hydrateFieldDefaults(field)
	fieldDef.FieldType = types.FieldType{
		Text: &types.FieldTypeText{},
	}
	return fieldDef, nil
}

// hydrateMultilineTextField hydrates a types.FieldTypeMultilineText
func (f FlatForms) hydrateMultilineTextField(field *Field) (*types.FieldDefinition, error) {
	fieldDef := hydrateFieldDefaults(field)
	fieldDef.FieldType = types.FieldType{
		MultilineText: &types.FieldTypeMultilineText{},
	}
	return fieldDef, nil
}

// hydrateDateField hydrates a types.FieldTypeDate
func (f FlatForms) hydrateDateField(field *Field) (*types.FieldDefinition, error) {
	fieldDef := hydrateFieldDefaults(field)
	fieldDef.FieldType = types.FieldType{
		Date: &types.FieldTypeDate{},
	}
	return fieldDef, nil
}

// hydrateMonthField hydrates a types.FieldTypeMonth
func (f FlatForms) hydrateMonthField(field *Field) (*types.FieldDefinition, error) {
	fieldDef := hydrateFieldDefaults(field)
	fieldDef.FieldType = types.FieldType{
		Month: &types.FieldTypeMonth{},
	}
	return fieldDef, nil
}

// hydrateWeekField hydrates a types.FieldTypeWeek
func (f FlatForms) hydrateWeekField(field *Field) (*types.FieldDefinition, error) {
	fieldDef := hydrateFieldDefaults(field)
	fieldDef.FieldType = types.FieldType{
		Week: &types.FieldTypeWeek{},
	}
	return fieldDef, nil
}

// hydrateSingleSelectField hydrates a types.FieldTypeSingleSelect
func (f FlatForms) hydrateSingleSelectField(field *Field) (*types.FieldDefinition, error) {
	fieldDef := hydrateFieldDefaults(field)
	fieldDef.FieldType = types.FieldType{
		SingleSelect: &types.FieldTypeSingleSelect{
			Options: f.hydrateSelectOptions(field.ID),
		},
	}
	return fieldDef, nil
}

// hydrateMultiSelectField hydrates a types.FieldTypeMultiSelect
func (f FlatForms) hydrateMultiSelectField(field *Field) (*types.FieldDefinition, error) {
	fieldDef := hydrateFieldDefaults(field)
	fieldDef.FieldType = types.FieldType{
		MultiSelect: &types.FieldTypeMultiSelect{
			Options: f.hydrateSelectOptions(field.ID),
		},
	}
	return fieldDef, nil
}

// hydrateSelectOptions hydrates a list of single/multi select options
func (f FlatForms) hydrateSelectOptions(fieldId string) []*types.SelectOption {
	var result []*types.SelectOption
	for _, option := range f.Options {
		if option.FieldID != fieldId {
			continue
		}
		result = append(result, &types.SelectOption{
			ID:   option.ID,
			Name: option.Name,
		})
	}
	return result
}

// hydrateQuantityField hydrates a types.FieldTypeQuantity
func (f FlatForms) hydrateQuantityField(field *Field) (*types.FieldDefinition, error) {
	fieldDef := hydrateFieldDefaults(field)
	fieldDef.FieldType = types.FieldType{
		Quantity: &types.FieldTypeQuantity{},
	}
	return fieldDef, nil
}

// hydrateReferenceField hydrates a types.FieldTypeReference
func (f FlatForms) hydrateReferenceField(field *Field) (*types.FieldDefinition, error) {
	fieldDef := hydrateFieldDefaults(field)
	fieldDef.FieldType = types.FieldType{
		Reference: &types.FieldTypeReference{
			DatabaseID: *field.ReferencedDatabaseID,
			FormID:     *field.ReferencedFormID,
		},
	}
	return fieldDef, nil
}

// hydrateCheckboxField hydrates a types.FieldTypeCheckbox
func (f FlatForms) hydrateCheckboxField(field *Field) (*types.FieldDefinition, error) {
	fieldDef := hydrateFieldDefaults(field)
	fieldDef.FieldType = types.FieldType{
		Checkbox: &types.FieldTypeCheckbox{},
	}
	return fieldDef, nil
}

// flattenForm will flatten a types.FormDefinition into a FlatForms
func flattenForm(form *types.FormDefinition) (FlatForms, error) {
	var folderId *string
	if len(form.FolderID) != 0 {
		folderId = pointers.String(form.FolderID)
	}
	var flattenedForm = &Form{
		ID:          form.ID,
		DatabaseID:  form.DatabaseID,
		RootOwnerID: form.ID,
		OwnerID:     form.ID,
		FolderID:    folderId,
		Name:        form.Name,
		Type:        string(form.Type),
	}
	result := FlatForms{
		Forms: []*Form{flattenedForm},
	}
	for _, field := range form.GetFields() {
		flatField, err := flattenField(form, form, field)
		if err != nil {
			return FlatForms{}, err
		}
		result = result.Add(flatField)
	}
	return result, nil
}

// flattenSubForm will flatten a types.FieldTypeSubForm into a FlatForms
func flattenSubForm(rootForm, owner, subForm types.FormInterface, field *types.FieldDefinition) (FlatForms, error) {
	var flattenedForm = &Form{
		ID:          field.ID,
		DatabaseID:  rootForm.GetDatabaseID(),
		RootOwnerID: rootForm.GetFormID(),
		OwnerID:     owner.GetFormID(),
		Name:        field.Name,
	}
	result := FlatForms{
		Forms: []*Form{flattenedForm},
	}
	for _, field := range field.FieldType.SubForm.GetFields() {
		flatField, err := flattenField(rootForm, subForm, field)
		if err != nil {
			return FlatForms{}, err
		}
		result = result.Add(flatField)
	}
	return result, nil
}

// flattenField converts a types.FieldDefinition into a flat storage structure
func flattenField(rootForm, form types.FormInterface, field *types.FieldDefinition) (FlatForms, error) {
	fieldKind, err := field.FieldType.GetFieldKind()
	if err != nil {
		return FlatForms{}, err
	}
	switch fieldKind {
	case types.FieldKindText:
		return flattenTextField(rootForm, form, field)
	case types.FieldKindSubForm:
		return flattenSubFormField(rootForm, form, field)
	case types.FieldKindReference:
		return flattenReferenceField(rootForm, form, field)
	case types.FieldKindMultilineText:
		return flattenMultilineTextField(rootForm, form, field)
	case types.FieldKindDate:
		return flattenDateField(rootForm, form, field)
	case types.FieldKindQuantity:
		return flattenQuantityField(rootForm, form, field)
	case types.FieldKindMonth:
		return flattenMonthField(rootForm, form, field)
	case types.FieldKindWeek:
		return flattenWeekField(rootForm, form, field)
	case types.FieldKindSingleSelect:
		return flattenSingleSelectField(rootForm, form, field)
	case types.FieldKindMultiSelect:
		return flattenMultiSelectField(rootForm, form, field)
	case types.FieldKindCheckbox:
		return flattenCheckboxField(rootForm, form, field)
	}
	return FlatForms{}, fmt.Errorf("unprocessable field kind %v", fieldKind)
}

// getFieldsDefault returns the default storage representation for a types.FieldDefinition
func getFieldsDefault(rootForm, form types.FormInterface, textField *types.FieldDefinition) *Field {
	return &Field{
		ID:          textField.ID,
		DatabaseID:  form.GetDatabaseID(),
		FormID:      form.GetFormID(),
		RootFormID:  rootForm.GetFormID(),
		Name:        textField.Name,
		Description: textField.Description,
		Code:        textField.Code,
		Key:         textField.Key,
		Required:    textField.Required,
	}
}

// flattenTextField flattens a types.FieldTypeText
func flattenTextField(rootForm, form types.FormInterface, field *types.FieldDefinition) (FlatForms, error) {
	storedField := getFieldsDefault(rootForm, form, field)
	storedField.Type = types.FieldKindText
	return FlatForms{Fields: []*Field{storedField}}, nil
}

// flattenMultilineTextField flattens a types.FieldTypeMultilineText
func flattenMultilineTextField(rootForm, form types.FormInterface, field *types.FieldDefinition) (FlatForms, error) {
	storedField := getFieldsDefault(rootForm, form, field)
	storedField.Type = types.FieldKindMultilineText
	return FlatForms{Fields: []*Field{storedField}}, nil
}

// flattenDateField flattens a types.FieldTypeDate
func flattenDateField(rootForm, form types.FormInterface, field *types.FieldDefinition) (FlatForms, error) {
	storedField := getFieldsDefault(rootForm, form, field)
	storedField.Type = types.FieldKindDate
	return FlatForms{Fields: []*Field{storedField}}, nil
}

// flattenMonthField flattens a types.FieldTypeMonth
func flattenMonthField(rootForm, form types.FormInterface, field *types.FieldDefinition) (FlatForms, error) {
	storedField := getFieldsDefault(rootForm, form, field)
	storedField.Type = types.FieldKindMonth
	return FlatForms{Fields: []*Field{storedField}}, nil
}

// flattenWeekField flattens a types.FieldTypeWeek
func flattenWeekField(rootForm, form types.FormInterface, field *types.FieldDefinition) (FlatForms, error) {
	storedField := getFieldsDefault(rootForm, form, field)
	storedField.Type = types.FieldKindWeek
	return FlatForms{Fields: []*Field{storedField}}, nil
}

// flattenQuantityField flattens a types.FieldTypeQuantity
func flattenQuantityField(rootForm, form types.FormInterface, field *types.FieldDefinition) (FlatForms, error) {
	storedField := getFieldsDefault(rootForm, form, field)
	storedField.Type = types.FieldKindQuantity
	return FlatForms{Fields: []*Field{storedField}}, nil
}

// flattenReferenceField flattens a types.FieldTypeReference
func flattenReferenceField(rootForm, form types.FormInterface, field *types.FieldDefinition) (FlatForms, error) {
	storedField := getFieldsDefault(rootForm, form, field)
	storedField.Type = types.FieldKindReference
	storedField.ReferencedDatabaseID = pointers.String(field.FieldType.Reference.DatabaseID)
	storedField.ReferencedFormID = pointers.String(field.FieldType.Reference.FormID)
	return FlatForms{Fields: []*Field{storedField}}, nil
}

// flattenSingleSelectField flattens a types.FieldTypeSingleSelect
func flattenSingleSelectField(rootForm, form types.FormInterface, field *types.FieldDefinition) (FlatForms, error) {
	storedField := getFieldsDefault(rootForm, form, field)
	storedField.Type = types.FieldKindSingleSelect
	flattened := FlatForms{Fields: []*Field{storedField}}
	flattenedOptions, err := flattenSelectOptions(rootForm, form, field.ID, field.FieldType.SingleSelect.Options)
	if err != nil {
		return FlatForms{}, err
	}
	return flattened.Add(flattenedOptions), nil
}

// flattenMultiSelectField flattens a types.FieldTypeMultiSelect
func flattenMultiSelectField(rootForm, form types.FormInterface, field *types.FieldDefinition) (FlatForms, error) {
	storedField := getFieldsDefault(rootForm, form, field)
	storedField.Type = types.FieldKindMultiSelect
	flattened := FlatForms{Fields: []*Field{storedField}}
	flattenedOptions, err := flattenSelectOptions(rootForm, form, field.ID, field.FieldType.MultiSelect.Options)
	if err != nil {
		return FlatForms{}, err
	}
	return flattened.Add(flattenedOptions), nil
}

// flattenSelectOptions flattens the options of a multi or single select field
func flattenSelectOptions(rootForm, form types.FormInterface, fieldID string, options []*types.SelectOption) (FlatForms, error) {
	flattened := FlatForms{}
	for _, option := range options {
		flattened = flattened.Add(FlatForms{
			Options: []*Option{
				{
					ID:         option.ID,
					DatabaseID: form.GetDatabaseID(),
					RootFormID: rootForm.GetFormID(),
					FormID:     form.GetFormID(),
					FieldID:    fieldID,
					Name:       option.Name,
				},
			},
		})
	}
	return flattened, nil
}

// flattenSubFormField flattens a types.FieldTypeSubForm
func flattenSubFormField(rootForm, form types.FormInterface, field *types.FieldDefinition) (FlatForms, error) {
	storedField := getFieldsDefault(rootForm, form, field)
	storedField.Type = types.FieldKindSubForm
	storedField.SubFormID = pointers.String(field.ID)
	flatField := FlatForms{Fields: []*Field{storedField}}
	subForm, err := form.FindSubForm(field.ID)
	if err != nil {
		return FlatForms{}, err
	}
	flatSubForm, err := flattenSubForm(rootForm, form, subForm, field)
	if err != nil {
		return FlatForms{}, err
	}
	return flatField.Add(flatSubForm), nil
}

// flattenCheckboxField flattens a types.FieldTypeCheckbox
func flattenCheckboxField(rootForm, form types.FormInterface, field *types.FieldDefinition) (FlatForms, error) {
	storedField := getFieldsDefault(rootForm, form, field)
	storedField.Type = types.FieldKindCheckbox
	return FlatForms{Fields: []*Field{storedField}}, nil
}
