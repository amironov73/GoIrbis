package irbis

import (
	"strconv"
	"strings"
)

type RecordField struct {
	Tag       int
	Value     string
	Subfields []SubField
}

func NewRecordField(tag int, value string) *RecordField {
	return &RecordField{Tag: tag, Value: value}
}

func (field *RecordField) Add(code rune, value string) *RecordField {
	subfield := NewSubField(code, value)
	field.Subfields = append(field.Subfields, *subfield)

	return field
}

func (field *RecordField) Clear() *RecordField {
	field.Subfields = []SubField{}
	return field
}

func (field *RecordField) Decode(text string) {
	parts := strings.SplitN(text, "#", 2)
	field.Tag, _ = strconv.Atoi(parts[0])
	body := parts[1]
	all := strings.Split(body, "^")
	if body[0] != '^' {
		field.Value = all[0]
		all = all[1:]
	}
	for _, one := range all {
		if len(one) != 0 {
			subfield := SubField{}
			subfield.Decode(one)
			field.Subfields = append(field.Subfields, subfield)
		}
	}
}

func (field *RecordField) Encode() string {
	result := strings.Builder{}
	result.WriteString(strconv.Itoa(field.Tag))
	result.WriteRune('#')
	result.WriteString(field.Value)
	for i := range field.Subfields {
		subfield := &field.Subfields[i]
		result.WriteString(subfield.Encode())
	}

	return result.String()
}

func (field *RecordField) GetFirstSubField(code rune) *SubField {
	for i := range field.Subfields {
		candidate := &field.Subfields[i]
		if SameRune(candidate.Code, code) {
			return candidate
		}
	}

	return nil
}

func (field *RecordField) GetFirstSubFieldValue(code rune) string {
	for i := range field.Subfields {
		candidate := &field.Subfields[i]
		if SameRune(candidate.Code, code) {
			return candidate.Value
		}
	}

	return ""
}
