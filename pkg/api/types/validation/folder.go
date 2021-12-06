package validation

import (
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"regexp"
	"strings"
)

const (
	errFolderNameRequired       = "Folder name is required"
	errFolderNameInvalid        = "Folder name is invalid"
	errFolderDatabaseIDRequired = "Database ID is required"
	errFolderInvalidDatabaseID  = "Invalid database ID"
	errFolderInvalidParentID    = "Invalid parent ID"
	folderNameMinLength         = 3
	folderNameMaxLength         = 64
)

var (
	folderNameRegex = regexp.MustCompile("^[a-zA-Z0-9]+( [a-zA-Z0-9]+)*$")
)

func ValidateFolder(folder *types.Folder) validation.ErrorList {
	var result validation.ErrorList

	namePath := validation.NewPath("name")
	databaseIdPath := validation.NewPath("databaseId")
	parentIdPath := validation.NewPath("parentId")

	if len(folder.Name) == 0 {
		result = append(result, validation.Required(namePath, errFolderNameRequired))
	} else {
		if len(folder.Name) > folderNameMaxLength {
			result = append(result, validation.TooLong(namePath, folder.Name, folderNameMaxLength))
		} else if len(folder.Name) < folderNameMinLength {
			result = append(result, validation.TooShort(namePath, folder.Name, folderNameMinLength))
		} else if strings.TrimSpace(folder.Name) != folder.Name || !folderNameRegex.MatchString(folder.Name) {
			result = append(result, validation.Invalid(namePath, folder.Name, errFolderNameInvalid))
		}
	}

	if len(folder.DatabaseID) == 0 {
		result = append(result, validation.Required(databaseIdPath, errFolderDatabaseIDRequired))
	} else if _, err := uuid.FromString(folder.DatabaseID); err != nil {
		result = append(result, validation.Invalid(databaseIdPath, folder.DatabaseID, errFolderInvalidDatabaseID))
	}

	if len(folder.ParentID) != 0 {
		if _, err := uuid.FromString(folder.ParentID); err != nil {
			result = append(result, validation.Invalid(parentIdPath, folder.ParentID, errFolderInvalidParentID))
		}
	}

	return result

}
