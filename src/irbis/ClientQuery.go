package irbis

import (
	"bytes"
	"strconv"
)

type ClientQuery struct {
	buffer *bytes.Buffer
}

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

func (query *ClientQuery) Add(value int) *ClientQuery {
	return query.AddAnsi(strconv.Itoa(value))
}

func (query *ClientQuery) AddAnsi(text string) *ClientQuery {
	buf := toWin1251(text)
	query.buffer.Write(buf)
	return query
}

func (query *ClientQuery) AddUtf(text string) *ClientQuery {
	buf := toUtf8(text)
	query.buffer.Write(buf)
	return query
}

func (query *ClientQuery) Encode() []byte {
	result := bytes.NewBuffer(nil)
	length := query.buffer.Len()
	prefix := strconv.Itoa(length) + "\n"
	result.WriteString(prefix)
	result.Write(query.buffer.Bytes())

	return result.Bytes()
}

func (query *ClientQuery) NewLine() *ClientQuery {
	query.buffer.WriteByte(10)
	return query
}
