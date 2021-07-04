package json

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/go-kita/encoding"
)

type asciiSafe struct {
	marshaler encoding.Marshaler
	pool      *sync.Pool
}

func (a *asciiSafe) Marshal(ctx context.Context, v interface{}) ([]byte, error) {
	data, err := a.marshaler.Marshal(ctx, v)
	if err != nil {
		return nil, err
	}
	buf := a.pool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		a.pool.Put(buf)
	}()
	reader := bytes.NewReader(data)
	for {
		ch, _, err := reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return buf.Bytes(), nil
			}
			return nil, err
		}
		if ch < '\u007F' || ch > '\uFFFF' {
			buf.WriteRune(ch)
			continue
		}
		_, _ = fmt.Fprintf(buf, "\\u%X", ch)
	}
}

// AsciiSafe wraps an encoding.Marshaler, the wrapper returned will do Ascii-SafetyL:
// escape all rune to Unicode expression if a rune is not a ASCII character.
//
// For example: replace char '李', bytes [0xE6, 0x9D, 0x8E] to '\u674E',
// bytes [0x5C, 0x75, 0x36, 0x37, 0x34, 0x45]
func AsciiSafe(marshaler encoding.Marshaler) encoding.Marshaler {
	return &asciiSafe{
		marshaler: marshaler,
		pool: &sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
}
