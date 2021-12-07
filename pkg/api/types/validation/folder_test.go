package validation

import (
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestValidateFolder(t *testing.T) {

	aUUID := uuid.NewV4().String()

	const validName = "A Folder"
	tests := []struct {
		name   string
		folder *types.Folder
		expect validation.ErrorList
	}{
		{
			name: "valid folder without parent",
			folder: &types.Folder{
				Name:       validName,
				DatabaseID: aUUID,
			},
			expect: nil,
		}, {
			name: "valid folder with parent",
			folder: &types.Folder{
				Name:       validName,
				DatabaseID: aUUID,
				ParentID:   aUUID,
			},
			expect: nil,
		}, {
			name: "folder without name",
			folder: &types.Folder{
				DatabaseID: aUUID,
			},
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("name"), errFolderNameRequired),
			},
		}, {
			name: "folder without database id",
			folder: &types.Folder{
				Name: validName,
			},
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("databaseId"), errFolderDatabaseIDRequired),
			},
		}, {
			name: "folder with invalid database id",
			folder: &types.Folder{
				Name:       validName,
				DatabaseID: "bla",
			},
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("databaseId"), "bla", errFolderInvalidDatabaseID),
			},
		}, {
			name: "folder with invalid parent id",
			folder: &types.Folder{
				Name:       validName,
				ParentID:   "bla",
				DatabaseID: aUUID,
			},
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("parentId"), "bla", errFolderInvalidParentID),
			},
		}, {
			name: "folder with name too short",
			folder: &types.Folder{
				Name:       strings.Repeat("a", folderNameMinLength-1),
				DatabaseID: aUUID,
			},
			expect: validation.ErrorList{
				validation.TooShort(validation.NewPath("name"), strings.Repeat("a", folderNameMinLength-1), folderNameMinLength),
			},
		}, {
			name: "folder with name too long",
			folder: &types.Folder{
				Name:       strings.Repeat("a", folderNameMaxLength+1),
				DatabaseID: aUUID,
			},
			expect: validation.ErrorList{
				validation.TooLong(validation.NewPath("name"), strings.Repeat("a", folderNameMaxLength+1), folderNameMaxLength),
			},
		}, {
			name: "folder with invalid name",
			folder: &types.Folder{
				Name:       " ! ",
				DatabaseID: aUUID,
			},
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("name"), " ! ", errFolderNameInvalid),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ValidateFolder(test.folder)
			assert.Equal(t, test.expect, got)
		})
	}

}
