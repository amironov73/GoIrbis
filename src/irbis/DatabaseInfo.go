package irbis

import (
	"strconv"
	"strings"
)

// DatabaseInfo Информация о базе данных ИРБИС.
type DatabaseInfo struct {
	// Name Имя базы данных.
	Name string

	// Description Описание базы данных в произвольной форме.
	Description string

	// MaxMfn Максимальный MFN.
	MaxMfn int

	// LogicallyDeletedRecords Логически удаленные записи.
	LogicallyDeletedRecords []int

	// PhysicallyDeletedRecords Физически удаленные записи.
	PhysicallyDeletedRecords []int

	// NonActualizedRecords Неактуализированные записи.
	NonActualizedRecords []int

	// LockedRecords Заблокированные записи.
	LockedRecords []int

	// DatabaseLocked Признак блокировки базы данных в целом.
	DatabaseLocked bool

	// ReadOnly База только для чтения
	ReadOnly bool
}

func parseLine(line string) (result []int) {
	items := strings.Split(line, ShortDelimiter)
	for _, item := range items {
		value, err := strconv.Atoi(item)
		if err == nil {
			result = append(result, value)
		}
	}
	return
}

// Parse Разбор ответа сервера (см. GetDatabaseInfo)
func (db *DatabaseInfo) Parse(lines []string) {
	db.LogicallyDeletedRecords = parseLine(lines[0])
	db.PhysicallyDeletedRecords = parseLine(lines[1])
	db.NonActualizedRecords = parseLine(lines[2])
	db.LockedRecords = parseLine(lines[3])
	db.MaxMfn, _ = strconv.Atoi(lines[4])
	flag, _ := strconv.Atoi(lines[5])
	db.DatabaseLocked = flag != 0
}

func ParseMenu(menu *MenuFile) (result []DatabaseInfo) {
	for i := range menu.Entries {
		entry := &menu.Entries[i]
		name := entry.Code
		if name == "*****" {
			break
		}
		description := entry.Comment
		readOnly := false
		if name[0] == '-' {
			name = name[1:]
			readOnly = true
		}
		db := DatabaseInfo{}
		db.Name = name
		db.Description = description
		db.ReadOnly = readOnly
		result = append(result, db)
	}
	return
}

func (db *DatabaseInfo) String() string {
	return db.Name
}
