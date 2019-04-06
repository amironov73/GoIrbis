package irbis

import (
	"bytes"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

type ServerResponse struct {
	Command       string
	ClientId      int
	QueryId       int
	AnswerSize    int
	ReturnCode    int
	ServerVersion string
	reader        *bytes.Reader
}

func NewServerResponse(conn net.Conn) *ServerResponse {
	result := &ServerResponse{}
	buffer, _ := ioutil.ReadAll(conn)
	result.reader = bytes.NewReader(buffer)
	result.Command = result.ReadAnsi()
	result.ClientId = result.ReadInteger()
	result.QueryId = result.ReadInteger()
	result.AnswerSize = result.ReadInteger()
	result.ServerVersion = result.ReadAnsi()
	result.ReadAnsi()
	result.ReadAnsi()
	result.ReadAnsi()
	result.ReadAnsi()
	result.ReadAnsi()
	return result
}

func (response *ServerResponse) CheckReturnCode(allowed ...int) bool {
	if response.GetReturnCode() < 0 {
		if contains(allowed, response.ReturnCode) {
			return true
		}
		return false
	}

	return true
}

func (response *ServerResponse) GetLine() []byte {
	//if response.EOT {
	//	return []byte{}
	//}

	result := bytes.Buffer{}
	for response.reader.Len() != 0 {
		one, err := response.reader.ReadByte()
		if err != nil {
			break
		}
		if one == 13 {
			one, err = response.reader.ReadByte()
			if err != nil {
				break
			}
			if one != 10 {
				_ = response.reader.UnreadByte()
			}
			break
		}

		result.WriteByte(one)
	}

	return result.Bytes()
}

func (response *ServerResponse) GetReturnCode() int {
	response.ReturnCode = response.ReadInteger()
	return response.ReturnCode
}

func (response *ServerResponse) ReadAnsi() string {
	line := response.GetLine()
	result := FromAnsi(line)
	return result
}

func (response *ServerResponse) ReadInteger() int {
	result, _ := strconv.Atoi(response.ReadAnsi())
	return result
}

func (response *ServerResponse) ReadRemainingAnsiLines() []string {
	text := response.ReadRemainingAnsiText()
	result := strings.Split(text, "\n")
	for i := range result {
		result[i] = strings.ReplaceAll(result[i], "\r", "")
	}
	return result
}

func (response *ServerResponse) ReadRemainingAnsiText() string {
	line, _ := ioutil.ReadAll(response.reader)
	result := FromAnsi(line)

	return result
}

func (response *ServerResponse) ReadRemainingUtfLines() []string {
	text := response.ReadRemainingUtfText()
	result := strings.Split(text, "\n")
	for i := range result {
		result[i] = strings.TrimSuffix(result[i], "\r")
	}
	return result
}

func (response *ServerResponse) ReadRemainingUtfText() string {
	line, _ := ioutil.ReadAll(response.reader)
	result := fromUtf8(line)

	return result
}

func (response *ServerResponse) ReadUtf() string {
	line := response.GetLine()
	result := fromUtf8(line)
	return result
}
