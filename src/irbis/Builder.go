package irbis

import (
	"fmt"
	"strings"
)

// Search expression builder.
type Search struct {
	_buffer string
}

func All() (result Search) {
	result._buffer = "I=$"
	return result
}

func (query Search) And(items ...interface{}) Search {
	query._buffer = "(" + query._buffer
	for _, item := range items {
		text := fmt.Sprintf("%v", item)
		query._buffer = query._buffer + " * " + Wrap(text)
	}

	query._buffer = query._buffer + ")"
	return query
}

func Equals(prefix string, items ...interface{}) (result Search) {
	result._buffer = Wrap(prefix + fmt.Sprintf("%v", items[0]))
	if len(items) > 1 {
		items = items[1:]
		result._buffer = "(" + result._buffer
		for _, item := range items {
			result._buffer = result._buffer + " + " + Wrap(prefix+fmt.Sprintf("%v", item))
		}
		result._buffer = result._buffer + ")"
	}
	return
}

func (query Search) Not(value interface{}) Search {
	query._buffer = "(" + query._buffer + " ^ " + Wrap(fmt.Sprintf("%v", value)) + ")"
	return query
}

func NeedWrap(text string) bool {
	if len(text) == 0 {
		return true
	}

	c := text[0]
	if c == '"' || c == '(' {
		return false
	}

	if strings.Contains(text, " ") ||
		strings.Contains(text, "+") ||
		strings.Contains(text, "*") ||
		strings.Contains(text, "^") ||
		strings.Contains(text, "#") ||
		strings.Contains(text, "(") ||
		strings.Contains(text, ")") ||
		strings.Contains(text, "+") ||
		strings.Contains(text, "\"") {
		return true
	}

	return false
}

func (query Search) Or(items ...interface{}) Search {
	query._buffer = "(" + query._buffer
	for _, item := range items {
		text := fmt.Sprintf("%v", item)
		query._buffer = query._buffer + " + " + Wrap(text)
	}

	query._buffer = query._buffer + ")"
	return query
}

func (query Search) SameField(items ...interface{}) Search {
	query._buffer = "(" + query._buffer
	for _, item := range items {
		text := fmt.Sprintf("%v", item)
		query._buffer = query._buffer + " (G) " + Wrap(text)
	}

	query._buffer = query._buffer + ")"
	return query
}

func (query Search) SameRepeat(items ...interface{}) Search {
	query._buffer = "(" + query._buffer
	for _, item := range items {
		text := fmt.Sprintf("%v", item)
		query._buffer = query._buffer + " (F) " + Wrap(text)
	}

	query._buffer = query._buffer + ")"
	return query
}

func (query Search) String() string {
	return query._buffer
}

func Wrap(text string) string {
	if NeedWrap(text) {
		return "\"" + text + "\""
	}

	return text
}

func convert(values []string) []interface{} {
	result := make([]interface{}, len(values))
	for i := range values {
		result[i] = values[i]
	}
	return result
}

func Keyword(values ...string) Search {
	return Equals("K=", convert(values)...)
}

func Author(values ...string) Search {
	return Equals("A=", convert(values)...)
}

func Title(values ...string) Search {
	return Equals("T=", convert(values)...)
}

func Number(values ...string) Search {
	return Equals("IN=", convert(values)...)
}

func Publisher(values ...string) Search {
	return Equals("O=", convert(values)...)
}

func Place(values ...string) Search {
	return Equals("MI=", convert(values)...)
}

func Subject(values ...string) Search {
	return Equals("S=", convert(values)...)
}

func Language(values ...string) Search {
	return Equals("J=", convert(values)...)
}

func Year(values ...string) Search {
	return Equals("G=", convert(values)...)
}

func Magazine(values ...string) Search {
	return Equals("TJ=", convert(values)...)
}

func DocumentKind(values ...string) Search {
	return Equals("V=", convert(values)...)
}

func Udc(values ...string) Search {
	return Equals("U=", convert(values)...)
}

func Bbk(values ...string) Search {
	return Equals("BBK=", convert(values)...)
}

func Rzn(values ...string) Search {
	return Equals("RZN=", convert(values)...)
}

func Mhr(values ...string) Search {
	return Equals("MHR=", convert(values)...)
}
