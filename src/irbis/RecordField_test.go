package irbis

import (
	"testing"
)

func getField1() *RecordField {
	field := NewRecordField(461, "")
	field.Add('1', "2001#").
	Add('a', "Златая цепь").
	Add('e', "Записки. Повести. Рассказы").
	Add('f', "Бондарин С. А.").
	Add('v', "С. 76-132")
	return field
}

func getField2() *RecordField {
	field := NewRecordField(461, "")
	field.Add('1', "2001#").
	Add('a', "Златая цепь").
	Add('e', "Записки. Повести. Рассказы").
	Add('f', "Бондарин С. А.").
	Add('v', "С. 76-132").
	Add('1', "200#1#").
	Add('a', "Руслан и Людмила").
	Add('f', "Пушкин А. С.")
	return field
}

func TestNewRecordField_1(t *testing.T) {
	tag := 200
	value := "Value"
	field := NewRecordField(tag, value)
	if field.Tag != tag || field.Value != value {
		t.FailNow()
	}
	if len(field.Subfields) != 0 {
		t.FailNow()
	}
}

func TestRecordField_Add_1(t *testing.T) {
	code := 'a'
	value := "Value"
	field := new(RecordField)
	field.Add(code, value)
	if len(field.Subfields) != 1 {
		t.FailNow()
	}
	sf := field.Subfields[0]
	if sf.Code != code || sf.Value != value {
		t.FailNow()
	}
}

func TestRecordField_Clear_1(t *testing.T) {
	field := new(RecordField)
	field.Add('a', "Value")
	field.Clear()
	if len(field.Subfields) != 0 {
		t.FailNow()
	}
}

func TestRecordField_Decode_1(t *testing.T) {
	text := "200#^aTitle^bText^eSubtitle"
	field := new(RecordField)
	field.Decode(text)
	if field.Tag != 200 || len(field.Value) != 0 ||
		len(field.Subfields) != 3 ||
		field.Subfields[0].Code != 'a' ||
		field.Subfields[0].Value != "Title" ||
		field.Subfields[1].Code != 'b' ||
		field.Subfields[1].Value != "Text" ||
		field.Subfields[2].Code != 'e' ||
		field.Subfields[2].Value != "Subtitle" {
		t.FailNow()
	}
}

func TestRecordField_Decode_2(t *testing.T) {
	text := "300#Comment text"
	field := new(RecordField)
	field.Decode(text)
	if field.Tag != 300 || len(field.Subfields) != 0 ||
		field.Value != "Comment text" {
		t.FailNow()
	}
}

func TestRecordField_Decode_3(t *testing.T) {
	text := "400#Value^aSubA^bSubB"
	field := new(RecordField)
	field.Decode(text)
	if field.Tag != 400 || len(field.Subfields) != 2 ||
		field.Value != "Value" ||
		field.Subfields[0].Code != 'a' ||
		field.Subfields[0].Value != "SubA" ||
		field.Subfields[1].Code != 'b' ||
		field.Subfields[1].Value != "SubB" {
		t.FailNow()
	}
}

func TestRecordField_Encode_1(t *testing.T) {
	field := NewRecordField(200, "")
	field.Add('a', "Title")
	field.Add('e', "Subtitle")
	encoded := field.Encode()
	if encoded != "200#^aTitle^eSubtitle" {
		t.FailNow()
	}
}

func TestRecordField_Encode_2(t *testing.T) {
	field := NewRecordField(300, "Comment")
	encoded := field.Encode()
	if encoded != "300#Comment" {
		t.FailNow()
	}
}

func TestRecordField_GetEmbeddedFields_1(t *testing.T) {
	field := NewRecordField(200, "")
	embedded := field.GetEmbeddedFields()
	if len(embedded) != 0 {
		t.FailNow()
	}
}

func TestRecordField_GetEmbeddedFields_2(t *testing.T) {
	field:=getField1()
	embedded:=field.GetEmbeddedFields()
	if len(embedded) != 1 {
		t.FailNow()
	}
	if embedded[0].Tag != 200 {
		t.FailNow()
	}
	if len(embedded[0].Subfields) != 4 {
		t.FailNow()
	}
	if embedded[0].Subfields[0].Code != 'a' {
		t.FailNow()
	}
	if embedded[0].Subfields[0].Value != "Златая цепь" {
		t.FailNow()
	}
}

func TestRecordField_GetEmbeddedFields_3(t *testing.T) {
	field:=getField2()
	embedded:=field.GetEmbeddedFields()
	if len(embedded) != 2 {
		t.FailNow()
	}
	if embedded[0].Tag != 200 {
		t.FailNow()
	}
	if len(embedded[0].Subfields) != 4 {
		t.FailNow()
	}
	if embedded[0].Subfields[0].Code != 'a' {
		t.FailNow()
	}
	if embedded[0].Subfields[0].Value != "Златая цепь" {
		t.FailNow()
	}
	if embedded[1].Tag != 200 {
		t.FailNow()
	}
	if len(embedded[1].Subfields) != 2 {
		t.FailNow()
	}
	if embedded[1].Subfields[0].Code != 'a' {
		t.FailNow()
	}
	if embedded[1].Subfields[0].Value != "Руслан и Людмила" {
		t.FailNow()
	}
}

