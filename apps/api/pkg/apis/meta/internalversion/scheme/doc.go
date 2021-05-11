package scheme

import (
	"github.com/nrc-no/core/apps/api/pkg/apis/meta/internalversion"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/serializer"
	utilruntime "github.com/nrc-no/core/apps/api/pkg/util/runtime"
)

var scheme = runtime.NewScheme()

var Codecs = serializer.NewCodecFactory(scheme)

var ParameterCodec = runtime.NewParameterCodec(scheme)

func init() {
	utilruntime.Must(internalversion.AddToScheme(scheme))
}
