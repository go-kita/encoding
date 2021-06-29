package xml

import (
	"bytes"
	"context"
	"encoding/xml"
	"sync"

	"github.com/go-kita/encoding"
)

func init() {
	Register(Name)
}

// Name is type name.
const Name = "xml"

var _ encoding.Marshaler = (*codec)(nil)
var _ encoding.Unmarshaler = (*codec)(nil)

type codec struct {
	buf *sync.Pool
}

var _codec = &codec{
	buf: &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	},
}

func (c *codec) Marshal(ctx context.Context, v interface{}) ([]byte, error) {
	buf := c.buf.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		c.buf.Put(buf)
	}()
	encoder := xml.NewEncoder(buf)
	for _, option := range encoderOptionFromContext(ctx) {
		option(encoder)
	}
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *codec) Unmarshal(ctx context.Context, data []byte, v interface{}) error {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = IanaTransformCharsetReader()
	for _, option := range decoderOptionFromContext(ctx) {
		option(decoder)
	}
	return decoder.Decode(v)
}

// Register register marshaler/unmarshaler.
func Register(name string) {
	encoding.RegisterMarshaler(name, _codec)
	encoding.RegisterUnmarshaler(name, _codec)
}
