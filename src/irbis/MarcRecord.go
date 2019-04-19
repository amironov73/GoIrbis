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
	Fields []*RecordField
}

// NewMarcRecord Конструктор записи.
func NewMarcRecord() *MarcRecord {
	return new(MarcRecord)
}

// Add Добавление поля в запись.
func (record *MarcRecord) Add(tag int, value string) *RecordField {
	field := NewRecordField(tag, value)
	record.Fields = append(record.Fields, field)
	return field
}

// AddNonEmpty добавляет в запись непустое поле.
func (record *MarcRecord) AddNonEmpty(tag int, value string) *MarcRecord {
	if len(value) != 0 {
		record.Add(tag, value)
	}
	return record
}

// Clear Очистка записи (удаление всех полей).
func (record *MarcRecord) Clear() {
	record.Fields = []*RecordField{}
}

// Clone клонирует запись со всеми полями.
func (record *MarcRecord) Clone() *MarcRecord {
	result := new(MarcRecord)
	result.Database = record.Database
	result.Mfn = record.Mfn
	result.Version = record.Version
	result.Status = record.Status
	result.Fields = make([]*RecordField, len(record.Fields))
	for i := range record.Fields {
		result.Fields[i] = record.Fields[i].Clone()
	}
	return result
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
			field := new(RecordField)
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
	for _, field := range record.Fields {
		result.WriteString(field.Encode())
		result.WriteString(delimiter)
	}

	return result.String()
}

// FM Получение значения поля с указанной меткой.
// Если поле не найдено, возвращается пустая строка.
func (record *MarcRecord) FM(tag int) string {
	for _, field := range record.Fields {
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
	for _, field := range record.Fields {
		if field.Tag == tag {
			for _, subfield := range field.Subfields {
				if SameRune(subfield.Code, code) {
					return subfield.Value
				}
			}
		}
	}

	return ""
}

// FMA Получение слайса со значениями полей с указанной меткой.
func (record *MarcRecord) FMA(tag int) (result []string) {
	for _, field := range record.Fields {
		if field.Tag == tag && field.Value != "" {
			result = append(result, field.Value)
		}
	}

	return
}

// FSMA Получение слайса со значениями подполей с указанными
// меткой и кодом. Если подполя не найдены, возвращается
// слайс нулевой длины.
func (record *MarcRecord) FSMA(tag int, code rune) (result []string) {
	for _, field := range record.Fields {
		if field.Tag == tag {
			for _, subfield := range field.Subfields {
				if SameRune(subfield.Code, code) && subfield.Value != "" {
					result = append(result, subfield.Value)
				}
			}
		}
	}

	return
}

// GetField Получение указанного повторения поля с указанной меткой.
// Если поле не найдено, возвращается nil.
func (record *MarcRecord) GetField(tag, occurrence int) *RecordField {
	for _, field := range record.Fields {
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
	for _, field := range record.Fields {
		if field.Tag == tag {
			result = append(result, field)
		}
	}

	return
}

// GetFirstField Получение первого вхождения поля с указанной меткой.
// Если такого поля нет, возвращается nil.
func (record *MarcRecord) GetFirstField(tag int) *RecordField {
	for _, field := range record.Fields {
		if field.Tag == tag {
			return field
		}
	}

	return nil
}

// HaveField выясняет, есть ли в записи поле с указанной меткой.
func (record *MarcRecord) HaveField(tag int) bool {
	for _, field := range record.Fields {
		if field.Tag == tag {
			return true
		}
	}
	return false
}

// InsertAt вставляет поле в указанную позицию.
func (record *MarcRecord) InsertAt(i int, tag int, value string) *RecordField {
	s := record.Fields
	s = append(s, nil)
	copy(s[i+1:], s[i:])
	field := NewRecordField(tag, value)
	s[i] = field
	record.Fields = s
	return field
}

// IsDeleted Запись удалена?
func (record *MarcRecord) IsDeleted() bool {
	return (record.Status & 3) != 0
}

// RemoveAt удаляет поле в указанной позиции.
func (record *MarcRecord) RemoveAt(i int) *MarcRecord {
	s := record.Fields
	record.Fields = s[:i+copy(s[i:], s[i+1:])]
	return record
}

// RemoveField удаляет все поля с указанной меткой.
func (record *MarcRecord) RemoveField(tag int) *MarcRecord {
	for i := 0; i < len(record.Fields); {
		field := record.Fields[i]
		if field.Tag == tag {
			record.RemoveAt(i)
		} else {
			i++
		}
	}
	return record
}

// Reset сбрасывает состояние записи, отвязывая её от базы данных.
// Поля при этом остаются нетронутыми.
func (record *MarcRecord) Reset() {
	record.Mfn = 0
	record.Status = 0
	record.Version = 0
	record.Database = ""
}

// SetField устанавливает значение первого повторения поля
// с указанной меткой. Если такого поля нет, оно создаётся.
func (record *MarcRecord) SetField(tag int, value string) *MarcRecord {
	field := record.GetFirstField(tag)
	if field == nil {
		field = NewRecordField(tag, value)
	}
	field.Value = value
	return record
}

// SetSubfield устанавливает значение подполя первого повторения
// поля с указанной меткой. Если необходимые поля или подполе
// отсутствуют, они создаются.
func (record *MarcRecord) SetSubfield(tag int, code rune, value string) *MarcRecord {
	field := record.GetFirstField(tag)
	if field == nil {
		field = NewRecordField(tag, "")
	}
	field.SetSubfield(code, value)
	return record
}

// String Выдает текстовое представление записи.
func (record *MarcRecord) String() string {
	return record.Encode("\n")
}
