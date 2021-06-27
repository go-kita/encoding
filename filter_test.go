package encoding

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"golang.org/x/text/encoding/ianaindex"
)

var _ Marshaler = (*errMarshaler)(nil)

type errMarshaler struct {
}

func (e *errMarshaler) Marshal(_ context.Context, _ interface{}) ([]byte, error) {
	return nil, errors.New("encoding: test error")
}

func TestFilter_Marshal(t *testing.T) {
	stars := []byte("***")
	f := &filter{
		fn: func(pre []byte) ([]byte, error) {
			return stars, nil
		},
		marshaler: &nopMarshaler{},
	}
	data, err := f.Marshal(context.Background(), "abc")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	if !bytes.Equal(data, stars) {
		t.Errorf("expect %s, got %s", string(stars), string(data))
	}
	f.marshaler = &errMarshaler{}
	_, err = f.Marshal(context.Background(), "abc")
	if err == nil {
		t.Errorf("expect an error, got nil")
	}
}

func TestFilter_Unmarshal(t *testing.T) {
	starts := []byte("***")
	f := &filter{
		fn: func(pre []byte) ([]byte, error) {
			return starts, nil
		},
		unmarshaler: &nopUnmarshaler{},
	}
	str := "xyz"
	err := f.Unmarshal(context.Background(), []byte("abc"), &str)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	if str != "xyz" {
		t.Errorf("expect xyz, got %s", str)
	}
	f.fn = func(pre []byte) ([]byte, error) {
		return nil, errors.New("encoding: test error")
	}
	err = f.Unmarshal(context.Background(), []byte("abc"), &str)
	if err == nil {
		t.Errorf("expect an error, got nil")
	}
}

func TestFilterMarshaler(t *testing.T) {
	marshaler := FilterMarshaler(nopMarshaler{},
		func(pre []byte) ([]byte, error) {
			return []byte("***"), nil
		},
		func(pre []byte) ([]byte, error) {
			return nil, errors.New("encoding: test error")
		},
	)
	if marshaler == nil {
		t.Errorf("expect a marshaler, got nil")
	}
}

func TestFilterUnmarshaler(t *testing.T) {
	unmarshaler := FilterUnmarshaler(nopUnmarshaler{},
		func(pre []byte) ([]byte, error) {
			return []byte("***"), nil
		},
		func(pre []byte) ([]byte, error) {
			return nil, errors.New("encoding: test error")
		},
	)
	if unmarshaler == nil {
		t.Errorf("expect an unmarshaler, got nil")
	}
}

func TestEncodeWith(t *testing.T) {
	encoding, _ := ianaindex.IANA.Encoding("ISO-8859-1")
	filterFunc := EncodeWith(encoding)
	_, err := filterFunc([]byte("中文"))
	if err == nil {
		t.Errorf("expect an error, got nil")
	}
	encoding, _ = ianaindex.IANA.Encoding("GBK")
	filterFunc = EncodeWith(encoding)
	data, err := filterFunc([]byte("中文"))
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	expect := []byte{0xD6, 0xD0, 0xCE, 0xC4}
	if !bytes.Equal(data, expect) {
		t.Errorf("expect %#X, got %#X", expect, data)
	}
	filterFunc = EncodeWith(nil)
	data, err = filterFunc([]byte("中文"))
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	expect = []byte{0xE4, 0xB8, 0xAD, 0xE6, 0x96, 0x87}
	if !bytes.Equal(data, expect) {
		t.Errorf("expect %#X, got %#X", expect, data)
	}
}

func TestEncodeWithCharset(t *testing.T) {
	filterFunc := EncodeWithCharset("GBK")
	data, err := filterFunc([]byte("中文"))
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	expect := []byte{0xD6, 0xD0, 0xCE, 0xC4}
	if !bytes.Equal(data, expect) {
		t.Errorf("expect %U, got %U", expect, data)
	}
	filterFunc = EncodeWithCharset("My-Fake")
	_, err = filterFunc([]byte("中文"))
	if err == nil {
		t.Errorf("expect an error, got nil")
	}
}

func TestDecodeWith(t *testing.T) {
	encoding, _ := ianaindex.IANA.Encoding("ISO-8859-1")
	filterFunc := DecodeWith(encoding)
	gbkData := []byte{0xD6, 0xD0, 0xCE, 0xC4}
	data, err := filterFunc(gbkData)
	if err != nil {
		t.Errorf("expect an error, got nil")
	}
	expect := []byte{0xC3, 0x96, 0xC3, 0x90, 0xC3, 0x8E, 0xC3, 0x84}
	if !bytes.Equal(data, expect) {
		t.Errorf("expect %#X, got %#X", expect, data)
	}
	encoding, _ = ianaindex.IANA.Encoding("GBK")
	filterFunc = DecodeWith(encoding)
	data, err = filterFunc(gbkData)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	expect = []byte{0xE4, 0xB8, 0xAD, 0xE6, 0x96, 0x87}
	if !bytes.Equal(data, expect) {
		t.Errorf("expect %#2X, got %#2X", expect, data)
	}
	filterFunc = DecodeWith(nil)
	data, err = filterFunc(gbkData)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	if !bytes.Equal(data, gbkData) {
		t.Errorf("expect %#X, got %#X", gbkData, data)
	}
}

func TestDecodingWithCharset(t *testing.T) {
	filterFunc := DecodingWithCharset("GBK")
	gbkData := []byte{0xD6, 0xD0, 0xCE, 0xC4}
	data, err := filterFunc(gbkData)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	expect := []byte("中文")
	if !bytes.Equal(data, expect) {
		t.Errorf("expect %#X, got %#X", expect, data)
	}
	filterFunc = DecodingWithCharset("MY-Fake")
	_, err = filterFunc(gbkData)
	if err == nil {
		t.Errorf("expect an error, got nil")
	}
}
