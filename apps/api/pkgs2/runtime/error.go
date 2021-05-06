package runtime

import (
	"github.com/nrc-no/core/apps/api/pkgs2/runtime/schema"
	"reflect"
)

type notRegisteredErr struct {
	schemeName string
	gvk        schema.GroupVersionKind
	target     GroupVersioner
	t          reflect.Type
}

func NewNotRegisteredErrForKind(schemeName string, gvk schema.GroupVersionKind) error {
	return &notRegisteredErr{schemeName: schemeName, gvk: gvk}
}

func NewNotRegisteredErrForType(schemeName string, t reflect.Type) error {
	return &notRegisteredErr{schemeName: schemeName, t: t}
}

func NewNotRegisteredErrForTarget(schemeName string, t reflect.Type, target GroupVersioner) error {
	return &notRegisteredErr{schemeName: schemeName, t: t, target: target}
}

func NewNotRegisteredGVKErrForTarget(schemeName string, gvk schema.GroupVersionKind, target GroupVersioner) error {
	return &notRegisteredErr{schemeName: schemeName, gvk: gvk, target: target}
}

func (k *notRegisteredErr) Error() string {
	if k.t != nil && k.target != nil {
		return fmt.Sprintf("%v is not suitable for converting to %q in scheme %q", k.t, k.target, k.schemeName)
	}
	nullGVK := schema.GroupVersionKind{}
	if k.gvk != nullGVK && k.target != nil {
		return fmt.Sprintf("%q is not suitable for converting to %q in scheme %q", k.gvk.GroupVersion(), k.target, k.schemeName)
	}
	if k.t != nil {
		return fmt.Sprintf("no kind is registered for the type %v in scheme %q", k.t, k.schemeName)
	}
	if len(k.gvk.Kind) == 0 {
		return fmt.Sprintf("no version %q has been registered in scheme %q", k.gvk.GroupVersion(), k.schemeName)
	}
	//if k.gvk.Version == APIVersionInternal {
	//  return fmt.Sprintf("no kind %q is registered for the internal version of group %q in scheme %q", k.gvk.Kind, k.gvk.Group, k.schemeName)
	//}

	return fmt.Sprintf("no kind %q is registered for version %q in scheme %q", k.gvk.Kind, k.gvk.GroupVersion(), k.schemeName)
}

// IsNotRegisteredError returns true if the error indicates the provided
// object or input data is not registered.
func IsNotRegisteredError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*notRegisteredErr)
	return ok
}

type missingKindErr struct {
	data string
}

func NewMissingKindErr(data string) error {
	return &missingKindErr{data}
}

func (k *missingKindErr) Error() string {
	return fmt.Sprintf("Object 'Kind' is missing in '%s'", k.data)
}

// IsMissingKind returns true if the error indicates that the provided object
// is missing a 'Kind' field.
func IsMissingKind(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*missingKindErr)
	return ok
}

type missingVersionErr struct {
	data string
}

func NewMissingVersionErr(data string) error {
	return &missingVersionErr{data}
}

func (k *missingVersionErr) Error() string {
	return fmt.Sprintf("Object 'apiVersion' is missing in '%s'", k.data)
}

// IsMissingVersion returns true if the error indicates that the provided object
// is missing a 'Version' field.
func IsMissingVersion(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*missingVersionErr)
	return ok
}
