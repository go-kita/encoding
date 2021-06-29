package xml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func TestXmlCharsetReader(t *testing.T) {
	gbkBytes, err := ioutil.ReadFile("test_data/gbk.xml")
	if err != nil {
		t.Fatal(err)
	}
	utf8Bytes, err := simplifiedchinese.GBK.NewDecoder().Bytes(gbkBytes)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		data          []byte
		charsetReader CharsetReader
		wantMatch     bool
	}{
		{gbkBytes, IanaTransformCharsetReader(), true},
		{utf8Bytes, IanaTransformCharsetReader(), false},
		{utf8Bytes, AsUtf8CharsetReader(), true},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			decoder := xml.NewDecoder(bytes.NewReader(test.data))
			decoder.CharsetReader = test.charsetReader
			for {
				token, err := decoder.Token()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					t.Fatal(err)
				}
				switch tt := token.(type) {
				case xml.ProcInst:
					t.Logf("%s", token)
				case xml.CharData:
					text := bytes.TrimSpace(tt)
					if len(text) == 0 {
						break
					}
					if string(text) == "key1" {
						break
					}
					t.Logf("text: %s", text)
					match := string(text) == "值1"
					if match != test.wantMatch {
						if test.wantMatch {
							t.Errorf("expect match, but not")
						} else {
							t.Errorf("expect not match, but match")
						}
					}
				}
			}
		})
	}
}

type kv struct {
	Key string `xml:"key"`
	Val string `xml:"val"`
}

func TestXmlMarshalWithProcInst(t *testing.T) {
	myKv := &kv{
		Key: "key1",
		Val: "值1",
	}
	var b bytes.Buffer
	encoder := xml.NewEncoder(&b)
	err := encoder.EncodeToken(xml.ProcInst{
		Target: "xml",
		Inst:   []byte(`version="1.0" encoding="GBK"`),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = encoder.EncodeToken(xml.CharData("\n"))
	if err != nil {
		t.Fatal(err)
	}
	encoder.Indent("", "\t")
	err = encoder.Encode(myKv)
	if err != nil {
		t.Fatal(err)
	}
	data := b.String()
	t.Logf("\n%s", data)
}
