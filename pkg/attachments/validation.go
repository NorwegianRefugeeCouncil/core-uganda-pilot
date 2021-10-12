package attachments

import (
	"fmt"
	"github.com/nrc-no/core/pkg/validation"
)

func ValidateAttachment(attachment *Attachment, path *validation.Path) validation.ErrorList {
	errList := validation.ErrorList{}

	// Validate UUIDs
	uuids := map[string]string{
		"id":           attachment.ID,
		"attachedToId": attachment.AttachedToID,
	}

	for name, uuid := range uuids {
		if uuid != "" && !validation.IsValidUUID(uuid) {
			errList = append(errList, validation.Invalid(path.Child(name), uuid, fmt.Sprintf("%s is not a valid UUID", name)))
		}
	}

	return errList
}