func TestRecordField_GetFirstSubField_1(t *testing.T) {
	field := NewRecordField(200, "")
	field.Add('a', "Title")
	field.Add('e', "Subtitle")
	sf := field.GetFirstSubField('a')
	if sf.Code != 'a' {
		t.FailNow()
	}
	sf = field.GetFirstSubField('z')
	if sf != nil {
		t.FailNow()
	}
}

func TestRecordField_GetFirstSubField_2(t *testing.T) {
	field := NewRecordField(300, "Comment")
	sf := field.GetFirstSubField('a')
	if sf != nil {
		t.FailNow()
	}
}

func TestRecordField_GetFirstSubFieldValue_1(t *testing.T) {
	field := NewRecordField(200, "")
	field.Add('a', "Title")
	field.Add('e', "Subtitle")
	sfv := field.GetFirstSubFieldValue('a')
	if sfv != "Title" {
		t.FailNow()
	}
	sfv = field.GetFirstSubFieldValue('z')
	if sfv != "" {
		t.FailNow()
	}
}

func TestRecordField_GetFirstSubFieldValue_2(t *testing.T) {
	field := NewRecordField(300, "Comment")
	sfv := field.GetFirstSubFieldValue('a')
	if sfv != "" {
		t.FailNow()
	}
}

func TestRecordField_InsertAt_1(t *testing.T) {
	field := NewRecordField(200, "")
	field.Add('a', "SubA")
	field.Add('c', "SubC")
	field.Add('d', "SubD")
	field.InsertAt(1, 'b', "SubB")
	if field.String() != "200#^aSubA^bSubB^cSubC^dSubD" {
		t.FailNow()
	}
	field.InsertAt(1, 'b', "SubB")
	if field.String() != "200#^aSubA^bSubB^bSubB^cSubC^dSubD" {
		t.FailNow()
	}
}

func TestRecordField_RemoveAt_1(t *testing.T) {
	field := NewRecordField(200, "")
	field.Add('a', "SubA")
	field.Add('b', "SubB")
	field.Add('c', "SubC")
	field.Add('d', "SubD")
	field.RemoveAt(1)
	if field.String() != "200#^aSubA^cSubC^dSubD" {
		t.FailNow()
	}
	field.RemoveAt(1)
	if field.String() != "200#^aSubA^dSubD" {
		t.FailNow()
	}
	field.RemoveAt(1)
	if field.String() != "200#^aSubA" {
		t.FailNow()
	}
	field.RemoveAt(0)
	if field.String() != "200#" {
		t.FailNow()
	}
}

func TestRecordField_RemoveSubfield_1(t *testing.T) {
	field := NewRecordField(200, "")
	field.Add('a', "SubA")
	field.Add('b', "SubB")
	field.Add('c', "SubC")
	field.Add('d', "SubD")
	field.RemoveSubfield('B')
	if field.String() != "200#^aSubA^cSubC^dSubD" {
		t.FailNow()
	}
	field.RemoveSubfield('b')
	if field.String() != "200#^aSubA^cSubC^dSubD" {
		t.FailNow()
	}
	field.RemoveSubfield('a')
	if field.String() != "200#^cSubC^dSubD" {
		t.FailNow()
	}
}

func TestRecordField_String_1(t *testing.T) {
	field := NewRecordField(200, "")
	field.Add('a', "Title")
	field.Add('e', "Subtitle")
	if field.String() != "200#^aTitle^eSubtitle" {
		t.FailNow()
	}
}

func TestRecordField_String_2(t *testing.T) {
	field := NewRecordField(300, "Comment")
	if field.String() != "300#Comment" {
		t.FailNow()
	}
}

func TestRecordField_Verify_1(t *testing.T) {
	field := NewRecordField(200, "")
	field.Add('a', "Title")
	field.Add('e', "Subtitle")
	if !field.Verify() {
		t.FailNow()
	}
}

func TestRecordField_Verify_2(t *testing.T) {
	field := NewRecordField(300, "Comment")
	if !field.Verify() {
		t.FailNow()
	}
}

func TestRecordField_Verify_3(t *testing.T) {
	field := NewRecordField(300, "")
	if field.Verify() {
		t.FailNow()
	}
}

func TestRecordField_Verify_4(t *testing.T) {
	field := NewRecordField(200, "")
	field.Add('a', "Title")
	field.Add('e', "")
	if field.Verify() {
		t.FailNow()
	}
}
