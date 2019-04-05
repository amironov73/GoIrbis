package irbis

import (
	"encoding/binary"
	"io"
	"os"
)

const XrfRecordSize = 12

// XrfRecord содержит информацию о смещении записи и её статус.
type XrfRecord struct {
	Low    int32
	High   int32
	Status int32
}

// XrfFile - обёртка над XRF-файлом.
type XrfFile struct {
	file *os.File
}

// OpenXrfFile открывает файл на чтение.
func OpenXrfFile(filename string) (result *XrfFile, err error) {
	var file *os.File
	file, err = os.Open(filename)
	if err != nil {
		return
	}

	result = new(XrfFile)
	result.file = file

	return result, nil
}

// Close закрывает файл.
func (xrf *XrfFile) Close() {
	_ = xrf.file.Close()
}

func GetXrfOffset(mfn int) int64 {
	return int64(mfn-1) * XrfRecordSize
}

func (xrf *XrfFile) ReadRecord(mfn int) (result XrfRecord, err error) {
	offset := GetXrfOffset(mfn)
	_, err = xrf.file.Seek(offset, io.SeekStart)
	if err != nil {
		return
	}
	err = binary.Read(xrf.file, binary.BigEndian, &result)
	return
}

func (xrf *XrfRecord) Offset() int64 {
	return (int64(xrf.High) << 32) + int64(xrf.Low)
}
