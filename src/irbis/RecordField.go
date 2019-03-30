package irbis

import (
	"strconv"
	"strings"
)

// RecordField Поле записи. Состоит из метки и (опционального) значения.
// Может содержать произвольное количество подполей.
type RecordField struct {
	// Tag Метка поля.
	Tag int

	// Value Значение поля до первого разделителя.
	Value string

	// Subfields Подполя.
	Subfields []SubField
}

// NewRecordField Конструктор: создает поле с указанными меткой и значением.
func NewRecordField(tag int, value string) *RecordField {
	return &RecordField{Tag: tag, Value: value}
}

// Add Добавление подполя с указанными кодом и значением.
func (field *RecordField) Add(code rune, value string) *RecordField {
	subfield := NewSubField(code, value)
	field.Subfields = append(field.Subfields, *subfield)

	return field
}

// Clear Очищает поле (удаляет значение и все подполя).
func (field *RecordField) Clear() *RecordField {
	field.Subfields = []SubField{}
	return field
}

// Decode Декодирование поля из протокольного представления.
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

// Encode Кодирование поля в протокольное представление.
func (field *RecordField) Encode() string {
	result := strings.Builder{}
	result.WriteString(strconv.Itoa(field.Tag))
	result.WriteRune('#')
	result.WriteString(field.Value)
	for i := range field.Subfields {
		subfield := &field.Subfields[i]
		result.WriteString(subfield.String())
	}

	return result.String()
}

// GetFirstSubField Возвращает первое вхождение подполя с указанным кодом.
// Если такого подполя нет, возвращается nil.
func (field *RecordField) GetFirstSubField(code rune) *SubField {
	for i := range field.Subfields {
		candidate := &field.Subfields[i]
		if SameRune(candidate.Code, code) {
			return candidate
		}
	}

	return nil
}

// GetFirstSubFieldValue Возвращает значение первого вхождения
// подполя с указанным кодом, либо пустую строку, если такого подполя нет.
func (field *RecordField) GetFirstSubFieldValue(code rune) string {
	for i := range field.Subfields {
		candidate := &field.Subfields[i]
		if SameRune(candidate.Code, code) {
			return candidate.Value
		}
	}

	return ""
}

// String Выдает текстовое представление поля со всеми его подполями.
func (field *RecordField) String() string {
	return field.Encode()
}

// Verify Верификация поля со всеми его подполями.
func (field *RecordField) Verify() (result bool) {
	result = (field.Tag != 0) &&
		((field.Value != "") || (len(field.Subfields) != 0))
	if result && (len(field.Subfields) != 0) {
		for i := range field.Subfields {
			result = field.Subfields[i].Verify()
			if !result {
				break
			}
		}
	}

	return
}
