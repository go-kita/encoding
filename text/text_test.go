package text

import (
	"bytes"
	"context"
	se "encoding"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/go-kita/encoding"
)

func TestRegister(t *testing.T) {
	m := encoding.GetMarshaler(Name)
	if m != _codec {
		t.Errorf("expect Marshaler of name %q is _codec, but not", Name)
	}
	u := encoding.GetUnmarshaler(Name)
	if u != _codec {
		t.Errorf("expect Unmarshaler of name %q is _codec, but not", Name)
	}
}

var _ se.TextMarshaler = (*textual)(nil)
var _ se.TextUnmarshaler = (*textual)(nil)

type textual struct {
	str string
}

func (t *textual) MarshalText() (text []byte, err error) {
	if t.str == "err" {
		return nil, errors.New("text: textual error")
	}
	return []byte(t.str), nil
}

func (t *textual) UnmarshalText(text []byte) error {
	t.str = string(text)
	if t.str == "err" {
		return errors.New("text: textual error")
	}
	return nil
}

type no uint

func (n no) String() string {
	return fmt.Sprintf("No.%d", n)
}

func TestCodec_Marshal(t *testing.T) {
	bg := context.Background()
	tests := []struct {
		ctx       context.Context
		v         interface{}
		wantData  []byte
		wantError bool
	}{
		{
			ctx:       bg,
			v:         &textual{"abc"},
			wantData:  []byte("abc"),
			wantError: false,
		},
		{
			ctx:       bg,
			v:         &textual{"err"},
			wantData:  nil,
			wantError: true,
		},
		{
			ctx:       bg,
			v:         no(1),
			wantData:  []byte("No.1"),
			wantError: false,
		},
		{
			ctx:       bg,
			v:         "hello",
			wantData:  []byte("hello"),
			wantError: false,
		},
		{
			ctx:       bg,
			v:         errors.New("text: text error"),
			wantData:  []byte("text: text error"),
			wantError: false,
		},
		{
			ctx:       bg,
			v:         false,
			wantData:  []byte("false"),
			wantError: false,
		},
		{
			ctx:       bg,
			v:         _codec,
			wantData:  []byte("&{}"),
			wantError: false,
		},
		{
			ctx:       bg,
			v:         "中文",
			wantData:  []byte{0xE4, 0xB8, 0xAD, 0xE6, 0x96, 0x87},
			wantError: false,
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			data, err := _codec.Marshal(test.ctx, test.v)
			if test.wantError && err == nil {
				t.Errorf("want err but got nil")
			}
			if !test.wantError && err != nil {
				t.Errorf("do not want err but got %v", err)
			}
			if len(test.wantData) > 0 && !bytes.Equal(data, test.wantData) {
				t.Errorf("want data %b(%s), but got %b(%s)",
					test.wantData, string(test.wantData), data, string(data))
			}
		})
	}
}

func TestCodec_Unmarshal(t *testing.T) {
	bg := context.Background()
	tests := []struct {
		ctx     context.Context
		data    []byte
		v       interface{}
		wantV   interface{}
		wantErr bool
	}{
		{
			ctx:     bg,
			data:    []byte("abc"),
			v:       (*textual)(nil),
			wantV:   nil,
			wantErr: true,
		},
		{
			ctx:     bg,
			data:    []byte("abc"),
			v:       &textual{},
			wantV:   &textual{str: "abc"},
			wantErr: false,
		},
		{
			ctx:     bg,
			data:    []byte("err"),
			v:       &textual{},
			wantV:   nil,
			wantErr: true,
		},
		{
			ctx:     bg,
			data:    []byte("abc"),
			v:       textual{}, // *textual is TextUnmarshaler, but textual is not.
			wantV:   nil,
			wantErr: true,
		},
		{
			ctx:     bg,
			data:    []byte("str"),
			v:       (*string)(nil),
			wantV:   nil,
			wantErr: true,
		},
		{
			ctx:     bg,
			data:    []byte("str"),
			v:       "", // string is immutable
			wantV:   nil,
			wantErr: true,
		},
		{
			ctx:     bg,
			data:    []byte("str"),
			v:       123,
			wantV:   nil,
			wantErr: true,
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			err := _codec.Unmarshal(test.ctx, test.data, test.v)
			if test.wantErr {
				if err == nil {
					t.Errorf("want err, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("do not want err, but got %v", err)
				} else {
					if !reflect.DeepEqual(test.v, test.wantV) {
						t.Errorf("expect result %v, got %v", test.wantV, test.v)
					}
				}
			}
		})
	}
}
