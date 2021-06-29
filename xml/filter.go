package xml

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/go-kita/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
)

// DecoderOption is a function which modifies a *xml.Decoder.
type DecoderOption func(decoder *xml.Decoder)

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

// CharsetReader is function for xml.Decoder.CharsetReader.
type CharsetReader func(charset string, input io.Reader) (io.Reader, error)

// WithCharsetReader returns a DecoderOption which modifies decoder's CharsetReader.
func WithCharsetReader(reader CharsetReader) DecoderOption {
	return func(decoder *xml.Decoder) {
		decoder.CharsetReader = reader
	}
}

// AsUtf8CharsetReader returns a CharsetReader which treat the input content as UTF-8 encoded,
// no matter what charset the input claims.
func AsUtf8CharsetReader() CharsetReader {
	return func(charset string, input io.Reader) (io.Reader, error) {
		return input, nil
	}
}

// IanaTransformCharsetReader return a CharsetReader which treat the input content
// as encoded as what charset it claims. The CharsetReader will look up supported
// encoding according to the charset name. If no charset found or the encoding is
// not supported by runtime platform, The CharsetReader returns a non-nil error.
// The CharsetReader decode the input content by the encoding found.
func IanaTransformCharsetReader() CharsetReader {
	return func(charset string, input io.Reader) (io.Reader, error) {
		e, err := ianaindex.IANA.Encoding(charset)
		if err != nil {
			return nil, err
		}
		if e == nil {
			return nil, fmt.Errorf("xml: unsupported encoding for charset: %s", charset)
		}
		return transform.NewReader(input, e.NewDecoder()), nil
	}
}

// EncoderOption is a function which modifies a *xml.Encoder
type EncoderOption func(encoder *xml.Encoder)

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

// WithEncoderOption returns a proxy of encoding.Marshaler
// It wraps EncoderOption into context.Context, and then calls Marshal method
// of the underlying encoding.Marshaler with the new context.Context.
func WithEncoderOption(marshaler encoding.Marshaler, opt ...EncoderOption) encoding.Marshaler {
	return &optMarshaler{
		opt:       opt,
		marshaler: marshaler,
	}
}

// WithEncodingProcInst returns a EncoderOption which writes a processing instruction.
// The content of processing instruction will be `<?xml version="1.0" encoding="#{charset-name}"?>`.
// NOTE: You should provide this option at most once. Only the first option this function produced will work.
func WithEncodingProcInst(name string) EncoderOption {
	return func(encoder *xml.Encoder) {
		_ = encoder.EncodeToken(xml.ProcInst{
			Target: "xml",
			Inst:   []byte(fmt.Sprintf(`version="1.0" encoding="%s"`, name)),
		})
		_ = encoder.EncodeToken(xml.CharData{'\n'})
	}
}

// WithIndent returns a EncoderOption which modifies the indent the encoder outputs.
func WithIndent(prefix string, indent string) EncoderOption {
	return func(encoder *xml.Encoder) {
		encoder.Indent(prefix, indent)
	}
}
