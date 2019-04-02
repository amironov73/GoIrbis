package irbis

import (
	"testing"
)

func TestNewSubField_1(t *testing.T) {
	code := 'a'
	value := "Value"
	sf := NewSubField(code, value)
	if sf.Code != code || sf.Value != value {
		t.FailNow()
	}
}

func TestSubField_Decode_1(t *testing.T) {
	text := "aValue"
	sf := new(SubField)
	sf.Decode(text)
	if sf.Code != 'a' || sf.Value != "Value" {
		t.FailNow()
	}
}

func TestSubField_Encode_1(t *testing.T) {
	sf := NewSubField('a', "Value")
	encoded := sf.Encode()
	if encoded != "^aValue" {
		t.FailNow()
	}
}

func TestSubField_String_1(t *testing.T) {
	sf := NewSubField('a', "Value")
	encoded := sf.String()
	if encoded != "^aValue" {
		t.FailNow()
	}
}

func TestSubField_Verify_1(t *testing.T) {
	sf := NewSubField('a', "Value")
	if !sf.Verify() {
		t.FailNow()
	}
}

func TestSubField_Verify_2(t *testing.T) {
	sf := NewSubField('a', "")
	if sf.Verify() {
		t.FailNow()
	}
}

func TestSubField_Verify_3(t *testing.T) {
	sf := new(SubField)
	if sf.Verify() {
		t.FailNow()
	}
}
