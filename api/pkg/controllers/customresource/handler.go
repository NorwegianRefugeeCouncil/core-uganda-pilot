package customresource

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/api/pkg/apis/core"
	v1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	informers "github.com/nrc-no/core/api/pkg/client/informers/core/v1"
	listers "github.com/nrc-no/core/api/pkg/client/listers/core/v1"
	"github.com/nrc-no/core/api/pkg/customresource"
	"github.com/nrc-no/core/api/pkg/customresource/conversion"
	"github.com/nrc-no/core/api/pkg/endpoints/handlers"
	"github.com/nrc-no/core/api/pkg/endpoints/request"
	structuralschema "github.com/nrc-no/core/api/pkg/openapi"
	"github.com/nrc-no/core/api/pkg/openapi/defaulting"
	"github.com/nrc-no/core/api/pkg/openapi/objectmeta"
	"github.com/nrc-no/core/api/pkg/openapi/pruning"
	customresourceregistry "github.com/nrc-no/core/api/pkg/registry/core/customresource"
	"github.com/nrc-no/core/api/pkg/registry/generic"
	store2 "github.com/nrc-no/core/api/pkg/store"
	schemaobjectmeta "k8s.io/apiextensions-apiserver/pkg/apiserver/schema/objectmeta"
	"k8s.io/apiextensions-apiserver/pkg/crdserverscheme"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/runtime/serializer/versioning"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	utilwaitgroup "k8s.io/apimachinery/pkg/util/waitgroup"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// crdHandler is a generic REST handler that is able to dynamically serve CustomResources
// as they are added to the registry
type crdHandler struct {
	customStorageLock       sync.Mutex
	customStorage           atomic.Value
	restOptionsGetter       generic.RESTOptionsGetter
	codecs                  runtime.NegotiatedSerializer
	converterFactory        *conversion.CRConverterFactory
	versionDiscoveryHandler *CRDVersionDiscoveryHandler
	groupDiscoveryHandler   *CRDGroupDiscoveryHandler
	crdLister               listers.CustomResourceDefinitionLister
	hasSynced               func() bool
	scheme                  *runtime.Scheme
}

// crdInfo contains information about a CustomResource (for serving purposes). It is lazily-created
// as the requests are coming in, and stored in a cache as they are a bit expensive to build
type crdInfo struct {
	spec          *v1.CustomResourceDefinitionSpec
	acceptedNames *v1.CustomResourceDefinitionNames
	deprecated    map[string]bool
	storages      map[string]*customresourceregistry.CustomResourceStore
	requestScopes map[string]*handlers.RequestScope
	waitgroup     *utilwaitgroup.SafeWaitGroup
}

// Maps a CustomResourceDefinition UID to it's crdInfo
type crdStorageMap map[types.UID]*crdInfo

