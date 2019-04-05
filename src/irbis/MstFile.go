package irbis

import (
	"encoding/binary"
	"io"
	"os"
	"strconv"
)

const MstControlRecordSize = 36 // Размер управляющей записи в байтах.

// MstLeader - лидер MST-записи.
type MstLeader struct {
	Mfn          int32 // Номер записи в файле документов.
	Length       int32 // Длина записи.
	PreviousLow  int32 // Ссылка на предыдущую версию записи (младшая часть).
	PreviousHigh int32 // Ссылка на предыдущую версию записи (старшая часть).
	Base         int32 // Смещенеи полей переменной длины.
	Nvf          int32 // Число полей в записи.
	Version      int32 // Номер версии записи.
	Status       int32 // Статус записи.
}

func (leader *MstLeader) PreviousOffset() int64 {
	return (int64(leader.PreviousHigh) << 32) + int64(leader.PreviousLow)
}

// MstDictionaryEntry - запись в словаре MST.
type MstDictionaryEntry struct {
	Tag      int32 // Метка поля.
	Position int32 // Смещение.
	Length   int32 // Длина данных.
}

// MstField поле в MST-записи.
type MstField struct {
	Tag  int32  // Метка поля.
	Text string // Значение поля (без деления на подполя).
}

// Decode декодирует запись.
func (field *MstField) Decode() *RecordField {
	result := new(RecordField)
	result.Tag = int(field.Tag)
	result.DecodeBody(field.Text)
	return result
}

func (field *MstField) String() string {
	return strconv.Itoa(int(field.Tag)) + "#" + field.Text
}

// MstRecord - запись в MST-файле.
type MstRecord struct {
	Leader     MstLeader            // Лидер записи.
	Dictionary []MstDictionaryEntry // Словарь.
	Fields     []MstField           // Поля.
}

// Decode декодирует запись.
func (record *MstRecord) Decode() *MarcRecord {
	result := NewMarcRecord()
	result.Mfn = int(record.Leader.Mfn)
	result.Status = int(record.Leader.Status)
	result.Version = int(record.Leader.Version)
	nfields := len(record.Fields)
	result.Fields = make([]*RecordField, nfields)
	for i := 0; i < nfields; i++ {
		result.Fields[i] = record.Fields[i].Decode()
	}
	return result
}

// MstControlRecord - управляющая запись MST-файла.
type MstControlRecord struct {
	CtlMfn           int32 // Резерв.
	NextMfn          int32 // Номер, который будет назначен следующей созданной записи.
	NextPositionLow  int32 // Смещение свободного места в файле (младшая часть).
	NextPositionHigh int32 // Смещение свободного места в файле (старшая часть).
	MftType          int32 // Резерв.
	RecCnt           int32 // Резерв.
	Reserv1          int32 // Резерв.
	Reserv2          int32 // Резерв.
	Blocked          int32 // Индикатор блокировки базы данных.
}

// XrfFile - обёртка над MST-файлом.
type MstFile struct {
	file    *os.File
	Control MstControlRecord // Управляющая запись.
}

// OpenXrfFile открывает файл на чтение.
func OpenMstFile(filename string) (result *MstFile, err error) {
	var file *os.File
	file, err = os.Open(filename)
	if err != nil {
		return
	}

	result = new(MstFile)
	result.file = file
	err = binary.Read(file, binary.BigEndian, &result.Control)
	if err != nil {
		_ = file.Close()
		result = nil
		return
	}

	return result, nil
}

func (control *MstControlRecord) NextPosition() int64 {
	return (int64(control.NextPositionHigh) << 32) + int64(control.NextPositionLow)
}

// Close закрывает файл.
func (mst *MstFile) Close() {
	_ = mst.file.Close()
}

// ReadRecord читает запись по указанному смещению
func (mst *MstFile) ReadRecord(position int64) (result *MstRecord, err error) {
	_, err = mst.file.Seek(position, io.SeekStart)
	if err != nil {
		return
	}

	result = new(MstRecord)
	err = binary.Read(mst.file, binary.BigEndian, &result.Leader)
	if err != nil {
		result = nil
		return
	}

	nvf := result.Leader.Nvf
	result.Dictionary = make([]MstDictionaryEntry, nvf)
	err = binary.Read(mst.file, binary.BigEndian, &result.Dictionary)
	if err != nil {
		result = nil
		return
	}

	result.Fields = make([]MstField, nvf)

	remaining := result.Leader.Length - result.Leader.Base
	temp := make([]byte, remaining)
	_, err = mst.file.Read(temp)
	if err != nil {
		result = nil
		return
	}

	for i := int32(0); i < nvf; i++ {
		result.Fields[i].Tag = result.Dictionary[i].Tag
		ofs := result.Dictionary[i].Position
		raw := temp[ofs : ofs+result.Dictionary[i].Length]
		result.Fields[i].Text = string(raw)
	}

	return
}
