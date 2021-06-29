package xml

import (
	"context"
	"encoding/xml"
	"github.com/go-kita/encoding"
	"reflect"
	"testing"
)

type val struct {
	ID   int    `xml:"id"`
	Name string `xml:"name"`
}

type named struct {
	XMLName xml.Name `xml:"byName"`
	Name    string   `xml:"name"`
}

func Test_codec_Marshal(t *testing.T) {
	type args struct {
		ctx context.Context
		v   interface{}
	}
	tests := []struct {
		name    string
		c       *codec
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "string",
			c:    _codec,
			args: args{
				ctx: context.Background(),
				v:   "abc",
			},
			want:    []byte(`<string>abc</string>`),
			wantErr: false,
		},
		{
			name: "val",
			c:    _codec,
			args: args{
				ctx: context.Background(),
				v: val{
					ID:   0,
					Name: "n",
				},
			},
			want:    []byte(`<val><id>0</id><name>n</name></val>`),
			wantErr: false,
		},
		{
			name: "named",
			c:    _codec,
			args: args{
				ctx: context.Background(),
				v: named{
					Name: "nn",
				},
			},
			want:    []byte(`<byName><name>nn</name></byName>`),
			wantErr: false,
		},
		{
			name: "func",
			c:    _codec,
			args: args{
				ctx: context.Background(),
				v:   Register,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Marshal(tt.args.ctx, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("codec.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("codec.Marshal() = %v(%s), want %v(%s)",
					got, got, tt.want, tt.want)
			}
		})
	}
}

func Test_codec_Unmarshal(t *testing.T) {
	type args struct {
		ctx  context.Context
		data []byte
		v    interface{}
	}
	strV := ""
	strW := "str"
	valV := val{}
	valW := val{
		ID:   1,
		Name: "n",
	}
	tests := []struct {
		name    string
		c       *codec
		args    args
		wantErr bool
		want    interface{}
	}{
		{
			name: "string-p",
			c:    _codec,
			args: args{
				ctx:  context.Background(),
				data: []byte(`<string>str</string>`),
				v:    &strV,
			},
			wantErr: false,
			want:    &strW,
		},
		{
			name: "string",
			c:    _codec,
			args: args{
				ctx:  context.Background(),
				data: []byte(`<string>str</string>`),
				v:    strV,
			},
			wantErr: true,
		},
		{
			name: "string1",
			c:    _codec,
			args: args{
				ctx:  context.Background(),
				data: []byte(`<string1>str</string1>`),
				v:    &strV,
			},
			wantErr: false,
			want:    &strW,
		},
		{
			name: "val-p",
			c:    _codec,
			args: args{
				ctx:  context.Background(),
				data: []byte(`<val><id>1</id><name>n</name></val>`),
				v:    &valV,
			},
			wantErr: false,
			want:    &valW,
		},
		{
			name: "val",
			c:    _codec,
			args: args{
				ctx:  context.Background(),
				data: []byte(`<val><id>1</id><name>n</name></val>`),
				v:    valV,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.Unmarshal(tt.args.ctx, tt.args.data, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("codec.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				if !reflect.DeepEqual(tt.args.v, tt.want) {
					t.Errorf("codec.Unmarshal: want = %v, got %v", tt.want, tt.args.v)
				}
			}
		})
	}
}

func TestRegister(t *testing.T) {
	m := encoding.GetMarshaler(Name)
	if m != _codec {
		t.Errorf("expect Marshaler of name %s is _codec, but not", Name)
	}
	u := encoding.GetUnmarshaler(Name)
	if u != _codec {
		t.Errorf("expect Unmarshaler of name %s is _codec, but not", Name)
	}
}
