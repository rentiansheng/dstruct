package jsoniter

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/modern-go/reflect2"
)

var (
	defaultFrozenConfig = initFronzeConfig()
)

func initFronzeConfig() *frozenConfig {
	api := &frozenConfig{}
	api.streamPool = &sync.Pool{
		New: func() interface{} {
			return NewStream(api, nil, 512)
		},
	}
	api.iteratorPool = &sync.Pool{
		New: func() interface{} {
			return NewIterator(api)
		},
	}
	api.initCache()
	encoderExtension := EncoderExtension{}
	decoderExtension := DecoderExtension{}

	api.encoderExtension = encoderExtension
	api.decoderExtension = decoderExtension
	return api
}

func CreateDecoderOfType(typ reflect2.Type) ValDecoder {
	ctx := &ctx{
		frozenConfig: defaultFrozenConfig,
		prefix:       "",
		decoders:     map[reflect2.Type]ValDecoder{},
		encoders:     map[reflect2.Type]ValEncoder{},
	}
	decoder := createDecoderOfJsonRawMessage(ctx, typ)
	if decoder != nil {
		return decoder
	}
	decoder = createDecoderOfJsonNumber(ctx, typ)
	if decoder != nil {
		return decoder
	}
	decoder = createDecoderOfMarshaler(ctx, typ)
	if decoder != nil {
		return decoder
	}
	decoder = createDecoderOfAny(ctx, typ)
	if decoder != nil {
		return decoder
	}
	decoder = createDecoderOfNative(ctx, typ)
	if decoder != nil {
		return decoder
	}
	switch typ.Kind() {
	case reflect.Interface:
		ifaceType, isIFace := typ.(*reflect2.UnsafeIFaceType)
		if isIFace {
			return &ifaceDecoder{valType: ifaceType}
		}
		return &efaceDecoder{}
	case reflect.Struct:
		return decoderOfStruct(ctx, typ)
	case reflect.Array:
		return decoderOfArray(ctx, typ)
	case reflect.Slice:
		return decoderOfSlice(ctx, typ)
	case reflect.Map:
		return decoderOfMap(ctx, typ)
	case reflect.Ptr:
		return decoderOfOptional(ctx, typ)
	default:
		return &lazyErrorDecoder{err: fmt.Errorf("%s%s is unsupported type", ctx.prefix, typ.String())}
	}
}
