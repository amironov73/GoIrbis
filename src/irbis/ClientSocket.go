package irbis

import (
	"net"
	"strconv"
)

type ClientSocket interface {
	TalkToServer(query *ClientQuery) *ServerResponse
}

type Tcp4ClientSocket struct {
	connection *Connection
}

func NewTcp4ClientSocket(connection *Connection) *Tcp4ClientSocket {
	result := new(Tcp4ClientSocket)
	result.connection = connection
	return result
}

func (client *Tcp4ClientSocket) TalkToServer(query *ClientQuery) *ServerResponse {
	connection := client.connection
	address := connection.Host + ":" + strconv.Itoa(connection.Port)
	socket, err := net.Dial("tcp", address)
	if err != nil {
		return nil
	}

	defer func() { _ = socket.Close() }()

	chunks := query.Encode()
	for i := range chunks {
		_, err = socket.Write(chunks[i])
		if err != nil {
			return nil
		}
	}

	result := NewServerResponse(socket)

	return result
}
