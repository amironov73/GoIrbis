package irbis

import (
	"strconv"
)

// ClientQuery формирует клиентский запрос из запрашиваемых элементов (строк и их фрагментов).
type ClientQuery struct {
	chunks [][]byte
}

// NewClientQuery формирует заголовок клиентского запроса.
func NewClientQuery(connection *Connection, command string) *ClientQuery {
	result := ClientQuery{}
	result.AddAnsi(command).NewLine()
	result.AddAnsi(connection.Workstation).NewLine()
	result.AddAnsi(command).NewLine()
	result.Add(connection.ClientId).NewLine()
	result.Add(connection.QueryId).NewLine()
	result.AddAnsi(connection.Password).NewLine()
	result.AddAnsi(connection.Username).NewLine()
	result.NewLine()
	result.NewLine()
	result.NewLine()
	return &result
}

// Add добавляет в запрос целое число.
func (query *ClientQuery) Add(value int) *ClientQuery {
	return query.AddAnsi(strconv.Itoa(value))
}

// AddAnsi добавляет в запрос строку в кодировке ANSI.
func (query *ClientQuery) AddAnsi(text string) *ClientQuery {
	buf := ToAnsi(text)
	query.chunks = append(query.chunks, buf)
	return query
}

// AddFormat добавляет строку формата, предварительно подготовив её.
// Также добавляется перевод строки.
func (query *ClientQuery) AddFormat(format string) bool {
	if len(format) == 0 {
		query.NewLine()
		return false
	}

	prepared := prepareFormat(trimLeft(format))

	if format[0] == '@' {
		query.AddAnsi(prepared)
	} else if format[0] == '!' {
		query.AddUtf(prepared)
	} else {
		query.AddUtf("!" + prepared)
	}
	query.NewLine()
	return true
}

// AddUtf добавляет в запрос строку в кодировке UTF-8.
func (query *ClientQuery) AddUtf(text string) *ClientQuery {
	buf := toUtf8(text)
	query.chunks = append(query.chunks, buf)
	return query
}

// Encode выдает сетевой пакет, который нужно отправить серверу.
func (query *ClientQuery) Encode() [][]byte {
	length := 0
	for i := range query.chunks {
		length += len(query.chunks[i])
	}
	prefix := strconv.Itoa(length) + "\n"
	result := [][]byte{toUtf8(prefix)}
	result = append(result, query.chunks...)

	return result
}

// NewLine добавляет в запрос перевод строки (\n).
func (query *ClientQuery) NewLine() *ClientQuery {
	query.chunks = append(query.chunks, []byte{10})
	return query
}
