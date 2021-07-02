package encoding

import (
	"context"
	"testing"
	"unsafe"
)

var _ Marshaler = nopMarshaler{}

type nopMarshaler struct {
}

func (n nopMarshaler) Marshal(_ context.Context, _ interface{}) ([]byte, error) {
	return nil, nil
}

func TestRegisterMarshaler(t *testing.T) {
	origin := _marshalerSuppliers.Swap(unsafe.Pointer(&map[string]Marshaler{}))
	defer _marshalerSuppliers.Store(origin)
	_ = RegisterMarshaler("", func() Marshaler {
		return nopMarshaler{}
	})
	l := len(*(*map[string]Marshaler)(_marshalerSuppliers.Load()))
	if l != 0 {
		t.Errorf("expect size 0, got %d", l)
	}
	_ = RegisterMarshaler("nop", nil)
	l = len(*(*map[string]Marshaler)(_marshalerSuppliers.Load()))
	if l != 0 {
		t.Errorf("expect size 0, got %d", l)
	}
	l1 := len(*(*map[string]Marshaler)(_marshalerSuppliers.Load()))
	m := RegisterMarshaler("nop", func() Marshaler {
		return nopMarshaler{}
	})
	if m != nil {
		t.Errorf("expect nil, but got a func")
	}
	l2 := len(*(*map[string]Marshaler)(_marshalerSuppliers.Load()))
	if l2 <= l1 {
		t.Errorf("expect size raised, but not, before %d, after %d", l1, l2)
	}
	m = RegisterMarshaler("nop", func() Marshaler {
		return nopMarshaler{}
	})
	if m == nil {
		t.Errorf("expect not nil, but got nil")
	}
	l3 := len(*(*map[string]Marshaler)(_marshalerSuppliers.Load()))
	if l3 != l2 {
		t.Errorf("expect size unchanged, but not, before %d, after %d", l2, l3)
	}
}

func TestGetMarshaler(t *testing.T) {
	origin := _marshalerSuppliers.Swap(unsafe.Pointer(&map[string]Marshaler{}))
	defer _marshalerSuppliers.Store(origin)
	m := GetMarshaler("nop")
	if m != nil {
		t.Errorf("expect nil, got %q", m)
	}
	_ = RegisterMarshaler("nop", func() Marshaler {
		return nopMarshaler{}
	})
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
	origin := _unmarshalerSuppliers.Swap(unsafe.Pointer(&map[string]Unmarshaler{}))
	defer _unmarshalerSuppliers.Store(origin)
	_ = RegisterUnmarshaler("", func() Unmarshaler {
		return nopUnmarshaler{}
	})
	l := len(*(*map[string]Unmarshaler)(_unmarshalerSuppliers.Load()))
	if l != 0 {
		t.Errorf("expect size 0, got %d", l)
	}
	_ = RegisterUnmarshaler("nop", nil)
	l = len(*(*map[string]Unmarshaler)(_unmarshalerSuppliers.Load()))
	if l != 0 {
		t.Errorf("expect size 0, got %d", l)
	}
	l1 := len(*(*map[string]Unmarshaler)(_unmarshalerSuppliers.Load()))
	u := RegisterUnmarshaler("nop", func() Unmarshaler {
		return nopUnmarshaler{}
	})
	if u != nil {
		t.Errorf("expect nil, but got a func")
	}
	l2 := len(*(*map[string]Unmarshaler)(_unmarshalerSuppliers.Load()))
	if l2 <= l1 {
		t.Errorf("expect size raised, but not, before %d, after %d", l1, l2)
	}
	u = RegisterUnmarshaler("nop", func() Unmarshaler {
		return nopUnmarshaler{}
	})
	if u == nil {
		t.Errorf("expect not nil, but got nil")
	}
	l3 := len(*(*map[string]Unmarshaler)(_unmarshalerSuppliers.Load()))
	if l3 != l2 {
		t.Errorf("expect size unchanged, but not, before %d, after %d", l2, l3)
	}
}

func TestGetUnmarshaler(t *testing.T) {
	origin := _unmarshalerSuppliers.Swap(unsafe.Pointer(&map[string]Unmarshaler{}))
	defer _unmarshalerSuppliers.Store(origin)
	u := GetUnmarshaler("nop")
	if u != nil {
		t.Errorf("expect nil, got %q", u)
	}
	_ = RegisterUnmarshaler("nop", func() Unmarshaler {
		return nopUnmarshaler{}
	})
	u = GetUnmarshaler("nop")
	if u == nil {
		t.Errorf("expect not nil, got nil")
	}
}
