package storage

import (
	"fmt"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/runtime/serializer/recognizer"
	"github.com/nrc-no/core/apps/api/pkg/storage/backend"
	"mime"
)

type StorageCodecConfig struct {
	StorageMediaType  string
	StorageSerializer runtime.StorageSerializer
	StorageVersion    schema.GroupVersion
	MemoryVersion     schema.GroupVersion
	Config            backend.Config

	EncoderDecoratorFn func(runtime.Encoder) runtime.Encoder
	DecoderDecoratorFn func([]runtime.Decoder) []runtime.Decoder
}

func NewStorageCodec(opts StorageCodecConfig) (runtime.Codec, runtime.GroupVersioner, error) {

	mediaType, _, err := mime.ParseMediaType(opts.StorageMediaType)
	if err != nil {
		return nil, nil, fmt.Errorf("%s is not a valid mime-type", mediaType)
	}

	serializer, ok := runtime.SerializerInfoForMediaType(opts.StorageSerializer.SupportedMediaTypes(), mediaType)
	if !ok {
		return nil, nil, fmt.Errorf("unable to find serializer for media-type %s", mediaType)
	}
	s := serializer.Serializer

	var encoder runtime.Encoder = s
	if opts.EncoderDecoratorFn != nil {
		encoder = opts.EncoderDecoratorFn(encoder)
	}

	decoders := []runtime.Decoder{
		s,
		opts.StorageSerializer.UniversalDeserializer(),
	}
	if opts.DecoderDecoratorFn != nil {
		decoders = opts.DecoderDecoratorFn(decoders)
	}

	encodeVersioner := runtime.NewMultiGroupVersioner(
		opts.StorageVersion,
		schema.GroupKind{Group: opts.StorageVersion.Group},
		schema.GroupKind{Group: opts.MemoryVersion.Group},
	)

	encoder = opts.StorageSerializer.EncoderForVersion(
		encoder,
		encodeVersioner,
	)
	decoder := opts.StorageSerializer.DecoderToVersion(
		recognizer.NewDecoder(decoders...),
		runtime.NewCoercingMultiGroupVersioner(
			opts.MemoryVersion,
			schema.GroupKind{Group: opts.MemoryVersion.Group},
			schema.GroupKind{Group: opts.StorageVersion.Group},
		),
	)

	return runtime.NewCodec(encoder, decoder), encodeVersioner, nil

}
