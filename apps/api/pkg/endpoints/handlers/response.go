package handlers

import (
	"context"
	"encoding/hex"
	"fmt"
	metainternalscheme "github.com/nrc-no/core/apps/api/pkg/apis/meta/internalversion/scheme"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/negotiation"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/writers"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/util/exceptions"
	"net/http"
)

// transformResponseObject takes an object loaded from storage and performs any necessary transformations.
// Will write the complete response object.
func transformResponseObject(ctx context.Context, scope *RequestScope, req *http.Request, w http.ResponseWriter, statusCode int, mediaType negotiation.MediaTypeOptions, result runtime.Object) {
	options, err := optionsForTransform(mediaType, req)
	if err != nil {
		scope.Error(err, w, req)
		return
	}
	obj, err := transformObject(ctx, result, options, mediaType, scope, req)
	if err != nil {
		scope.Error(err, w, req)
		return
	}
	kind, serializer, _ := targetEncodingForTransform(scope, mediaType, req)
	writers.WriteObjectNegotiated(serializer, scope, kind.GroupVersion(), w, req, statusCode, obj)
}

// optionsForTransform will load and validate any additional query parameter options for
// a conversion or return an error.
func optionsForTransform(mediaType negotiation.MediaTypeOptions, req *http.Request) (interface{}, error) {
	//switch target := mediaType.Convert; {
	//case target == nil:
	//case target.Kind == "Table" && (target.GroupVersion() == metav1.SchemeGroupVersion):
	//  opts := &metav1.TableOptions{}
	//  if err := metainternalversionscheme.ParameterCodec.DecodeParameters(req.URL.Query(), metav1.SchemeGroupVersion, opts); err != nil {
	//    return nil, err
	//  }
	//  switch errs := validation.ValidateTableOptions(opts); len(errs) {
	//  case 0:
	//    return opts, nil
	//  case 1:
	//    return nil, errors.NewBadRequest(fmt.Sprintf("Unable to convert to Table as requested: %v", errs[0].Error()))
	//  default:
	//    return nil, errors.NewBadRequest(fmt.Sprintf("Unable to convert to Table as requested: %v", errs))
	//  }
	//}
	return nil, nil
}

// transformObject takes the object as returned by storage and ensures it is in
// the client's desired form, as well as ensuring any API level fields like self-link
// are properly set.
func transformObject(ctx context.Context, obj runtime.Object, opts interface{}, mediaType negotiation.MediaTypeOptions, scope *RequestScope, req *http.Request) (runtime.Object, error) {
	//if co, ok := obj.(runtime.CacheableObject); ok {
	//  if mediaType.Convert != nil {
	//    // Non-nil mediaType.Convert means that some conversion of the object
	//    // has to happen. Currently conversion may potentially modify the
	//    // object or assume something about it (e.g. asTable operates on
	//    // reflection, which won't work for any wrapper).
	//    // To ensure it will work correctly, let's operate on base objects
	//    // and not cache it for now.
	//    //
	//    // TODO: Long-term, transformObject should be changed so that it
	//    // implements runtime.Encoder interface.
	//    return doTransformObject(ctx, co.GetObject(), opts, mediaType, scope, req)
	//  }
	//}
	return doTransformObject(ctx, obj, opts, mediaType, scope, req)
}

func doTransformObject(ctx context.Context, obj runtime.Object, opts interface{}, mediaType negotiation.MediaTypeOptions, scope *RequestScope, req *http.Request) (runtime.Object, error) {
	if _, ok := obj.(*metav1.Status); ok {
		return obj, nil
	}
	// if err := setObjectSelfLink(ctx, obj, req, scope.Namer); err != nil {
	//  return nil, err
	//}

	switch target := mediaType.Convert; {
	case target == nil:
		return obj, nil

	//case target.Kind == "PartialObjectMetadata":
	//  return asPartialObjectMetadata(obj, target.GroupVersion())
	//
	//case target.Kind == "PartialObjectMetadataList":
	//  return asPartialObjectMetadataList(obj, target.GroupVersion())
	//
	//case target.Kind == "Table":
	//  options, ok := opts.(*metav1.TableOptions)
	//  if !ok {
	//    return nil, fmt.Errorf("unexpected TableOptions, got %T", opts)
	//  }
	//  return asTable(ctx, obj, options, scope, target.GroupVersion())

	default:
		accepted, _ := negotiation.MediaTypesForSerializer(metainternalscheme.Codecs)
		err := negotiation.NewNotAcceptableError(accepted)
		return nil, err
	}
}

// targetEncodingForTransform returns the appropriate serializer for the input media type
func targetEncodingForTransform(scope *RequestScope, mediaType negotiation.MediaTypeOptions, req *http.Request) (schema.GroupVersionKind, runtime.NegotiatedSerializer, bool) {
	//switch target := mediaType.Convert; {
	//case target == nil:
	//case (target.Kind == "PartialObjectMetadata" || target.Kind == "PartialObjectMetadataList" || target.Kind == "Table") &&
	//  (target.GroupVersion() == metav1beta1.SchemeGroupVersion || target.GroupVersion() == metav1.SchemeGroupVersion):
	//  return *target, metainternalversionscheme.Codecs, true
	//}
	return scope.Kind, scope.Serializer, false
}

// transformDecodeError adds additional information into a bad-request api error when a decode fails.
func transformDecodeError(typer runtime.ObjectTyper, baseErr error, into runtime.Object, gvk *schema.GroupVersionKind, body []byte) error {
	objGVKs, _, err := typer.ObjectKinds(into)
	if err != nil {
		return exceptions.NewBadRequest(err.Error())
	}
	objGVK := objGVKs[0]
	if gvk != nil && len(gvk.Kind) > 0 {
		return exceptions.NewBadRequest(fmt.Sprintf("%s in version %q cannot be handled as a %s: %v", gvk.Kind, gvk.Version, objGVK.Kind, baseErr))
	}
	summary := summarizeData(body, 30)
	return exceptions.NewBadRequest(fmt.Sprintf("the object provided is unrecognized (must be of type %s): %v (%s)", objGVK.Kind, baseErr, summary))
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
