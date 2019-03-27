package irbis

import (
	//	"golang.org/x/text/encoding/charmap"
	"strings"
	"unicode"
)

const IrbisDelimiter = "\x1F\x1E"
const ShortDelimiter = "\x1E"

func contains(s []int, item int) bool {
	for _, one := range s {
		if one == item {
			return true
		}
	}

	return false
}

func toWin1251(text string) []byte {
	return cp1251FromUnicode(text)

	//encoder := charmap.Windows1251.NewEncoder()
	//result, _ := encoder.Bytes([]byte(text))
	//return result
}

func fromWin1251(buffer []byte) string {
	return cp1251ToUnicode(buffer)

	//decoder := charmap.Windows1251.NewDecoder()
	//temp, _ := decoder.Bytes(buffer)
	//result := string(temp)
	//return result
}

func toUtf8(text string) []byte {
	result := []byte(text)
	return result
}

func fromUtf8(buffer []byte) string {
	return string(buffer)
}

func removeComments(text string) string  {
	if len(text) == 0 || !strings.Contains(text, "/*") {
		return text
	}

	result := strings.Builder{}
	state := '\x00'
	chars := []rune(text)
	index := 0
	length := len(chars)
	result.Grow(length)

	for index < length {
		c := chars[index]

		switch state {
		case '\'', '"', '|':
			if c == state {
				state = '\x00'
			}
			result.WriteRune(c)

		default:
			if c == '/' {
				if (index + 1 < length) && (chars[index + 1] == '*') {
					for index < length {
						c = chars[index]
						if (c == '\r') || (c == '\n') {
							result.WriteRune(c)
							break
						}
						index++
					}
				} else {
					result.WriteRune(c)
				}
			} else if (c == '\'') || (c == '"') || (c == '|') {
				state = c
				result.WriteRune(c)
			} else {
				result.WriteRune(c)
			}
		}

		index++
	}

	return result.String()
}

func prepareFormat(text string) string {
	text = removeComments(text)
	length := len(text)
	if length == 0 {
		return text
	}

	flag := false
	chars := []rune(text)
	for i := range chars {
		if chars[i] < ' ' {
			flag = true
			break
		}
	}

	if !flag {
		return text
	}

	result := strings.Builder{}
	result.Grow(length)
	for i := range chars {
		c := chars[i]
		if c >= ' ' {
			result.WriteRune(c)
		}
	}

	return result.String()
}

func IrbisToDos(text string) string {
	return strings.ReplaceAll(text, IrbisDelimiter, "\n")
}

func IrbisToLines(text string) []string {
	return strings.Split(text, IrbisDelimiter)
}

func PickOne(lines ...string) string {
	for _, line := range lines {
		if len(line) != 0 {
			return line
		}
	}

	return ""
}

func SameRune(left, right rune) bool {
	return unicode.ToUpper(left) == unicode.ToUpper(right)
}

func SameString(left, right string) bool {
	return strings.EqualFold(left, right)
}