// clone clones the above crdStorageMap
func (in crdStorageMap) clone() crdStorageMap {
	if in == nil {
		return nil
	}
	out := make(crdStorageMap, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

// NewCustomResourceDefinitionHandler Builds a crdHandler
func NewCustomResourceDefinitionHandler(
	restOptionsGetter generic.RESTOptionsGetter,
	codecs runtime.NegotiatedSerializer,
	crdInformer informers.CustomResourceDefinitionInformer,
	versionDiscoveryHandler *CRDVersionDiscoveryHandler,
	groupDiscoveryHandler *CRDGroupDiscoveryHandler,
	scheme *runtime.Scheme,
) (*crdHandler, error) {
	ret := &crdHandler{
		customStorageLock:       sync.Mutex{},
		customStorage:           atomic.Value{},
		restOptionsGetter:       restOptionsGetter,
		crdLister:               crdInformer.Lister(),
		hasSynced:               crdInformer.Informer().HasSynced,
		codecs:                  codecs,
		versionDiscoveryHandler: versionDiscoveryHandler,
		groupDiscoveryHandler:   groupDiscoveryHandler,
		scheme:                  scheme,
	}
	ret.customStorage.Store(crdStorageMap{})
	crConverterFactory, err := conversion.NewCRConverterFactory()
	if err != nil {
		return nil, err
	}
	ret.converterFactory = crConverterFactory
	return ret, nil
}

// ServeHTTP implements the http.Handler interface and is able to serve
// http requests to deliver CustomResources. It mimics the mechanisms
// of the server.APIInstaller, though, is able to do this at runtime.
func (r *crdHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	// Retrieves the request.RequestInfo info from the context.
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		responsewriters.ErrorNegotiated(
			errors.NewInternalError(fmt.Errorf("no RequestInfo found in the context")),
			r.codecs,
			schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion},
			w, req,
		)
		return
	}

	// If the request is not targeting a resource, then it is either
	// 1. A group discovery request eg. /apis/{group}
	// 2. A version discovery request ev. /apis/{group}/{version}
	// 3. Otherwise, return not found
	if !requestInfo.IsResourceRequest {
		pathParts := splitPath(requestInfo.Path)

		// Assumes this is 2. a version discovery request
		if len(pathParts) == 3 {
			if !r.hasSynced() {
				responsewriters.ErrorNegotiated(serverStartingError(), r.codecs, schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}, w, req)
				return
			}
			r.versionDiscoveryHandler.ServeHTTP(w, req)
			return
		}

		// Assumes this is 1. a group discovery request
		if len(pathParts) == 2 {
			if !r.hasSynced() {
				responsewriters.ErrorNegotiated(serverStartingError(), r.codecs, schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}, w, req)
				return
			}
			r.groupDiscoveryHandler.ServeHTTP(w, req)
			return
		}

		// not found
		http.NotFoundHandler().ServeHTTP(w, req)
		return
	}

	// This is the name of the CRD, as the CRD **must** have a name
	// equal to {resource}.{group}
	crdName := requestInfo.Resource + "." + requestInfo.APIGroup

	// retrieve the CustomResourceDefinition for that name
	crd, err := r.crdLister.Get(crdName)
	if err != nil {
		responsewriters.ErrorNegotiated(
			err,
			r.codecs,
			schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion},
			w,
			req,
		)
		return
	}

	// gets or build the crdInfo for the request
	crInfo, err := r.getOrCreateServingInfoFor(ctx, crd.UID, crd.Name)
	if err != nil {
		responsewriters.ErrorNegotiated(
			err,
			r.codecs,
			schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion},
			w,
			req,
		)
		return
	}

	// Dynamically serves the custom resource request
	r.serveResource(w, req, requestInfo, crInfo, crd)(w, req)

}

// serveResource is a dynamic http handler able to serve http requests that target
// runtime CustomResources
func (r *crdHandler) serveResource(w http.ResponseWriter, req *http.Request, requestInf *request.RequestInfo, crdInfo *crdInfo, crd *v1.CustomResourceDefinition) http.HandlerFunc {

	// Retrieves the request scope from the pre-built map
	scope := crdInfo.requestScopes[requestInf.APIVersion]

	// Retrieves the request storage from the pre-built storage map
	storage := crdInfo.storages[requestInf.APIVersion].CustomResource

	// Map the request to the appropriate handler
	switch requestInf.Verb {
	case "get":
		return handlers.GetResource(scope, storage)
	case "list":
		return handlers.ListResource(scope, storage, nil)
	case "create":
		return handlers.CreateResource(scope, storage)
	case "update":
		return handlers.UpdateResource(scope, storage)
	case "delete":
		return handlers.DeleteResource(scope, storage)
	default:
		responsewriters.ErrorNegotiated(
			errors.NewMethodNotSupported(schema.GroupResource{Group: requestInf.APIGroup, Resource: requestInf.Resource}, requestInf.Verb),
			r.codecs, schema.GroupVersion{Group: requestInf.APIGroup, Version: requestInf.APIVersion}, w, req,
		)
		return nil
	}
}

