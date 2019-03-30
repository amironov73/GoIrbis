package irbis

import (
	"strconv"
	"strings"
)

// MarcRecord Библиографическая запись.
// Составная единица базы данных.
// Состоит из произвольного количества полей.
type MarcRecord struct {
	// Database Имя базы данных, в которой хранится запись.
	Database string

	// Mfn MFN записи.
	Mfn int

	// Version Версия записи.
	Version int

	// Status Статус записи.
	Status int

	// Fields Поля записи.
	Fields []RecordField
}

// NewMarcRecord Конструктор записи.
func NewMarcRecord() *MarcRecord {
	return &MarcRecord{}
}

// Add Добавление поля в запись.
func (record *MarcRecord) Add(tag int, value string) *RecordField {
	field := NewRecordField(tag, value)
	record.Fields = append(record.Fields, *field)
	return &record.Fields[len(record.Fields)-1]
}

// Clear Очистка записи (удаление всех полей).
func (record *MarcRecord) Clear() {
	record.Fields = []RecordField{}
}

// Decode Декодирование записи из протокольного представления.
func (record *MarcRecord) Decode(lines []string) {
	firstLine := strings.Split(lines[0], "#")
	record.Mfn, _ = strconv.Atoi(firstLine[0])
	record.Status, _ = strconv.Atoi(firstLine[1])
	secondLine := strings.Split(lines[1], "#")
	record.Version, _ = strconv.Atoi(secondLine[1])
	length := len(lines)
	for i := 2; i < length; i++ {
		line := lines[i]
		if len(line) != 0 {
			field := RecordField{}
			field.Decode(line)
			record.Fields = append(record.Fields, field)
		}
	}
}

// Encode Кодирование записи в протокольное представление.
func (record *MarcRecord) Encode(delimiter string) string {
	result := strings.Builder{}
	result.WriteString(strconv.Itoa(record.Mfn))
	result.WriteRune('#')
	result.WriteString(strconv.Itoa(record.Status))
	result.WriteString(delimiter)
	result.WriteString("0#")
	result.WriteString(strconv.Itoa(record.Version))
	result.WriteString(delimiter)
	for i := range record.Fields {
		field := &record.Fields[i]
		result.WriteString(field.Encode())
		result.WriteString(delimiter)
	}

	return result.String()
}

// FM Получение значения поля с указанной меткой.
// Если поле не найдено, возвращается пустая строка.
func (record *MarcRecord) FM(tag int) string {
	for i := range record.Fields {
		field := &record.Fields[i]
		if field.Tag == tag {
			return field.Value
		}
	}

	return ""
}

// FSM Получение значения подполя с указанным кодом
// в поле с указанной меткой.
// Если поле или подполе не найдено, возвращается пустая строка.
func (record *MarcRecord) FSM(tag int, code rune) string {
	for i := range record.Fields {
		field := &record.Fields[i]
		if field.Tag == tag {
			for j := range field.Subfields {
				subfield := &field.Subfields[j]
				if SameRune(subfield.Code, code) {
					return subfield.Value
				}
			}
		}
	}

	return ""
}

// FMA Получение слайса полей с указанной меткой.
func (record *MarcRecord) FMA(tag int) (result []string) {
	for i := range record.Fields {
		field := &record.Fields[i]
		if field.Tag == tag && field.Value != "" {
			result = append(result, field.Value)
		}
	}

	return
}

// FSMA Получение слайса подполей с указанным кодом
// в полях с указанной меткой.
func (record *MarcRecord) FSMA(tag int, code rune) (result []string) {
	for i := range record.Fields {
		field := &record.Fields[i]
		if field.Tag == tag {
			for j := range field.Subfields {
				subfield := &field.Subfields[j]
				if SameRune(subfield.Code, code) && subfield.Value != "" {
					result = append(result, subfield.Value)
				}
			}
		}
	}

	return
}

// GetField Получение поля с указанной меткой с учетом повторения.
// Если поле не найдено, возвращается nil.
func (record *MarcRecord) GetField(tag, occurrence int) *RecordField {
	for i := range record.Fields {
		field := &record.Fields[i]
		if field.Tag == tag {
			if occurrence == 0 {
				return field
			}
			occurrence--
		}
	}

	return nil
}

// GetFields Получение слайса полей с указанной меткой.
func (record *MarcRecord) GetFields(tag int) (result []*RecordField) {
	for i := range record.Fields {
		field := &record.Fields[i]
		if field.Tag == tag {
			result = append(result, field)
		}
	}

	return
}

// GetFirstField Получение первого вхождения поля с указанной меткой.
// Если такого поля нет, возвращается nil.
func (record *MarcRecord) GetFirstField(tag int) *RecordField {
	for i := range record.Fields {
		field := &record.Fields[i]
		if field.Tag == tag {
			return field
		}
	}

	return nil
}

// IsDeleted Запись удалена?
func (record *MarcRecord) IsDeleted() bool {
	return (record.Status & 3) != 0
}

// String Выдает текстовое представление записи.
func (record *MarcRecord) String() string {
	return record.Encode("\n")
}
