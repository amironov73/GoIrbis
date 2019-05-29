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
	Subfields []*SubField
}

// NewRecordField Конструктор: создает поле с указанными меткой и значением.
func NewRecordField(tag int, value string) *RecordField {
	return &RecordField{Tag: tag, Value: value}
}

// Add Добавление подполя с указанными кодом и значением.
func (field *RecordField) Add(code rune, value string) *RecordField {
	subfield := NewSubField(code, value)
	field.Subfields = append(field.Subfields, subfield)
	return field
}

// AddNonEmpty добавляет подполе, при условии, что его значение не пустое.
func (field *RecordField) AddNonEmpty(code rune, value string) *RecordField {
	if len(value) != 0 {
		field.Add(code, value)
	}
	return field
}

// Clear Очищает поле (удаляет значение и все подполя).
// Метка поля остаётся нетронутой.
func (field *RecordField) Clear() *RecordField {
	field.Subfields = []*SubField{}
	return field
}

// Clone клонирует поле со всеми подполями.
func (field *RecordField) Clone() *RecordField {
	result := new(RecordField)
	result.Tag = field.Tag
	result.Value = strings.Repeat(field.Value, 1)
	result.Subfields = make([]*SubField, len(field.Subfields))
	for i := range field.Subfields {
		result.Subfields[i] = field.Subfields[i].Clone()
	}
	return result
}

// DecodeBody декодирует только текст поля и подполей (без метки).
func (field *RecordField) DecodeBody(body string) {
	all := strings.Split(body, "^")
	if body[0] != '^' {
		field.Value = all[0]
		all = all[1:]
	}
	for _, one := range all {
		if len(one) != 0 {
			subfield := new(SubField)
			subfield.Decode(one)
			field.Subfields = append(field.Subfields, subfield)
		}
	}
}

// Decode Декодирование поля из протокольного представления
// (метка, значение и подполя).
func (field *RecordField) Decode(text string) {
	parts := strings.SplitN(text, "#", 2)
	field.Tag, _ = strconv.Atoi(parts[0])
	body := parts[1]
	field.DecodeBody(body)
}

// Encode Кодирование поля в протокольное представление
// (метка, значение и подполя).
func (field *RecordField) Encode() string {
	result := strings.Builder{}
	result.WriteString(strconv.Itoa(field.Tag))
	result.WriteRune('#')
	result.WriteString(field.Value)
	for _, subfield := range field.Subfields {
		result.WriteString(subfield.String())
	}

	return result.String()
}

// EncodeBody Кодирование поля в протокольное представление
// (только значение и подполя).
func (field *RecordField) EncodeBody() string {
	result := strings.Builder{}
	result.WriteString(field.Value)
	for _, subfield := range field.Subfields {
		result.WriteString(subfield.String())
	}

	return result.String()
}

// GetEmbeddedFields получает слайс встроенных полей из данного поля.
func (field *RecordField) GetEmbeddedFields() (result []*RecordField) {
	var found *RecordField = nil
	for _, one:=range field.Subfields {
		if one.Code == '1' {
			if found != nil {
				if len(found.Subfields) != 0 || len(found.Value) != 0 {
					result = append(result, found)
				}
			}
			value := one.Value
			if len(value) == 0 {
				continue
			}
			tag, _ := strconv.Atoi(value[:3])
			found = NewRecordField(tag, "")
			if tag < 10 && len(value) > 3 {
				found.Value = value[3:]
			}
		} else {
			if found != nil {
				found.Subfields = append(found.Subfields, one)
			}
		}
	}
	if found != nil {
		if len(found.Subfields) != 0 || len(found.Value) != 0 {
			result = append(result, found)
		}
	}
	return
}

// GetFirstSubField Возвращает первое вхождение подполя с указанным кодом.
// Если такого подполя нет, возвращается nil.
func (field *RecordField) GetFirstSubField(code rune) *SubField {
	for _, candidate := range field.Subfields {
		if SameRune(candidate.Code, code) {
			return candidate
		}
	}

	return nil
}

// GetFirstSubFieldValue Возвращает значение первого вхождения
// подполя с указанным кодом, либо пустую строку, если такого подполя нет.
func (field *RecordField) GetFirstSubFieldValue(code rune) string {
	for _, candidate := range field.Subfields {
		if SameRune(candidate.Code, code) {
			return candidate.Value
		}
	}

	return ""
}

// GetValueOrFirstSubField Выдаёт значение для ^*.
func (field *RecordField) GetValueOrFirstSubField() string {
	result := field.Value
	if len(result) == 0 && len(field.Subfields) != 0 {
		result = field.Subfields[0].Value
	}
	return result
}

// HaveSubField выясняет, есть ли подполе с указанным кодом.
func (field *RecordField) HaveSubField(code rune) bool {
	for _, candidate := range field.Subfields {
		if SameRune(candidate.Code, code) {
			return true
		}
	}

	return false
}

// InsertAt вставляет подполе в указанную позицию.
func (field *RecordField) InsertAt(i int, code rune, value string) *RecordField {
	s := field.Subfields
	s = append(s, nil)
	copy(s[i+1:], s[i:])
	s[i] = NewSubField(code, value)
	field.Subfields = s
	return field
}

// RemoveAt удаляет подполе в указанной позиции.
func (field *RecordField) RemoveAt(i int) *RecordField {
	s := field.Subfields
	field.Subfields = s[:i+copy(s[i:], s[i+1:])]
	return field
}

// RemoveSubfield удаляет все подполя с указанным кодом.
func (field *RecordField) RemoveSubfield(code rune) *RecordField {
	for i := 0; i < len(field.Subfields); {
		sf := field.Subfields[i]
		if SameRune(sf.Code, code) {
			field.RemoveAt(i)
		} else {
			i++
		}
	}
	return field
}

// ReplaceSubfield заменяет значение  подполя.
func (field *RecordField) ReplaceSubfield(code rune, oldValue, newValue string) *RecordField {
	for _, subfield := range field.Subfields {
		if SameRune(subfield.Code, code) &&
			SameString(subfield.Value, oldValue) {
			subfield.Value = newValue
		}
	}
	return field
}

// SetSubfield устанавливает значение первого повторения подполя с указанным кодом. Если value==nil, подполе удаляется.
func (field *RecordField) SetSubfield(code rune, value string) *RecordField {
	if len(value) == 0 {
		field.RemoveSubfield(code)
	} else {
		subfield := field.GetFirstSubField(code)
		if subfield == nil {
			subfield = NewSubField(code, value)
			field.Subfields = append(field.Subfields, subfield)
		}
		subfield.Value = value
	}
	return field
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
