package irbis

import (
	"strconv"
	"strings"
)

// DatabaseInfo Информация о базе данных ИРБИС.
type DatabaseInfo struct {
	Name                     string // Name Имя базы данных.
	Description              string // Description Описание базы данных в произвольной форме.
	MaxMfn                   int    // MaxMfn Максимальный MFN.
	LogicallyDeletedRecords  []int  // LogicallyDeletedRecords Логически удаленные записи.
	PhysicallyDeletedRecords []int  // PhysicallyDeletedRecords Физически удаленные записи.
	NonActualizedRecords     []int  // NonActualizedRecords Неактуализированные записи.
	LockedRecords            []int  // LockedRecords Заблокированные записи.
	DatabaseLocked           bool   // DatabaseLocked Признак блокировки базы данных в целом.
	ReadOnly                 bool   // ReadOnly База только для чтения
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
	for _, entry := range menu.Entries {
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