// getOrCreateServingInfoFor will either return an already built crdInfo for the given uid,
// otherwise it will build it and store it in the crdStorageMap
func (r *crdHandler) getOrCreateServingInfoFor(ctx context.Context, uid types.UID, name string) (*crdInfo, error) {

	// tries to find the storage map if already exists
	storageMap := r.customStorage.Load().(crdStorageMap)
	if ret, ok := storageMap[uid]; ok {
		return ret, nil
	}

	// lock the storage map because we're accessing this concurrently possibly
	r.customStorageLock.Lock()
	defer r.customStorageLock.Unlock()

	// gets the CustomResourceDefinition
	crd, err := r.crdLister.Get(name)
	if err != nil {
		return nil, err
	}

	// Tries again to find the storage map. Perhaps it was
	// loaded in the meantime?
	storageMap = r.customStorage.Load().(crdStorageMap)
	if ret, ok := storageMap[crd.UID]; ok {
		return ret, nil
	}

	// builds objects needed for creating the crdInfo object
	requestScopes := map[string]*handlers.RequestScope{}
	storages := map[string]*customresourceregistry.CustomResourceStore{}
	structuralSchemes := map[string]*structuralschema.Structural{}

	// loops through all crd versions
	for _, v := range crd.Spec.Versions {

		// Gets the CustomResourceDefinitionValidation
		val, err := customresource.GetSchemaForVersion(crd, v.Name)
		if err != nil {
			utilruntime.HandleError(err)
			return nil, fmt.Errorf("the server could not serve the CR schema: %v", err)
		}

		// Converts the schema to the internal (hub) api version
		internalValidation := &core.CustomResourceDefinitionValidation{}
		if err := v1.Convert_v1_CustomResourceDefinitionValidation_To_core_CustomResourceDefinitionValidation(&val, internalValidation, nil); err != nil {
			return nil, fmt.Errorf("failed converting CRD validation to internal version: %v", err)
		}

		// Converts the schema to the Structural type
		s, err := structuralschema.NewStructural(&internalValidation.OpenAPIV3Schema)
		if err != nil {
			utilruntime.HandleError(fmt.Errorf("fialed to convert schema to structural: %v", err))
			return nil, fmt.Errorf("the server could not serve the CR schema: %v", err)
		}

		// Save the structural schema
		structuralSchemes[v.Name] = s

	}

	// Retrieves converters for the CRD
	safeConverter, unsafeConverter, err := r.converterFactory.NewConverter(crd)
	if err != nil {
		return nil, err
	}

	// Finds which version we'll use for storage (to store in the database)
	// Since a CustomResourceDefinition should define at least one version
	// that will be the "storage" version (persisted in the database)
	var storageVersion string
	for _, v := range crd.Spec.Versions {
		if v.Storage {
			storageVersion = v.Name
		}
	}
	// No storageVersion found in this CRD. Should not happen but still.
	if len(storageVersion) == 0 {
		return nil, fmt.Errorf("no storage version for CR")
	}

	// Loop through all versions
	for _, v := range crd.Spec.Versions {

		// We need a parameterScheme to be able to decode ListResourcesOptions, GetOptions, etc
		// from the URL query parameters
		parameterScheme := runtime.NewScheme()
		parameterScheme.AddUnversionedTypes(schema.GroupVersion{Group: crd.Spec.Group, Version: v.Name},
			&metav1.ListOptions{},
			&metav1.GetOptions{},
			&metav1.DeleteOptions{},
		)
		parameterCodec := runtime.NewParameterCodec(parameterScheme)

		// This is the CR resource
		resource := schema.GroupVersionResource{Group: crd.Spec.Group, Version: v.Name, Resource: crd.Spec.Names.Plural}

		// This is the CR kind
		kind := schema.GroupVersionKind{Group: crd.Spec.Group, Version: v.Name, Kind: crd.Spec.Names.Kind}

		// The CR typer
		typer := newUnstructuredObjectTyper(parameterScheme)

		// The CR creator able to create new instances of a CR on the fly
		creator := unstructuredCreator{}

		// Gets the v1 validation schema
		validationSchema, err := customresource.GetSchemaForVersion(crd, v.Name)
		if err != nil {
			utilruntime.HandleError(err)
			return nil, fmt.Errorf("the server could not serve the CR schema")
		}

		// Converts the v1 validation schema to internal (hub) api version
		var internalValidationSchema = &core.CustomResourceDefinitionValidation{}
		if err := v1.Convert_v1_CustomResourceDefinitionValidation_To_core_CustomResourceDefinitionValidation(&validationSchema, internalValidationSchema, nil); err != nil {
			return nil, fmt.Errorf("failed to convert CRD validation to internal version: %v", err)
		}

		//TOOD: validator, _, err := NewSchemaVa

		// Builds the CR REST storage interface
		storage, err := customresourceregistry.NewStorage(
			resource.GroupResource(),
			kind,
			schema.GroupVersionKind{Group: crd.Spec.Group, Version: v.Name, Kind: crd.Spec.Names.Kind + "List"},
			customresourceregistry.NewStrategy(
				typer,
				false,
				kind,
				nil,
				structuralSchemes,
			),
			crdConversionRESTOptionsGetter{
				RESTOptionsGetter:     r.restOptionsGetter,
				converter:             safeConverter,
				decoderVersion:        schema.GroupVersion{Group: crd.Spec.Group, Version: v.Name},
				encoderVersion:        schema.GroupVersion{Group: crd.Spec.Group, Version: storageVersion},
				structuralSchemas:     structuralSchemes,
				structuralSchemaGK:    kind.GroupKind(),
				preserveUnknownFields: true,
				scheme:                r.scheme,
			})
		if err != nil {
			return nil, err
		}

		// Store the CR storage
		storages[v.Name] = storage

		// Builds the negotiated serializer for the CR
		negotiatedSerializer := unstructuredNegotiatedSerializer{
			scheme:                r.scheme,
			typer:                 typer,
			creator:               creator,
			converter:             safeConverter,
			structuralSchemas:     structuralSchemes,
			structuralSchemaGK:    kind.GroupKind(),
			preserveUnknownFields: true,
		}
		var standardSerializers []runtime.SerializerInfo
		for _, s := range negotiatedSerializer.SupportedMediaTypes() {
			if s.MediaType == runtime.ContentTypeProtobuf {
				continue
			}
			standardSerializers = append(standardSerializers, s)
		}

		// Builds the RequestScope scope for the CR
		requestScopes[v.Name] = &handlers.RequestScope{
			Namer:          handlers.ContextBasedNaming{},
			Serializer:     negotiatedSerializer,
			ParameterCodec: parameterCodec,
			Creater:        creator,
			Convertor:      safeConverter,
			Defaulter: unstructuredDefaulter{
				parameterScheme,
				structuralSchemes,
				kind.GroupKind()},
			Typer:           typer,
			UnsafeConvertor: unsafeConverter,
			Resource:        schema.GroupVersionResource{Group: crd.Spec.Group, Version: v.Name, Resource: crd.Spec.Names.Plural},
			Kind:            kind,
			HubGroupVersion: kind.GroupVersion(),
		}

	}

	// Finalize the crdInfo
	ret := &crdInfo{
		spec:          &crd.Spec,
		deprecated:    map[string]bool{},
		storages:      storages,
		requestScopes: requestScopes,
		waitgroup:     &utilwaitgroup.SafeWaitGroup{},
		acceptedNames: &crd.Spec.Names,
	}

	// Clone the storageMap and store in
	storageMap2 := storageMap.clone()
	storageMap2[crd.UID] = ret
	r.customStorage.Store(storageMap2)

	return ret, nil

}

