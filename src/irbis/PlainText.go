package irbis

import (
	"os"
	"strconv"
	"strings"
)

func (record *MarcRecord) ExportPlainText(file *os.File) error {
	_, err := file.WriteString(record.ToPlainText())
	if err != nil {
		return err
	}
	_, err = file.WriteString("*****\n")
	return err
}

func (record *MarcRecord) ToPlainText() string {
	result := strings.Builder{}
	for i := range record.Fields {
		field := record.Fields[i]
		result.WriteString(strconv.Itoa(field.Tag))
		result.WriteRune('#')
		result.WriteString(field.Value)
		for j := range field.Subfields {
			subfield := field.Subfields[j]
			result.WriteString(subfield.String())
		}
		result.WriteRune(10)
	}

	return result.String()
}
