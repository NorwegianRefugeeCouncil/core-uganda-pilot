package install

import (
	corefuzzer "github.com/nrc-no/coreapi/pkg/apis/core/fuzzer"
	"k8s.io/apimachinery/pkg/api/apitesting/roundtrip"
	"testing"
)

func TestRoundTripTypes(t *testing.T) {
	roundtrip.RoundTripTestForAPIGroup(t, Install, corefuzzer.Funcs)
}