// crdConversionRESTOptionsGetter is an RESTOptionsGetter that
// uses a custom codec that is suitable to handle runtime
// objects (not compile time)
type crdConversionRESTOptionsGetter struct {
	generic.RESTOptionsGetter
	converter             runtime.ObjectConvertor
	encoderVersion        schema.GroupVersion
	decoderVersion        schema.GroupVersion
	structuralSchemas     map[string]*structuralschema.Structural // by version
	structuralSchemaGK    schema.GroupKind
	preserveUnknownFields bool
	scheme                *runtime.Scheme
}

func (t crdConversionRESTOptionsGetter) GetRESTOptions(resource schema.GroupResource) (generic.RESTOptions, error) {
	ret, err := t.RESTOptionsGetter.GetRESTOptions(resource)
	if err == nil {
		d := schemaCoercingDecoder{delegate: ret.StorageConfig.Codec, validator: unstructuredSchemaCoercer{
			// drop invalid fields while decoding old CRs (before we haven't had any ObjectMeta validation)
			dropInvalidMetadata:   true,
			repairGeneration:      true,
			structuralSchemas:     t.structuralSchemas,
			structuralSchemaGK:    t.structuralSchemaGK,
			preserveUnknownFields: t.preserveUnknownFields,
		}}
		c := schemaCoercingConverter{delegate: t.converter, validator: unstructuredSchemaCoercer{
			structuralSchemas:     t.structuralSchemas,
			structuralSchemaGK:    t.structuralSchemaGK,
			preserveUnknownFields: t.preserveUnknownFields,
		}}
		ret.StorageConfig.Codec = versioning.NewCodec(
			ret.StorageConfig.Codec,
			d,
			c,
			&unstructuredCreator{},
			crdserverscheme.NewUnstructuredObjectTyper(),
			&unstructuredDefaulter{
				delegate:           t.scheme,
				structuralSchemaGK: t.structuralSchemaGK,
				structuralSchemas:  t.structuralSchemas,
			},
			t.encoderVersion,
			t.decoderVersion,
			"crdRESTOptions",
		)
	}
	return ret, err
}

