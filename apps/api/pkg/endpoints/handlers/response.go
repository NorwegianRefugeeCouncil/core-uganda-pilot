package handlers

import (
  "context"
  "fmt"
  "github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/negotiation"
  "github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/responsewriters"
  "k8s.io/apimachinery/pkg/api/errors"
  "k8s.io/apimachinery/pkg/api/meta"
  "k8s.io/apimachinery/pkg/apis/meta/internalversion/scheme"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/apimachinery/pkg/apis/meta/v1/validation"
  "k8s.io/apimachinery/pkg/runtime"
  "k8s.io/apimachinery/pkg/runtime/schema"
  "net/http"
)

func transformObject(ctx context.Context, obj runtime.Object, opts interface{}, mediaType negotiation.MediaTypeOptions, scope *RequestScope, req *http.Request) (runtime.Object, error) {
  if co, ok := obj.(runtime.CacheableObject); ok {
    if mediaType.Convert != nil {
      // Non-nil mediaType.Convert means that some conversion of the object
      // has to happen. Currently conversion may potentially modify the
      // object or assume something about it (e.g. asTable operates on
      // reflection, which won't work for any wrapper).
      // To ensure it will work correctly, let's operate on base objects
      // and not cache it for now.
      //
      // TODO: Long-term, transformObject should be changed so that it
      // implements runtime.Encoder interface.
      return doTransformObject(ctx, co.GetObject(), opts, mediaType, scope, req)
    }
  }
  return doTransformObject(ctx, obj, opts, mediaType, scope, req)
}

func doTransformObject(ctx context.Context, obj runtime.Object, opts interface{}, mediaType negotiation.MediaTypeOptions, scope *RequestScope, req *http.Request) (runtime.Object, error) {
  if _, ok := obj.(*metav1.Status); ok {
    return obj, nil
  }
  //if err := setObjectSelfLink(ctx, obj, req, scope.Namer); err != nil {
  //  return nil, err
  //}

  switch target := mediaType.Convert; {
  case target == nil:
    return obj, nil

  case target.Kind == "PartialObjectMetadata":
    return asPartialObjectMetadata(obj, target.GroupVersion())

  case target.Kind == "PartialObjectMetadataList":
    return asPartialObjectMetadataList(obj, target.GroupVersion())

  case target.Kind == "Table":
    options, ok := opts.(*metav1.TableOptions)
    if !ok {
      return nil, fmt.Errorf("unexpected TableOptions, got %T", opts)
    }
    return asTable(ctx, obj, options, scope, target.GroupVersion())

  default:
    accepted, _ := negotiation.MediaTypesForSerializer(scheme.Codecs)
    err := negotiation.NewNotAcceptableError(accepted)
    return nil, err
  }
}

// optionsForTransform will load and validate any additional query parameter options for
// a conversion or return an error.
func optionsForTransform(mediaType negotiation.MediaTypeOptions, req *http.Request) (interface{}, error) {
  switch target := mediaType.Convert; {
  case target == nil:
  case target.Kind == "Table" && (target.GroupVersion() == metav1.SchemeGroupVersion):
    opts := &metav1.TableOptions{}
    if err := scheme.ParameterCodec.DecodeParameters(req.URL.Query(), metav1.SchemeGroupVersion, opts); err != nil {
      return nil, err
    }
    switch errs := validation.ValidateTableOptions(opts); len(errs) {
    case 0:
      return opts, nil
    case 1:
      return nil, errors.NewBadRequest(fmt.Sprintf("Unable to convert to Table as requested: %v", errs[0].Error()))
    default:
      return nil, errors.NewBadRequest(fmt.Sprintf("Unable to convert to Table as requested: %v", errs))
    }
  }
  return nil, nil
}

// transformResponseObject takes an object loaded from storage and performs any necessary transformations.
// Will write the complete response object.
func transformResponseObject(ctx context.Context, scope *RequestScope, req *http.Request, w http.ResponseWriter, statusCode int, mediaType negotiation.MediaTypeOptions, result runtime.Object) {
  options, err := optionsForTransform(mediaType, req)
  if err != nil {
    scope.err(err, w, req)
    return
  }
  obj, err := transformObject(ctx, result, options, mediaType, scope, req)
  if err != nil {
    scope.err(err, w, req)
    return
  }
  kind, serializer, _ := targetEncodingForTransform(scope, mediaType, req)
  responsewriters.WriteObjectNegotiated(serializer, scope, kind.GroupVersion(), w, req, statusCode, obj)
}

