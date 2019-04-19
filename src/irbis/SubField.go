package irbis

import "strings"

// SubField Подполе записи. Состоит из кода и значения.
type SubField struct {
	// Code Код подполя
	Code rune

	// Value Значение подполя
	Value string
}

// NewSubField Конструктор, создает подполе
// с указанными кодом и значением.
func NewSubField(code rune, value string) *SubField {
	return &SubField{Code: code, Value: value}
}

// Clone клонирует подполе.
func (subfield *SubField) Clone() *SubField {
	result := new(SubField)
	result.Code = subfield.Code
	result.Value = strings.Repeat(subfield.Value, 1)
	return result
}

// Decode Декодирование подполя из протокольного представления.
func (subfield *SubField) Decode(text string) {
	runes := []rune(text)
	subfield.Code = runes[0]
	subfield.Value = text[1:]
}

// Encode Кодирование подполя в протокольное представление.
func (subfield *SubField) Encode() string {
	return "^" + string(subfield.Code) + subfield.Value
}

// String Выдает текстовое представление подполя.
func (subfield *SubField) String() string {
	return subfield.Encode()
}

// Verify Верификация подполя.
func (subfield *SubField) Verify() bool {
	return subfield.Code != '\x00' && subfield.Value != ""
}