type unstructuredNegotiatedSerializer struct {
	scheme    *runtime.Scheme
	typer     runtime.ObjectTyper
	creator   runtime.ObjectCreater
	converter runtime.ObjectConvertor

	structuralSchemas     map[string]*structuralschema.Structural // by version
	structuralSchemaGK    schema.GroupKind
	preserveUnknownFields bool
}

func (s unstructuredNegotiatedSerializer) SupportedMediaTypes() []runtime.SerializerInfo {
	return []runtime.SerializerInfo{
		{
			MediaType:        "application/json",
			MediaTypeType:    "application",
			MediaTypeSubType: "json",
			EncodesAsText:    true,
			Serializer:       json.NewSerializer(json.DefaultMetaFactory, s.creator, s.typer, false),
			PrettySerializer: json.NewSerializer(json.DefaultMetaFactory, s.creator, s.typer, true),
			StreamSerializer: &runtime.StreamSerializerInfo{
				EncodesAsText: true,
				Serializer:    json.NewSerializer(json.DefaultMetaFactory, s.creator, s.typer, false),
				Framer:        json.Framer,
			},
		},
		{
			MediaType:        "application/yaml",
			MediaTypeType:    "application",
			MediaTypeSubType: "yaml",
			EncodesAsText:    true,
			Serializer:       json.NewYAMLSerializer(json.DefaultMetaFactory, s.creator, s.typer),
		},
		// We're not supporting protobuf for now
		//{
		//	MediaType:        "application/vnd.kubernetes.protobuf",
		//	MediaTypeType:    "application",
		//	MediaTypeSubType: "vnd.kubernetes.protobuf",
		//	Serializer:       protobuf.NewSerializer(s.creator, s.typer),
		//	StreamSerializer: &runtime.StreamSerializerInfo{
		//		Serializer: protobuf.NewRawSerializer(s.creator, s.typer),
		//		Framer:     protobuf.LengthDelimitedFramer,
		//	},
		//},
	}
}

func (s unstructuredNegotiatedSerializer) EncoderForVersion(encoder runtime.Encoder, gv runtime.GroupVersioner) runtime.Encoder {
	return versioning.NewCodec(encoder, nil, s.converter, s.scheme, s.scheme, s.scheme, gv, nil, "crdNegotiatedSerializer")
}

func (s unstructuredNegotiatedSerializer) DecoderToVersion(decoder runtime.Decoder, gv runtime.GroupVersioner) runtime.Decoder {
	d := schemaCoercingDecoder{delegate: decoder, validator: unstructuredSchemaCoercer{structuralSchemas: s.structuralSchemas, structuralSchemaGK: s.structuralSchemaGK, preserveUnknownFields: s.preserveUnknownFields}}
	return versioning.NewCodec(nil, d, runtime.UnsafeObjectConvertor(s.scheme), s.scheme, s.scheme, unstructuredDefaulter{
		delegate:           s.scheme,
		structuralSchemas:  s.structuralSchemas,
		structuralSchemaGK: s.structuralSchemaGK,
	}, nil, gv, "unstructuredNegotiatedSerializer")
}

