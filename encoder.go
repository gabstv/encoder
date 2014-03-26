package encoder

// Original code borrowed from https://github.com/PuerkitoBio/martini-api-example
// TextEncoder and XmlEncoder has been removed. If someone really needs it, let me know.

// Supported tags:
// 	 - "out" if it sets to "false", value won't be set to field
import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"reflect"
	"strings"
)

const (
	JSON = "json"
	XML  = "xml"
)

// An Encoder implements an encoding format of values to be sent as response to
// requests on the API endpoints.
type Encoder interface {
	Encode(v ...interface{}) ([]byte, error)
	SetIndent(indent bool)
}

// Because `panic`s are caught by martini's Recovery handler, it can be used
// to return server-side errors (500). Some helpful text message should probably
// be sent, although not the technical error (which is printed in the log).
func Must(data []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return data
}

type JsonEncoder struct {
	indent bool
}
type XMLEncoder struct {
	indent bool
}

// jsonEncoder is an Encoder that produces JSON-formatted responses.
func (e JsonEncoder) Encode(v ...interface{}) ([]byte, error) {
	/*var data interface{} = v
	var result interface{}

	if v == nil {
		// So that empty results produces `[]` and not `null`
		data = []interface{}{}
	} else if len(v) == 1 {
		data = v[0]
	}

	t := reflect.TypeOf(data)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct {
		result = copyStruct(reflect.ValueOf(data), t).Interface()
	} else {
		result = data
	}

	b, err := json.Marshal(result)

	return b, err*/
	var buffer bytes.Buffer
	for _, val := range v {
		var b []byte
		if e.indent {
			b, _ = json.MarshalIndent(val, "", "  ")
		} else {
			b, _ = json.Marshal(val)
		}
		buffer.Write(b)
	}
	return buffer.Bytes(), nil
}

func (v JsonEncoder) SetIndent(indent bool) {
	v.indent = indent
}

// XMLEncoder is an Encoder that produces XML-formatted responses.
func (e XMLEncoder) Encode(v ...interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	for _, val := range v {
		var b []byte
		if e.indent {
			b, _ = xml.MarshalIndent(val, "", "  ")
		} else {
			b, _ = xml.Marshal(val)
		}
		buffer.Write(b)
	}
	return buffer.Bytes(), nil
}

func (v XMLEncoder) SetIndent(indent bool) {
	v.indent = indent
}

type XJEncoder struct {
	XMLEnc  XMLEncoder
	JsonEnc JsonEncoder
	Request *http.Request
	Default string
}

func (x XJEncoder) Encode(v ...interface{}) ([]byte, error) {
	ct := x.Request.Header.Get("Content-Type")
	if len(ct) > 0 {
		ct = strings.TrimSpace(strings.Split(ct, ";")[0])
	}
	if ct == "application/xml" {
		return x.XMLEnc.Encode(v...)
	}
	if ct == "application/json" {
		return x.JsonEnc.Encode(v...)
	}
	if len(x.Default) < 1 {
		return x.JsonEnc.Encode(v...)
	}
	if x.Default == JSON {
		return x.JsonEnc.Encode(v...)
	}
	if x.Default == XML {
		return x.XMLEnc.Encode(v...)
	}
	return nil, errors.New("No valid Content-Type request header was provided and no valil default format provided.")
}

func (v XJEncoder) SetIndent(indent bool) {
	v.JsonEnc.SetIndent(indent)
	v.XMLEnc.SetIndent(indent)
}

func NewXJEncoder(r *http.Request, defaultEncoder string) XJEncoder {
	return XJEncoder{
		XMLEnc:  XMLEncoder{},
		JsonEnc: JsonEncoder{},
		Request: r,
		Default: defaultEncoder,
	}
}

func copyStruct(v reflect.Value, t reflect.Type) reflect.Value {
	result := reflect.New(t).Elem()

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		if tag := t.Field(i).Tag.Get("out"); tag == "false" {
			continue
		}
		if tag := t.Field(i).Tag.Get("json"); tag == "-" {
			continue
		}

		vfield := v.Field(i)

		if vfield.Kind() == reflect.Interface {
			vfield = vfield.Elem()

			for vfield.Kind() == reflect.Ptr {
				vfield = vfield.Elem()
			}

			result.Field(i).Set(copyStruct(vfield, reflect.TypeOf(vfield.Interface())))
			continue
		}

		if vfield.Kind() == reflect.Struct {
			result.Field(i).Set(copyStruct(vfield, t.Field(i).Type))
			continue
		}

		result.Field(i).Set(v.Field(i))
	}

	return result
}
