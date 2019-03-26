package irbis

import (
	"strconv"
	"strings"
)

type RawRecord struct {
	Database string
	Mfn      int
	Status   int
	Version  int
	Fields   []string
}

func (record *RawRecord) Decode(lines []string) {
	// TODO implement
}

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
