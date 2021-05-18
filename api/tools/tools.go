// +build tools

package tools

import (
	_ "k8s.io/code-generator/cmd/conversion-gen/generators"
	_ "k8s.io/code-generator/cmd/defaulter-gen"
	_ "sigs.k8s.io/controller-tools/pkg/deepcopy"
)
