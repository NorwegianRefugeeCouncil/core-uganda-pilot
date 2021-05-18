package fuzzer

import (
	fuzz "github.com/google/gofuzz"
	"github.com/nrc-no/coreapi/pkg/apis/core"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var Funcs = func(codecs serializer.CodecFactory) []interface{} {
	return []interface{}{
		func(s *core.FormDefinitionSpec, c fuzz.Continue) {
			c.FuzzNoCustom(s)
		},
	}
}
