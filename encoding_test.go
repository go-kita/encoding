package encoding

import (
	"context"
	"fmt"
	"testing"
	"unsafe"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
)

func TestMustWithCharset(t *testing.T) {
	t.Run("expect_panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expect panic but not")
			}
		}()
		MustWithCharset(context.Background(), "my-charset-xxx")
	})
	t.Run("common", func(t *testing.T) {
		ctx := MustWithCharset(context.Background(), "ISO-8859-1")
		e, ok := ctx.Value(encodingKey{}).(encoding.Encoding)
		if !ok || e == nil {
			t.Errorf("expect got Encoding but not")
		}
		name, _ := ianaindex.IANA.Name(e)
		t.Logf("encoding name: %s", name)
	})
}

func TestFromContext(t *testing.T) {
	t.Run("from_background", func(t *testing.T) {
		_, err := FromContext(context.Background())
		if !IsErrNoEncoding(err) {
			t.Errorf("expect errNoEncoding, but got %q", err)
		}
	})
	t.Run("common", func(t *testing.T) {
		e, _ := ianaindex.IANA.Encoding("ISO-8859-1")
		ctx := context.WithValue(context.Background(), encodingKey{}, e)
		ee, err := FromContext(ctx)
		if err != nil {
			t.Errorf("expect no err, but got %v", err)
		}
		if ee != e {
			t.Errorf("expect %v but got %v", e, ee)
		}
	})
}

func TestIsErrNoEncoding(t *testing.T) {
	tests := []struct {
		err  error
		want bool
	}{
		{fmt.Errorf("err: abc"), false},
		{errNoEncoding, true},
		{fmt.Errorf("biz: %w", errNoEncoding), true},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := IsErrNoEncoding(test.err)
			if got != test.want {
				t.Errorf("want %t, got %t", test.want, got)
			}
		})
	}
}

var _ Marshaler = nopMarshaler{}

type nopMarshaler struct {
}

func (n nopMarshaler) Marshal(_ context.Context, _ interface{}) ([]byte, error) {
	return nil, nil
}

func TestRegisterMarshaler(t *testing.T) {
	origin := _marshalers.Swap(unsafe.Pointer(&map[string]Marshaler{}))
	defer _marshalers.Store(origin)
	_ = RegisterMarshaler("", nopMarshaler{})
	l := len(*(*map[string]Marshaler)(_marshalers.Load()))
	if l != 0 {
		t.Errorf("expect size 0, got %d", l)
	}
	_ = RegisterMarshaler("nop", nil)
	l = len(*(*map[string]Marshaler)(_marshalers.Load()))
	if l != 0 {
		t.Errorf("expect size 0, got %d", l)
	}
	l1 := len(*(*map[string]Marshaler)(_marshalers.Load()))
	m := RegisterMarshaler("nop", nopMarshaler{})
	if m != nil {
		t.Errorf("expect nil, but got %q", m)
	}
	l2 := len(*(*map[string]Marshaler)(_marshalers.Load()))
	if l2 <= l1 {
		t.Errorf("expect size raised, but not, before %d, after %d", l1, l2)
	}
	m = RegisterMarshaler("nop", nopMarshaler{})
	if m == nil {
		t.Errorf("expect not nil, but got nil")
	}
	l3 := len(*(*map[string]Marshaler)(_marshalers.Load()))
	if l3 != l2 {
		t.Errorf("expect size unchanged, but not, before %d, after %d", l2, l3)
	}
}

func TestGetMarshaler(t *testing.T) {
	origin := _marshalers.Swap(unsafe.Pointer(&map[string]Marshaler{}))
	defer _marshalers.Store(origin)
	m := GetMarshaler("nop")
	if m != nil {
		t.Errorf("expect nil, got %q", m)
	}
	_ = RegisterMarshaler("nop", nopMarshaler{})
	m = GetMarshaler("nop")
	if m == nil {
		t.Errorf("expect not nil, got nil")
	}
}

var _ Unmarshaler = nopUnmarshaler{}

type nopUnmarshaler struct {
}

func (n nopUnmarshaler) Unmarshal(_ context.Context, _ []byte, _ interface{}) error {
	return nil
}

func TestRegisterUnmarshaler(t *testing.T) {
	origin := _unmarshalers.Swap(unsafe.Pointer(&map[string]Unmarshaler{}))
	defer _unmarshalers.Store(origin)
	_ = RegisterUnmarshaler("", nopUnmarshaler{})
	l := len(*(*map[string]Unmarshaler)(_unmarshalers.Load()))
	if l != 0 {
		t.Errorf("expect size 0, got %d", l)
	}
	_ = RegisterUnmarshaler("nop", nil)
	l = len(*(*map[string]Unmarshaler)(_unmarshalers.Load()))
	if l != 0 {
		t.Errorf("expect size 0, got %d", l)
	}
	l1 := len(*(*map[string]Unmarshaler)(_unmarshalers.Load()))
	u := RegisterUnmarshaler("nop", nopUnmarshaler{})
	if u != nil {
		t.Errorf("expect nil, but got %q", u)
	}
	l2 := len(*(*map[string]Unmarshaler)(_unmarshalers.Load()))
	if l2 <= l1 {
		t.Errorf("expect size raised, but not, before %d, after %d", l1, l2)
	}
	u = RegisterUnmarshaler("nop", nopUnmarshaler{})
	if u == nil {
		t.Errorf("expect not nil, but got nil")
	}
	l3 := len(*(*map[string]Unmarshaler)(_unmarshalers.Load()))
	if l3 != l2 {
		t.Errorf("expect size unchanged, but not, before %d, after %d", l2, l3)
	}
}

func TestGetUnmarshaler(t *testing.T) {
	origin := _unmarshalers.Swap(unsafe.Pointer(&map[string]Unmarshaler{}))
	defer _unmarshalers.Store(origin)
	u := GetUnmarshaler("nop")
	if u != nil {
		t.Errorf("expect nil, got %q", u)
	}
	_ = RegisterUnmarshaler("nop", nopUnmarshaler{})
	u = GetUnmarshaler("nop")
	if u == nil {
		t.Errorf("expect not nil, got nil")
	}
}
