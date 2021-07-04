package json

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-kita/encoding"
)

func TestAsciiSafe(t *testing.T) {
	type args struct {
		marshaler encoding.Marshaler
	}
	tests := []struct {
		name string
		args args
		want func(encoding.Marshaler) bool
	}{
		{
			name: "default",
			args: args{
				marshaler: &codec{buf: _bufPool},
			},
			want: func(marshaler encoding.Marshaler) bool {
				data, err := marshaler.Marshal(context.Background(), "hello, 李")
				if err != nil {
					t.Logf("expect no err, got %v", err)
					return false
				}
				expect := []byte("\"hello, \\u674E\"\n")
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
			if got := AsciiSafe(tt.args.marshaler); !tt.want(got) {
				t.Fail()
			}
		})
	}
	t.Logf("%x", []byte("\\u674E"))
}