// targetEncodingForTransform returns the appropriate serializer for the input media type
func targetEncodingForTransform(scope *RequestScope, mediaType negotiation.MediaTypeOptions, req *http.Request) (schema.GroupVersionKind, runtime.NegotiatedSerializer, bool) {
  switch target := mediaType.Convert; {
  case target == nil:
  case (target.Kind == "PartialObjectMetadata" || target.Kind == "PartialObjectMetadataList" || target.Kind == "Table") &&
    (target.GroupVersion() == metav1.SchemeGroupVersion):
    return *target, scheme.Codecs, true
  }
  return scope.Kind, scope.Serializer, false
}

func asTable(ctx context.Context, result runtime.Object, opts *metav1.TableOptions, scope *RequestScope, groupVersion schema.GroupVersion) (runtime.Object, error) {
  switch groupVersion {
  case metav1.SchemeGroupVersion:
  default:
    return nil, newNotAcceptableError(fmt.Sprintf("no Table exists in group version %s", groupVersion))
  }

  obj, err := scope.TableConvertor.ConvertToTable(ctx, result, opts)
  if err != nil {
    return nil, err
  }

  table := (*metav1.Table)(obj)

  for i := range table.Rows {
    item := &table.Rows[i]
    switch opts.IncludeObject {
    case metav1.IncludeObject:
      item.Object.Object, err = scope.Convertor.ConvertToVersion(item.Object.Object, scope.Kind.GroupVersion())
      if err != nil {
        return nil, err
      }
    // TODO: rely on defaulting for the value here?
    case metav1.IncludeMetadata, "":
      m, err := meta.Accessor(item.Object.Object)
      if err != nil {
        return nil, err
      }
      // TODO: turn this into an internal type and do conversion in order to get object kind automatically set?
      partial := meta.AsPartialObjectMetadata(m)
      partial.GetObjectKind().SetGroupVersionKind(groupVersion.WithKind("PartialObjectMetadata"))
      item.Object.Object = partial
    case metav1.IncludeNone:
      item.Object.Object = nil
    default:
      err = errors.NewBadRequest(fmt.Sprintf("unrecognized includeObject value: %q", opts.IncludeObject))
      return nil, err
    }
  }

  return table, nil
}

func asPartialObjectMetadata(result runtime.Object, groupVersion schema.GroupVersion) (runtime.Object, error) {
  if meta.IsListType(result) {
    err := newNotAcceptableError(fmt.Sprintf("you requested PartialObjectMetadata, but the requested object is a list (%T)", result))
    return nil, err
  }
  switch groupVersion {
  case metav1.SchemeGroupVersion:
  default:
    return nil, newNotAcceptableError(fmt.Sprintf("no PartialObjectMetadataList exists in group version %s", groupVersion))
  }
  m, err := meta.Accessor(result)
  if err != nil {
    return nil, err
  }
  partial := meta.AsPartialObjectMetadata(m)
  partial.GetObjectKind().SetGroupVersionKind(groupVersion.WithKind("PartialObjectMetadata"))
  return partial, nil
}

func asPartialObjectMetadataList(result runtime.Object, groupVersion schema.GroupVersion) (runtime.Object, error) {
  li, ok := result.(metav1.ListInterface)
  if !ok {
    return nil, newNotAcceptableError(fmt.Sprintf("you requested PartialObjectMetadataList, but the requested object is not a list (%T)", result))
  }

  gvk := groupVersion.WithKind("PartialObjectMetadata")
  switch {
  case groupVersion == metav1.SchemeGroupVersion:
    list := &metav1.PartialObjectMetadataList{}
    err := meta.EachListItem(result, func(obj runtime.Object) error {
      m, err := meta.Accessor(obj)
      if err != nil {
        return err
      }
      partial := meta.AsPartialObjectMetadata(m)
      partial.GetObjectKind().SetGroupVersionKind(gvk)
      list.Items = append(list.Items, *partial)
      return nil
    })
    if err != nil {
      return nil, err
    }
    list.ResourceVersion = li.GetResourceVersion()
    list.Continue = li.GetContinue()
    return list, nil
  default:
    return nil, newNotAcceptableError(fmt.Sprintf("no PartialObjectMetadataList exists in group version %s", groupVersion))
  }
}

// errNotAcceptable indicates Accept negotiation has failed
type errNotAcceptable struct {
  message string
}

func newNotAcceptableError(message string) error {
  return errNotAcceptable{message}
}

func (e errNotAcceptable) Error() string {
  return e.message
}

func (e errNotAcceptable) Status() metav1.Status {
  return metav1.Status{
    Status:  metav1.StatusFailure,
    Code:    http.StatusNotAcceptable,
    Reason:  metav1.StatusReason("NotAcceptable"),
    Message: e.Error(),
  }
}
