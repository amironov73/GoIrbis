package irbis

import (
	"strconv"
	"strings"
)

// RawRecord запись с нераскодированными полями/подполями.
type RawRecord struct {
	Database string    // имя базы данных
	Mfn      int       // MFN записи
	Status   int       // Статус записи
	Version  int       // Версия записи
	Fields   []string  // Слайс с нераскодированными полями
}

// NewRawRecord Конструктор записи.
func NewRawRecord() *RawRecord {
	return new(RawRecord)
}

// Decode декодирует запись из протокольного представления.
func (record *RawRecord) Decode(lines []string) {
	// TODO implement
}

// Encode кодирует запись в протокольное представление.
func (record *RawRecord) Encode(delimiter string) string {
	result := strings.Builder{}
	result.WriteString(strconv.Itoa(record.Mfn))
	result.WriteRune('#')
	result.WriteString(strconv.Itoa(record.Status))
	result.WriteString(delimiter)
	result.WriteString("0#")
	result.WriteString(strconv.Itoa(record.Version))
	result.WriteString(delimiter)
	for _, field := range record.Fields {
		result.WriteString(field)
		result.WriteString(delimiter)
	}

	return result.String()
}

// IsDeleted Запись удалена?
func (record *RawRecord) IsDeleted() bool {
	return (record.Status & 3) != 0
}

// Reset сбрасывает состояние записи, отвязывая её от базы данных.
// Поля при этом остаются нетронутыми.
func (record *RawRecord) Reset() {
	record.Mfn = 0
	record.Status = 0
	record.Version = 0
	record.Database = ""
}

// String Выдает текстовое представление записи.
func (record *RawRecord) String() string {
	return record.Encode("\n")
}
