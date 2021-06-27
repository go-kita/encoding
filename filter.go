package encoding

import (
	"context"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
)

// FilterFunc is a function which handle binary data.
// Possible tasks are:
//   - Textual content charset encoding / decoding.
//   - Compression / Decompression
//   - Encryption / Decryption.
//   - Signing / Verifying.
//   - ...
type FilterFunc func(pre []byte) ([]byte, error)

type filter struct {
	fn          FilterFunc
	marshaler   Marshaler
	unmarshaler Unmarshaler
}

func (f *filter) Marshal(ctx context.Context, v interface{}) ([]byte, error) {
	data, err := f.marshaler.Marshal(ctx, v)
	if err != nil {
		return nil, err
	}
	return f.fn(data)
}

func (f *filter) Unmarshal(ctx context.Context, data []byte, v interface{}) error {
	data, err := f.fn(data)
	if err != nil {
		return err
	}
	return f.unmarshaler.Unmarshal(ctx, data, v)
}

// FilterMarshaler decorates a marshaler with FilterFunc. The binary data the
// decorated Marshaler produced will be reprocessed by the FilterFunc.
func FilterMarshaler(marshaler Marshaler, filterFunc ...FilterFunc) Marshaler {
	for _, fn := range filterFunc {
		marshaler = &filter{
			fn:        fn,
			marshaler: marshaler,
		}
	}
	return marshaler
}

// FilterUnmarshaler decorates a unmarshaler with FilterFunc. The binary data the
// decorated Unmarshaler to produce will be processed by the FilterFunc first.
func FilterUnmarshaler(unmarshaler Unmarshaler, filterFunc ...FilterFunc) Unmarshaler {
	for _, fn := range filterFunc {
		unmarshaler = &filter{
			fn:          fn,
			unmarshaler: unmarshaler,
		}
	}
	return unmarshaler
}

// EncodeWith produces a FilterFunc which re-encodes the binary data to specific encoding.
// If the encoding provided is nil, the returned FilterFunc just keep the origin binary data.
func EncodeWith(encoding encoding.Encoding) FilterFunc {
	return func(pre []byte) ([]byte, error) {
		if encoding == nil {
			return pre, nil
		}
		return encoding.NewEncoder().Bytes(pre)
	}
}

// EncodeWithCharset produces a FilterFunc which re-encodes the binary data to
// encoding of specific name. If the charset / encoding is not supported by runtime
// platform, the produced FilterFunc won't encode the data and return a non-nil error.
func EncodeWithCharset(name string) FilterFunc {
	return func(pre []byte) ([]byte, error) {
		e, err := ianaindex.IANA.Encoding(name)
		if err != nil {
			return nil, err
		}
		return e.NewEncoder().Bytes(pre)
	}
}

// DecodeWith produces a FilterFunc which decodes the binary data from specific encoding.
// If the encoding provided is nil, the returned FilterFunc just keep the origin binary data.
func DecodeWith(encoding encoding.Encoding) FilterFunc {
	return func(pre []byte) ([]byte, error) {
		if encoding == nil {
			return pre, nil
		}
		return encoding.NewDecoder().Bytes(pre)
	}
}

// DecodingWithCharset produces a FilterFunc which decodes the binary data from
// encoding of specific name. If the charset / encoding is not supported by runtime
// platform, the produced FilterFunc won't decode the data and return a non-nil error.
func DecodingWithCharset(name string) FilterFunc {
	return func(pre []byte) ([]byte, error) {
		e, err := ianaindex.IANA.Encoding(name)
		if err != nil {
			return nil, err
		}
		return e.NewDecoder().Bytes(pre)
	}
}