type UnstructuredObjectTyper struct {
	Delegate          runtime.ObjectTyper
	UnstructuredTyper runtime.ObjectTyper
}

func newUnstructuredObjectTyper(Delegate runtime.ObjectTyper) UnstructuredObjectTyper {
	return UnstructuredObjectTyper{
		Delegate:          Delegate,
		UnstructuredTyper: crdserverscheme.NewUnstructuredObjectTyper(),
	}
}

func (t UnstructuredObjectTyper) ObjectKinds(obj runtime.Object) ([]schema.GroupVersionKind, bool, error) {
	// Delegate for things other than Unstructured.
	if _, ok := obj.(runtime.Unstructured); !ok {
		return t.Delegate.ObjectKinds(obj)
	}
	return t.UnstructuredTyper.ObjectKinds(obj)
}

func (t UnstructuredObjectTyper) Recognizes(gvk schema.GroupVersionKind) bool {
	return t.Delegate.Recognizes(gvk) || t.UnstructuredTyper.Recognizes(gvk)
}

type unstructuredCreator struct{}

func (c unstructuredCreator) New(kind schema.GroupVersionKind) (runtime.Object, error) {
	ret := &unstructured.Unstructured{}
	ret.SetGroupVersionKind(kind)
	return ret, nil
}

type unstructuredDefaulter struct {
	delegate           runtime.ObjectDefaulter
	structuralSchemas  map[string]*structuralschema.Structural // by version
	structuralSchemaGK schema.GroupKind
}

func (d unstructuredDefaulter) Default(in runtime.Object) {
	// Delegate for things other than Unstructured, and other GKs
	u, ok := in.(runtime.Unstructured)
	if !ok || u.GetObjectKind().GroupVersionKind().GroupKind() != d.structuralSchemaGK {
		d.delegate.Default(in)
		return
	}

	defaulting.Default(u.UnstructuredContent(), d.structuralSchemas[u.GetObjectKind().GroupVersionKind().Version])
}

// schemaCoercingDecoder calls the delegate decoder, and then applies the Unstructured schema validator
// to coerce the schema.
type schemaCoercingDecoder struct {
	delegate  runtime.Decoder
	validator unstructuredSchemaCoercer
}

var _ runtime.Decoder = schemaCoercingDecoder{}

func (d schemaCoercingDecoder) Decode(data []byte, defaults *schema.GroupVersionKind, into runtime.Object) (runtime.Object, *schema.GroupVersionKind, error) {
	obj, gvk, err := d.delegate.Decode(data, defaults, into)
	if err != nil {
		return nil, gvk, err
	}
	if u, ok := obj.(*unstructured.Unstructured); ok {
		if err := d.validator.apply(u); err != nil {
			return nil, gvk, err
		}
	}

	return obj, gvk, nil
}

// schemaCoercingConverter calls the delegate converter and applies the Unstructured validator to
// coerce the schema.
type schemaCoercingConverter struct {
	delegate  runtime.ObjectConvertor
	validator unstructuredSchemaCoercer
}

var _ runtime.ObjectConvertor = schemaCoercingConverter{}

func (v schemaCoercingConverter) Convert(in, out, context interface{}) error {
	if err := v.delegate.Convert(in, out, context); err != nil {
		return err
	}

	if u, ok := out.(*unstructured.Unstructured); ok {
		if err := v.validator.apply(u); err != nil {
			return err
		}
	}

	return nil
}

func (v schemaCoercingConverter) ConvertToVersion(in runtime.Object, gv runtime.GroupVersioner) (runtime.Object, error) {
	out, err := v.delegate.ConvertToVersion(in, gv)
	if err != nil {
		return nil, err
	}

	if u, ok := out.(*unstructured.Unstructured); ok {
		if err := v.validator.apply(u); err != nil {
			return nil, err
		}
	}

	return out, nil
}

