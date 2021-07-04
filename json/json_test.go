package json

import (
	"context"
	"reflect"
	"testing"
)

func Test_codec_Marshal(t *testing.T) {
	type args struct {
		ctx context.Context
		v   interface{}
	}
	c := &codec{buf: _bufPool}
	tests := []struct {
		name    string
		c       *codec
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "string",
			c:    c,
			args: args{
				ctx: context.Background(),
				v:   "abc",
			},
			want:    []byte("\"abc\"\n"),
			wantErr: false,
		},
		{
			name: "struct",
			c:    c,
			args: args{
				ctx: context.Background(),
				v: struct {
					Key string `json:"key"`
					V   string `json:"vv"`
				}{
					Key: "x",
					V:   "y",
				},
			},
			want:    []byte("{\"key\":\"x\",\"vv\":\"y\"}\n"),
			wantErr: false,
		},
		{
			name: "func",
			c:    c,
			args: args{
				ctx: context.Background(),
				v:   func() {},
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
				t.Errorf("codec.Marshal() = %v, want %v", got, tt.want)
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
	type vv struct {
		V string `json:"v"`
	}
	strV := ""
	strW := "abc"
	c := &codec{buf: _bufPool}
	tests := []struct {
		name    string
		c       *codec
		args    args
		wantErr bool
		want    func(v interface{}) bool
	}{
		{
			name: "str",
			c:    c,
			args: args{
				ctx:  context.Background(),
				data: []byte(`"abc"`),
				v:    &strV,
			},
			wantErr: false,
			want: func(v interface{}) bool {
				if strW != strV {
					t.Logf("expect strV %q, but is %q", strW, strV)
					return false
				}
				return true
			},
		},
		{
			name: "err",
			c:    c,
			args: args{
				ctx:  context.Background(),
				data: []byte("aaa"),
				v:    &strV,
			},
			wantErr: true,
		},
		{
			name: "struct",
			c:    c,
			args: args{
				ctx:  context.Background(),
				data: []byte(`{"v":"vv"}`),
				v:    &vv{},
			},
			wantErr: false,
			want: func(v interface{}) bool {
				if v.(*vv).V != "vv" {
					t.Logf("expect %q, got %q", &vv{V: "vv"}, v)
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Unmarshal(tt.args.ctx, tt.args.data, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("codec.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if !tt.wantErr && !tt.want(tt.args.v) {
					t.Fail()
				}
			}
		})
	}
}
