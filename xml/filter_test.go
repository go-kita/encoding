package xml

import (
	"context"
	"encoding/xml"
	"golang.org/x/text/encoding/simplifiedchinese"
	"reflect"
	"testing"

	"github.com/go-kita/encoding"
)

func Test_contextWithDecodeOption(t *testing.T) {
	type args struct {
		ctx context.Context
		opt []DecoderOption
	}
	tests := []struct {
		name string
		args args
		want func(context.Context) bool
	}{
		{
			name: "test1",
			args: args{
				ctx: context.Background(),
				opt: []DecoderOption{
					func(decoder *xml.Decoder) {
					},
				},
			},
			want: func(ctx context.Context) bool {
				opt, ok := ctx.Value(decoderOptionKey{}).([]DecoderOption)
				if !ok {
					t.Logf("expect contains, but not")
					return false
				}
				if len(opt) == 0 {
					t.Logf("expect not empty, but empty")
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contextWithDecoderOption(tt.args.ctx, tt.args.opt...); !tt.want(got) {
				t.Fail()
			}
		})
	}
}

func Test_decoderOptionFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want []DecoderOption
	}{
		{
			name: "no-opts",
			args: args{
				ctx: context.Background(),
			},
			want: nil,
		},
		{
			name: "empty-opts",
			args: args{
				ctx: context.WithValue(context.Background(), decoderOptionKey{}, []DecoderOption{}),
			},
			want: []DecoderOption{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decoderOptionFromContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decoderOptionFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_optUnmarshaler_Unmarshal(t *testing.T) {
	type args struct {
		ctx  context.Context
		data []byte
		v    interface{}
	}
	nop := func(decoder *xml.Decoder) {
	}
	strV := ""
	tests := []struct {
		name    string
		o       *optUnmarshaler
		args    args
		wantErr bool
	}{
		{
			name: "nop",
			o: &optUnmarshaler{
				opt:         []DecoderOption{nop},
				unmarshaler: _codec,
			},
			args: args{
				ctx:  context.Background(),
				data: []byte("<string>str</string>"),
				v:    &strV,
			},
			wantErr: false,
		},
		{
			name: "asUtf8",
			o: &optUnmarshaler{
				opt:         []DecoderOption{WithCharsetReader(AsUtf8CharsetReader())},
				unmarshaler: _codec,
			},
			args: args{
				ctx:  context.Background(),
				data: []byte("<string>str</string>"),
				v:    &strV,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.o.Unmarshal(tt.args.ctx, tt.args.data, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("optUnmarshaler.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWithDecoderOption(t *testing.T) {
	gbkData, err := simplifiedchinese.GBK.NewEncoder().Bytes(
		[]byte(`<?xml version="1.0" encoding="GBK" ?><string>中文</string>`))
	if err != nil {
		t.Fatal(err)
	}
	u := WithDecoderOption(_codec, WithCharsetReader(IanaTransformCharsetReader()))
	strV := ""
	strW := "中文"
	err = u.Unmarshal(context.Background(), gbkData, &strV)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	if strV != strW {
		t.Errorf("expect %s, got %s", strW, strV)
	}
}

func Test_contextWithEncoderOption(t *testing.T) {
	type args struct {
		ctx context.Context
		opt []EncoderOption
	}
	tests := []struct {
		name string
		args args
		want func(context.Context) bool
	}{
		{
			name: "default",
			args: args{
				ctx: context.Background(),
				opt: []EncoderOption{
					func(encoder *xml.Encoder) {
					},
				},
			},
			want: func(ctx context.Context) bool {
				opt, ok := ctx.Value(encoderOptionKey{}).([]EncoderOption)
				if !ok {
					t.Logf("expect contains, but not")
					return false
				}
				if len(opt) == 0 {
					t.Logf("expect not empty, but empty")
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contextWithEncoderOption(tt.args.ctx, tt.args.opt...); !tt.want(got) {
				t.Fail()
			}
		})
	}
}

func TestWithEncoderOption(t *testing.T) {
	type args struct {
		marshaler encoding.Marshaler
		opt       []EncoderOption
		v         interface{}
	}
	type val struct {
		V string `xml:"v"`
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "gbk",
			args: args{
				marshaler: _codec,
				opt:       []EncoderOption{WithEncodingProcInst("GBK")},
				v:         "v",
			},
			want: []byte(`<?xml version="1.0" encoding="GBK"?>
<string>v</string>`),
		},
		{
			name: "indent",
			args: args{
				marshaler: _codec,
				opt:       []EncoderOption{WithIndent("", "\t")},
				v:         val{V: "v"},
			},
			want: []byte("<val>\n\t<v>v</v>\n</val>"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithEncoderOption(tt.args.marshaler, tt.args.opt...)
			data, err := got.Marshal(context.Background(), tt.args.v)
			if err != nil {
				t.Errorf("want no error, got %v", err)
			}
			if !reflect.DeepEqual(data, tt.want) {
				t.Errorf("want %s, got %s", tt.want, data)
			}
		})
	}
}
