package validation

import (
	"fmt"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// TestValidateForm checks that all the validation methods are called and that,
// if these validation methods return errors, that they are included in the response
func TestValidateForm(t *testing.T) {

	tests := []struct {
		name       string
		expectPath string
		expectStr  func(form *types.FormDefinition) string
		validateFn *validateStrFn
	}{
		{
			name:       "formName",
			expectPath: "name",
			expectStr: func(form *types.FormDefinition) string {
				return form.Name
			},
			validateFn: &validateFormNameFn,
		},
		{
			name:       "folderId",
			expectPath: "folderId",
			expectStr: func(form *types.FormDefinition) string {
				return form.FolderID
			},
			validateFn: &validateFolderIdFn,
		},
		{
			name:       "databaseId",
			expectPath: "databaseId",
			expectStr: func(form *types.FormDefinition) string {
				return form.DatabaseID
			},
			validateFn: &validateDatabaseIdFn,
		},
	}

	for _, test := range tests {
		err := &validation.Error{}
		called := false
		var calledPath string
		var calledStr string
		oldFn := *test.validateFn
		*test.validateFn = func(str string, path *validation.Path) validation.ErrorList {
			called = true
			calledPath = path.String()
			calledStr = str
			return []*validation.Error{err}
		}
		form := &types.FormDefinition{
			ID:         "id",
			Code:       "code",
			DatabaseID: "databaseId",
			FolderID:   "folderId",
			Name:       "name",
			Fields:     nil,
		}
		errs := ValidateForm(form)
		*test.validateFn = oldFn
		assert.True(t, called)
		assert.Contains(t, errs, err)
		assert.Equal(t, test.expectPath, calledPath)
		assert.Equal(t, test.expectStr(form), calledStr)
	}
}

func TestValidateFormName(t *testing.T) {
	p := validation.NewPath("")
	tests := []struct {
		name     string
		formName string
		want     validation.ErrorList
	}{
		{
			name:     "valid",
			formName: "my form name",
			want:     []*validation.Error{},
		},
		{
			name:     "tooLong",
			formName: strings.Repeat("a", 65),
			want: []*validation.Error{
				validation.TooLong(p, strings.Repeat("a", 65), 64),
			},
		},
		{
			name:     "tooShort",
			formName: "a",
			want: []*validation.Error{
				validation.TooShort(p, "a", 3),
			},
		},
		{
			name:     "leadingWhiteSpaces",
			formName: " myForm",
			want: []*validation.Error{
				validation.Invalid(p, " myForm", formNameNoLeadingTrailingWhitespaces),
			},
		},
		{
			name:     "trailingWhiteSpaces",
			formName: "myForm ",
			want: []*validation.Error{
				validation.Invalid(p, "myForm ", formNameNoLeadingTrailingWhitespaces),
			},
		},
		{
			name:     "empty",
			formName: "",
			want: []*validation.Error{
				validation.Required(p, formNameRequired),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateFormName(tt.formName, p)
			assert.ElementsMatchf(t, tt.want, got, "")
		})
	}
}

func TestValidateFormDatabaseID(t *testing.T) {
	p := validation.NewPath("")
	tests := []struct {
		name       string
		databaseID string
		want       validation.ErrorList
	}{
		{
			name:       "valid",
			databaseID: uuid.NewV4().String(),
			want:       []*validation.Error{},
		},
		{
			name:       "empty",
			databaseID: "",
			want: []*validation.Error{
				validation.Required(p, formDatabaseIdRequired),
			},
		},
		{
			name:       "invalid uuid",
			databaseID: "abc",
			want: []*validation.Error{
				validation.Invalid(p, "abc", uuidInvalid),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				got := ValidateFormDatabaseID(tt.databaseID, p)
				assert.ElementsMatchf(t, tt.want, got, "")
			})
		})
	}
}

func TestValidateFormFolderID(t *testing.T) {
	p := validation.NewPath("")
	tests := []struct {
		name     string
		folderId string
		want     validation.ErrorList
	}{
		{
			name:     "valid",
			folderId: uuid.NewV4().String(),
			want:     []*validation.Error{},
		},
		{
			name:     "empty",
			folderId: "",
			want:     []*validation.Error{},
		},
		{
			name:     "invalid uuid",
			folderId: "abc",
			want: []*validation.Error{
				validation.Invalid(p, "abc", uuidInvalid),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				got := ValidateFormFolderID(tt.folderId, p)
				assert.ElementsMatchf(t, tt.want, got, "")
			})
		})
	}
}

func TestValidateFieldType(t *testing.T) {
	p := validation.NewPath("")
	tests := []struct {
		name string
		obj  types.FieldType
		want validation.ErrorList
	}{
		{
			name: "valid",
			obj:  types.FieldType{Text: &types.FieldTypeText{}},
			want: []*validation.Error{},
		},
		{
			name: "empty",
			obj:  types.FieldType{},
			want: []*validation.Error{
				validation.Invalid(p, types.FieldType{}, errOneFieldTypeRequired),
			},
		},
		{
			name: "multipleSpecified",
			obj:  types.FieldType{Text: &types.FieldTypeText{}, MultilineText: &types.FieldTypeMultilineText{}},
			want: []*validation.Error{
				validation.Invalid(p, types.FieldType{
					Text:          &types.FieldTypeText{},
					MultilineText: &types.FieldTypeMultilineText{},
				}, fmt.Sprintf(errFieldTypesMultipleF, []string{"text", "multilineText"})),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateFieldType(tt.obj, p)
			assert.ElementsMatchf(t, tt.want, got, "")
		})
	}
}
