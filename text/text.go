// Package text defines and registers Marshaler/Unmarshaler handling textual type.
package text

import (
	"context"
	se "encoding"
	"errors"
	"fmt"
	"reflect"

	"github.com/go-kita/encoding"
)

func init() {
	Register(Name)
}

// Name is type name.
const Name = "text"

var _ encoding.Marshaler = (*codec)(nil)
var _ encoding.Unmarshaler = (*codec)(nil)

type codec struct {
}

var _codec = &codec{}

func (s *codec) Marshal(_ context.Context, v interface{}) (data []byte, err error) {
	switch vv := v.(type) {
	case se.TextMarshaler:
		data, err = vv.MarshalText()
	case fmt.Stringer:
		data = []byte(vv.String())
	case string:
		data = []byte(vv)
	default:
		data = []byte(fmt.Sprintf("%v", vv))
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *codec) Unmarshal(_ context.Context, data []byte, v interface{}) (err error) {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			if !rv.CanAddr() {
				return errors.New("text: cannot unmarshal to unaddressable value")
			}
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	switch vv := v.(type) {
	case se.TextUnmarshaler:
		if err := vv.UnmarshalText(data); err != nil {
			return err
		}
		return nil
	case *string:
		*vv = string(data)
		return nil
	default:
		return fmt.Errorf("text: can not unmarshal type %T", v)
	}
}

// Register register marshaler/unmarshaler.
func Register(name string) {
	encoding.RegisterMarshaler(name, func() encoding.Marshaler { return &codec{} })
	encoding.RegisterUnmarshaler(name, func() encoding.Unmarshaler { return &codec{} })
}
