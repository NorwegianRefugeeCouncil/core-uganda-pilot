package handlers

import (
	"encoding/hex"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// transformDecodeError adds additional information into a bad-request api error when a decode fails.
func transformDecodeError(typer runtime.ObjectTyper, baseErr error, into runtime.Object, gvk *schema.GroupVersionKind, body []byte) error {
	objGVKs, _, err := typer.ObjectKinds(into)
	if err != nil {
		return errors.NewBadRequest(err.Error())
	}
	objGVK := objGVKs[0]
	if gvk != nil && len(gvk.Kind) > 0 {
		return errors.NewBadRequest(fmt.Sprintf("%s in version %q cannot be handled as a %s: %v", gvk.Kind, gvk.Version, objGVK.Kind, baseErr))
	}
	summary := summarizeData(body, 30)
	return errors.NewBadRequest(fmt.Sprintf("the object provided is unrecognized (must be of type %s): %v (%s)", objGVK.Kind, baseErr, summary))
}

func summarizeData(data []byte, maxLength int) string {
	switch {
	case len(data) == 0:
		return "<empty>"
	case data[0] == '{':
		if len(data) > maxLength {
			return string(data[:maxLength]) + " ..."
		}
		return string(data)
	default:
		if len(data) > maxLength {
			return hex.EncodeToString(data[:maxLength]) + " ..."
		}
		return hex.EncodeToString(data)
	}
}
