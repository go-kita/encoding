package json

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-kita/encoding"
)

func TestWithDecoderOption(t *testing.T) {
	type args struct {
		unmarshaler encoding.Unmarshaler
		opt         []DecoderOption
	}
	type testStruct struct {
		V string
	}
	tests := []struct {
		name string
		args args
		want func(u encoding.Unmarshaler) bool
	}{
		{
			name: "disableUnknownFields",
			args: args{
				unmarshaler: &codec{buf: _bufPool},
				opt:         []DecoderOption{DisallowUnknownFields()},
			},
			want: func(u encoding.Unmarshaler) bool {
				err := u.Unmarshal(context.Background(), []byte(`{"k":"v"}`), &testStruct{})
				if err == nil {
					t.Logf("expect err, got nil")
					return false
				}
				return true
			},
		},
		{
			name: "useNumber",
			args: args{
				unmarshaler: &codec{buf: _bufPool},
				opt:         []DecoderOption{UseNumber()},
			},
			want: func(u encoding.Unmarshaler) bool {
				var num json.Number
				err := u.Unmarshal(context.Background(), []byte(`100`), &num)
				if err != nil {
					t.Logf("expect nil err, got %v", err)
					return false
				}
				t.Logf("%v", num)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDecoderOption(tt.args.unmarshaler, tt.args.opt...); !tt.want(got) {
				t.Fail()
			}
		})
	}
}

func TestWithEncoderOption(t *testing.T) {
	type args struct {
		marshaler encoding.Marshaler
		opt       []EncoderOption
	}
	type testStruct struct {
		V string
	}
	tests := []struct {
		name string
		args args
		want func(encoding.Marshaler) bool
	}{
		{
			name: "escapeHtml",
			args: args{
				marshaler: &codec{buf: _bufPool},
				opt:       []EncoderOption{EscapeHTML(true)},
			},
			want: func(marshaler encoding.Marshaler) bool {
				data, err := marshaler.Marshal(context.Background(), "<i>en</i>")
				if err != nil {
					t.Logf("expect no err, got %v", err)
					return false
				}
				expect := []byte("\"\\u003ci\\u003een\\u003c/i\\u003e\"\n")
				if !reflect.DeepEqual(data, expect) {
					t.Logf("\nexpect %v(%s),\ngot    %v(%s)", expect, expect, data, data)
					return false
				}
				return true
			},
		},
		{
			name: "indent",
			args: args{
				marshaler: &codec{buf: _bufPool},
				opt:       []EncoderOption{Indent("", "\t")},
			},
			want: func(marshaler encoding.Marshaler) bool {
				v := &testStruct{V: "v"}
				data, err := marshaler.Marshal(context.Background(), v)
				if err != nil {
					t.Logf("expect no err, got %v", err)
					return false
				}
				expect := []byte("{\n\t\"V\": \"v\"\n}\n")
				if !reflect.DeepEqual(data, expect) {
					t.Logf("\nexpect %v(%s),\ngot    %v(%s)", expect, expect, data, data)
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithEncoderOption(tt.args.marshaler, tt.args.opt...); !tt.want(got) {
				t.Fail()
			}
		})
	}
}
