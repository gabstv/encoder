package encoder

import (
	"encoding/json"
	"encoding/xml"
	// "fmt"
	"testing"
)

type Sample struct {
	Visible     string `json:"visible"`
	Hidden      string `json:"hidden" out:"false"`
	HiddenNOmit string `json:"hidden_n_omit,omitempty" out:"false"`
}

func TestEncoderJSON(t *testing.T) {
	src := &Sample{Visible: "visible", Hidden: "value of this field won't be exported", HiddenNOmit: "field will be completely omitted"}
	dst := &Sample{}

	enc := &JsonEncoder{}
	result, err := enc.Encode(src)
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(result, dst); err != nil {
		t.Fatal("Unmarshal error:", err)
	}

	if dst.Hidden != "" {
		t.Fatalf("Expected empty field 'Hidden', got %v\n", dst.Hidden)
	}
}

type SampleXML struct {
	Visible string `xml:"visible"`
	Hidden  string `xml:"-"`
	Omit    string `xml:"omit,omitempty"`
}

func TestEncoderXML(t *testing.T) {
	src := &SampleXML{Visible: "visible", Hidden: "value of this field won't be exported"}
	dst := &SampleXML{}

	enc := &XMLEncoder{}
	result, err := enc.Encode(src)
	if err != nil {
		t.Fatal(err)
	}

	if err := xml.Unmarshal(result, dst); err != nil {
		t.Fatal("Unmarshal error:", err)
	}

	if dst.Hidden != "" {
		t.Fatalf("Expected empty field 'Hidden', got %v\n", dst.Hidden)
	}
}

type A struct {
	Ref interface{} `json:"some_if"`
}

func TestInterface(t *testing.T) {
	src := &A{&Sample{"visible", "should_be_hidden", "should_be_omitted"}}
	dst := &A{&Sample{}}

	enc := &JsonEncoder{}
	result, err := enc.Encode(src)
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Println(string(result))

	if err := json.Unmarshal(result, dst); err != nil {
		t.Fatal("Unmarshal error:", err)
	}

	if dst.Ref.(*Sample).Hidden != "" {
		t.Fatalf("Expected empty field 'Hidden', got %v\n", dst.Ref.(*Sample).Hidden)
	}
}