func (v schemaCoercingConverter) ConvertFieldLabel(gvk schema.GroupVersionKind, label, value string) (string, string, error) {
	return v.delegate.ConvertFieldLabel(gvk, label, value)
}

// unstructuredSchemaCoercer adds to unstructured unmarshalling what json.Unmarshal does
// in addition for native types when decoding into Golang structs:
//
// - validating and pruning ObjectMeta
// - generic pruning of unknown fields following a structural schema
// - removal of non-defaulted non-nullable null map values.
type unstructuredSchemaCoercer struct {
	dropInvalidMetadata bool
	repairGeneration    bool

	structuralSchemas     map[string]*structuralschema.Structural
	structuralSchemaGK    schema.GroupKind
	preserveUnknownFields bool
}

func (v *unstructuredSchemaCoercer) apply(u *unstructured.Unstructured) error {
	// save implicit meta fields that don't have to be specified in the validation spec
	kind, foundKind, err := unstructured.NestedString(u.UnstructuredContent(), "kind")
	if err != nil {
		return err
	}
	apiVersion, foundApiVersion, err := unstructured.NestedString(u.UnstructuredContent(), "apiVersion")
	if err != nil {
		return err
	}
	objectMeta, foundObjectMeta, err := schemaobjectmeta.GetObjectMeta(u.Object, v.dropInvalidMetadata)
	if err != nil {
		return err
	}

	// compare group and kind because also other object like DeleteCollection options pass through here
	gv, err := schema.ParseGroupVersion(apiVersion)
	if err != nil {
		return err
	}
	if gv.Group == v.structuralSchemaGK.Group && kind == v.structuralSchemaGK.Kind {
		if !v.preserveUnknownFields {
			// TODO: switch over pruning and coercing at the root to  schemaobjectmeta.Coerce too
			pruning.Prune(u.Object, v.structuralSchemas[gv.Version], false)
			defaulting.PruneNonNullableNullsWithoutDefaults(u.Object, v.structuralSchemas[gv.Version])
		}
		if err := objectmeta.Coerce(nil, u.Object, v.structuralSchemas[gv.Version], false, v.dropInvalidMetadata); err != nil {
			return err
		}
		// fixup missing generation in very old CRs
		if v.repairGeneration && objectMeta.Generation == 0 {
			objectMeta.Generation = 1
		}
	}

	// restore meta fields, starting clean
	if foundKind {
		u.SetKind(kind)
	}
	if foundApiVersion {
		u.SetAPIVersion(apiVersion)
	}
	if foundObjectMeta {
		if err := schemaobjectmeta.SetObjectMeta(u.Object, objectMeta); err != nil {
			return err
		}
	}

	return nil
}

type CRDRESTOptionsGetter struct {
	StorageConfig           store2.Config
	StoragePrefix           string
	EnableWatchCache        bool
	DefaultWatchCacheSize   int
	EnableGarbageCollection bool
	DeleteCollectionWorkers int
	CountMetricPollPeriod   time.Duration
}

func (t CRDRESTOptionsGetter) GetRESTOptions(resource schema.GroupResource) (generic.RESTOptions, error) {
	ret := generic.RESTOptions{
		StorageConfig: &t.StorageConfig,
		//Decorator:               generic.UndecoratedStorage,
		//EnableGarbageCollection: t.EnableGarbageCollection,
		//DeleteCollectionWorkers: t.DeleteCollectionWorkers,
		//ResourcePrefix:          resource.Group + "/" + resource.Resource,
		//CountMetricPollPeriod:   t.CountMetricPollPeriod,
	}
	//if t.EnableWatchCache {
	//	ret.Decorator = genericregistry.StorageWithCacher()
	//}
	return ret, nil
}

// serverStartingError returns a ServiceUnavailble error with a retry-after time
func serverStartingError() error {
	err := errors.NewServiceUnavailable("server is starting")
	if err.ErrStatus.Details == nil {
		err.ErrStatus.Details = &metav1.StatusDetails{}
	}
	if err.ErrStatus.Details.RetryAfterSeconds == 0 {
		err.ErrStatus.Details.RetryAfterSeconds = int32(10)
	}
	return err
}
