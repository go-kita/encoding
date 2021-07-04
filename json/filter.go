package json

import (
	"context"
	"encoding/json"

	"github.com/go-kita/encoding"
)

// DecoderOption is a function which modifies a *json.Decoder.
type DecoderOption func(decoder *json.Decoder)

// decoderOptionKey is the context.Context key for storing/extracting DecoderOption.
type decoderOptionKey struct {
}

// contextWithDecoderOption wraps DecoderOption into a new context.Context.
func contextWithDecoderOption(ctx context.Context, opt ...DecoderOption) context.Context {
	return context.WithValue(ctx, decoderOptionKey{}, opt)
}

// decoderOptionFromContext extracts DecoderOption from a context.Context.
func decoderOptionFromContext(ctx context.Context) []DecoderOption {
	if opt, ok := ctx.Value(decoderOptionKey{}).([]DecoderOption); ok {
		return opt
	}
	return nil
}

type optUnmarshaler struct {
	opt         []DecoderOption
	unmarshaler encoding.Unmarshaler
}

func (o *optUnmarshaler) Unmarshal(ctx context.Context, data []byte, v interface{}) error {
	ctx = contextWithDecoderOption(ctx, o.opt...)
	return o.unmarshaler.Unmarshal(ctx, data, v)
}

// WithDecoderOption returns a proxy of encoding.Unmarshaler.
// It wraps DecoderOption into context.Context, and then calls Unmarshal method
// of the underlying encoding.Unmarshaler with the new context.Context.
func WithDecoderOption(unmarshaler encoding.Unmarshaler, opt ...DecoderOption) encoding.Unmarshaler {
	return &optUnmarshaler{
		opt:         opt,
		unmarshaler: unmarshaler,
	}
}

// DisallowUnknownFields produces a DecoderOption which modifies a json.Decoder disallow unknown fields.
func DisallowUnknownFields() DecoderOption {
	return func(decoder *json.Decoder) {
		decoder.DisallowUnknownFields()
	}
}

// UseNumber produces a DecoderOption which modifies a json.Decoder decoding number to json.Number.
func UseNumber() DecoderOption {
	return func(decoder *json.Decoder) {
		decoder.UseNumber()
	}
}

// EncoderOption is a function which modifies a *.json.Encoder
type EncoderOption func(encoder *json.Encoder)

// encoderOptionKey is the context.Context key for storing/extracting EncoderOption
type encoderOptionKey struct {
}

// contextWithEncoderOption wraps EncoderOption into a new context.Context
func contextWithEncoderOption(ctx context.Context, opt ...EncoderOption) context.Context {
	return context.WithValue(ctx, encoderOptionKey{}, opt)
}

// encoderOptionFromContext extracts EncoderOption from a context.Context.
func encoderOptionFromContext(ctx context.Context) []EncoderOption {
	if opt, ok := ctx.Value(encoderOptionKey{}).([]EncoderOption); ok {
		return opt
	}
	return nil
}

type optMarshaler struct {
	opt       []EncoderOption
	marshaler encoding.Marshaler
}

func (o *optMarshaler) Marshal(ctx context.Context, v interface{}) ([]byte, error) {
	ctx = contextWithEncoderOption(ctx, o.opt...)
	return o.marshaler.Marshal(ctx, v)
}

// WithEncoderOption returns a proxy of encoding.Marshaler.
// It wraps EncoderOption into context.Context, and then calls Marshal method
// of the underlying encoding.Marshaler with the new context.Context.
func WithEncoderOption(marshaler encoding.Marshaler, opt ...EncoderOption) encoding.Marshaler {
	return &optMarshaler{
		opt:       opt,
		marshaler: marshaler,
	}
}

// EscapeHTML produces an EncoderOption which modifies a json.Encoder escaping HTML.
func EscapeHTML(on bool) EncoderOption {
	return func(encoder *json.Encoder) {
		encoder.SetEscapeHTML(on)
	}
}

// Indent produces an EncoderOption which set prefix and indent of a json.Encoder.
func Indent(prefix, indent string) EncoderOption {
	return func(encoder *json.Encoder) {
		encoder.SetIndent(prefix, indent)
	}
}
