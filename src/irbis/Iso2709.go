package irbis

import (
	"io"
	"strings"
)

const IsoMarkerLength = 24
const IsoRecordDelimiter = byte(0x1D)
const IsoFieldDelimiter = byte(0x1E)
const IsoSubfieldDelimiter = byte(0x1F)

func encodeInt32(buffer []byte, position, length, value int) {
	length--
	for position += length; length >= 0; length-- {
		buffer[position] = byte(value%10) + byte('0')
		value /= 10
		position--
	}
}

func encodeText(buffer []byte, position int, text string) int {
	if len(text) != 0 {
		encoded := []byte(text)
		for i := 0; i < len(encoded); i++ {
			buffer[position] = encoded[i]
			position++
		}
	}

	return position
}

// DecodeBody декодирует только текст поля.
func (field *RecordField) decodeBody(body string) {
	delimiter := IsoSubfieldDelimiter
	all := strings.Split(body, string(delimiter))
	if body[0] != delimiter {
		field.Value = all[0]
		all = all[1:]
	}
	for _, one := range all {
		if len(one) != 0 {
			subfield := new(SubField)
			subfield.Decode(one)
			field.Subfields = append(field.Subfields, subfield)
		}
	}
}

func ReadIsoRecord(reader io.Reader, decoder func([]byte) string) *MarcRecord {
	result := NewMarcRecord()

	// считываем длину записи
	marker := make([]byte, 5)
	if _, err := reader.Read(marker); err != nil {
		panic(err)
	}

	// а затем и ее остаток
	recordLength := ParseInt32(marker)
	record := make([]byte, recordLength)
	if _, err := reader.Read(record[len(marker):]); err != nil {
		panic(err)
	}

	// Простая проверка, что мы имеем дело с нормальной ISO-записью
	if record[recordLength-1] != IsoRecordDelimiter {
		panic("Not ISO record")
	}

	lengthOfLength := ParseInt32(record[20:21])
	lengthOfOffset := ParseInt32(record[21:22])
	additionalData := ParseInt32(record[22:23])
	directoryLength := 3 + lengthOfLength + lengthOfOffset + additionalData
	indicatorLength := ParseInt32(record[10:11])
	baseAddress := ParseInt32(record[12:17])

	// Подсчитываем количество полей в записи,
	// чтобы уменьшить трафик памяти в result.Fields
	fieldCount := 0
	for ofs := IsoMarkerLength; ; ofs += directoryLength {
		if record[ofs] == IsoFieldDelimiter {
			break
		}
		fieldCount++
	}
	result.Fields = make([]*RecordField, 0, fieldCount)

	// Пошли по полям при помощи справочника
	for directory := IsoMarkerLength; ; directory += directoryLength {
		// Переходим к следующему полю.
		// Если нарвались на разделитель, значит, справочник закончился
		if record[directory] == IsoFieldDelimiter {
			break
		}

		tag := ParseInt32(record[directory : directory+3])
		ofs := directory + 3
		fieldLength := ParseInt32(record[ofs : ofs+lengthOfLength])
		ofs = directory + 3 + lengthOfLength
		fieldOffset := baseAddress + ParseInt32(record[ofs:ofs+lengthOfOffset])
		field := NewRecordField(tag, "")
		result.Fields = append(result.Fields, field)
		if tag < 10 {
			// Фиксированное поле
			// не может содержать подполей и индикаторов
			temp := record[fieldOffset : fieldOffset+fieldLength-1]
			field.Value = decoder(temp)
		} else {
			// Поле переменной длины
			// Содержит два однобайтных индикатора
			// Может содержать подполя

			start := fieldOffset + indicatorLength
			stop := fieldOffset + fieldLength - indicatorLength + 1
			temp := record[start:stop]
			text := decoder(temp)
			field.decodeBody(text)
		}
	}

	return result
}
