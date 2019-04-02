package irbis

import (
	"bytes"
	"strconv"
)

// ClientQuery формирует клиентский запрос из запрашиваемых элементов (строк и их фрагментов).
type ClientQuery struct {
	buffer *bytes.Buffer
}

// NewClientQuery формирует заголовок клиентского запроса.
func NewClientQuery(connection *IrbisConnection, command string) *ClientQuery {
	result := ClientQuery{}
	result.buffer = bytes.NewBuffer(nil)
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
	buf := toWin1251(text)
	query.buffer.Write(buf)
	return query
}

// AddUtf добавляет в запрос строку в кодировке UTF-8.
func (query *ClientQuery) AddUtf(text string) *ClientQuery {
	buf := toUtf8(text)
	query.buffer.Write(buf)
	return query
}

// Encode выдает сетевой пакет, который нужно отправить серверу.
func (query *ClientQuery) Encode() []byte {
	result := bytes.NewBuffer(nil)
	length := query.buffer.Len()
	prefix := strconv.Itoa(length) + "\n"
	result.WriteString(prefix)
	result.Write(query.buffer.Bytes())

	return result.Bytes()
}

// NewLine добавляет в запрос перевод строки (\n).
func (query *ClientQuery) NewLine() *ClientQuery {
	query.buffer.WriteByte(10)
	return query
}
