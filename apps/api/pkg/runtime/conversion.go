package runtime

import "fmt"

// DefaultMetaV1FieldSelectorConversion auto-accepts metav1 values for name and namespace.
// A cluster scoped resource specifying namespace empty works fine and specifying a particular
// namespace will return no results, as expected.
func DefaultMetaV1FieldSelectorConversion(label, value string) (string, string, error) {
	switch label {
	case "metadata.name":
		return label, value, nil
	case "metadata.namespace":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("%q is not a known field selector: only %q, %q", label, "metadata.name", "metadata.namespace")
	}
}
