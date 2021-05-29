package scheme

import (
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	discoveryv1 "github.com/nrc-no/core/api/pkg/apis/discovery/v1"
	coremetav1 "github.com/nrc-no/core/api/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var Scheme = runtime.NewScheme()
var Codecs = serializer.NewCodecFactory(Scheme)
var ParameterCodec = runtime.NewParameterCodec(Scheme)
var localSchemeBuilder = runtime.SchemeBuilder{
	corev1.AddToScheme,
	discoveryv1.AddToScheme,
}

var AddToScheme = localSchemeBuilder.AddToScheme

func init() {
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
	coremetav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
	utilruntime.Must(AddToScheme(Scheme))
}
