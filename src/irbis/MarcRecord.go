package irbis

import (
	"strconv"
	"strings"
)

type MarcRecord struct {
	Database string
	Mfn      int
	Version  int
	Status   int
	Fields   []RecordField
}

func NewMarcRecord() *MarcRecord {
	return &MarcRecord{}
}

func (record *MarcRecord) Add(tag int, value string) *RecordField {
	field := NewRecordField(tag, value)
	record.Fields = append(record.Fields, *field)
	return field // ???
}

func (record *MarcRecord) Clear() {
	record.Fields = []RecordField{}
}

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

func (record *MarcRecord) FM(tag int) string {
	for i := range record.Fields {
		field := &record.Fields[i]
		if field.Tag == tag {
			return field.Value
		}
	}

	return ""
}

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

func (record *MarcRecord) FMA(tag int) (result []string) {
	for i := range record.Fields {
		field := &record.Fields[i]
		if field.Tag == tag && field.Value != "" {
			result = append(result, field.Value)
		}
	}

	return
}

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

func (record *MarcRecord) GetFields(tag int) (result []*RecordField) {
	for i := range record.Fields {
		field := &record.Fields[i]
		if field.Tag == tag {
			result = append(result, field)
		}
	}

	return
}

func (record *MarcRecord) IsDeleted() bool {
	return (record.Status & 3) != 0
}
