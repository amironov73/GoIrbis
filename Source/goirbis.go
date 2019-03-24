package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

func PeekOne(lines ...string) string {
	for _, line := range lines {
		if len(line) != 0 {
			return line
		}
	}

	return ""
}

type SubField struct {
	Code  rune
	Value string
}

func NewSubField(code rune, value string) *SubField {
	return &SubField{Code: code, Value: value}
}

func DecodeSubField(text string) *SubField {
	runes := []rune(text)
	code := runes[0]
	value := text[1:]
	result := NewSubField(code, value)
	return result
}

type RecordField struct {
	Tag       int
	Value     string
	Subfields []SubField
}

func NewRecordField(tag int, value string) *RecordField {
	return &RecordField{Tag: tag, Value: value, Subfields: []SubField{}}
}

func (field *RecordField) Add(code rune, value string) *RecordField {
	subfield := NewSubField(code, value)
	field.Subfields = append(field.Subfields, *subfield)

	return field;
}

func (field *RecordField) Clear() *RecordField {
	field.Subfields = []SubField{}
	return field;
}

type MarcRecord struct {
	Database string
	Mfn      int
	Version  int
	Status   int
	Fields   []RecordField
}

func NewMarcRecord() *MarcRecord {
	return &MarcRecord{Fields: []RecordField{}}
}

func (record *MarcRecord) Add(tag int, value string) *RecordField {
	field := NewRecordField(tag, value)
	record.Fields = append(record.Fields, *field)
	return field; // ???
}

type RawRecord struct {
	Database string
	Mfn      int
	Status   int
	Version  int
	Fields   []string
}

func (record *RawRecord) Decode(lines []string) {
	// TODO implement
}

func (record *RawRecord) Encode(delimiter string) string {
	var result strings.Builder
	result.WriteString(strconv.Itoa(record.Mfn))
	result.WriteRune('#')
	result.WriteString(strconv.Itoa(record.Status))
	result.WriteString(delimiter)
	result.WriteString("0#")
	result.WriteString(strconv.Itoa(record.Version))
	result.WriteString(delimiter)
	for _, field := range record.Fields {
		result.WriteString(field)
		result.WriteString(delimiter)
	}

	return result.String()
}

type ClientQuery struct {
	Dummy int
}

func NewClientQuery(connection *IrbisConnection, command string) *ClientQuery {
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

func (query *ClientQuery) Add(value int) *ClientQuery {
	return query.AddAnsi(strconv.Itoa(value))
}

func (query *ClientQuery) AddAnsi(text string) *ClientQuery {
	// TODO implement
	return query
}

func (query *ClientQuery) AddUtf(text string) *ClientQuery {
	// TODO implement
	return query
}

func (query *ClientQuery) NewLine() *ClientQuery {
	// TODO implement
	return query
}

type ServerResponse struct {
	Command    string
	ClientId   int
	QueryId    int
	ReturnCode int
	EOT        bool
}

func NewServerResponse() *ServerResponse {
	return &ServerResponse{}
}

func (response *ServerResponse) CheckReturnCode(allowed ...int) bool {
	// TODO implement
	return false
}

func (response *ServerResponse) GetLine() []byte {
	// TODO impelement
	return []byte{}
}

func (response *ServerResponse) GetReturnCode() int {
	// TODO implement
	return 0
}

func (response *ServerResponse) ReadAnsi() string {
	// TODO implement
	return ""
}

func (response *ServerResponse) ReadInteger() int {
	result, _ := strconv.Atoi(response.ReadAnsi())
	return result
}

func (response *ServerResponse) ReadRemainingAnsiLines() []string {
	var result []string
	for !response.EOT {
		line := response.ReadAnsi()
		result = append(result, line)
	}
	return result
}

func (response *ServerResponse) ReadRemainingAnsiText() string {
	// TODO implement
	return ""
}

func (response *ServerResponse) ReadRemainingUtfLines() []string {
	var result []string
	for !response.EOT {
		line := response.ReadUtf()
		result = append(result, line)
	}
	return result
}

func (response *ServerResponse) ReadRemainingUtfText() string {
	// TODO implement
	return ""
}

func (response *ServerResponse) ReadUtf() string {
	// TODO implement
	return ""
}

type IrbisConnection struct {
	Host        string
	Port        int
	Username    string
	Password    string
	Database    string
	Workstation string
	ClientId    int
	QueryId     int
	Connected   bool
}

func NewIrbisConnection() *IrbisConnection {
	result := IrbisConnection{}
	result.Host = "127.0.0.1"
	result.Port = 6666
	result.Database = "IBIS"
	result.Workstation = "C"

	return &result
}

func (connection *IrbisConnection) ActualizeRecord(database string, mfn int) bool {
	if !connection.Connected {
		return false
	}

	query := NewClientQuery(connection, "F")
	query.AddAnsi(database).NewLine()
	query.Add(mfn).NewLine()
	response := connection.Execute(query)
	if !response.CheckReturnCode() {
		return false
	}

	return true
}

func (connection *IrbisConnection) Connect() bool {
	if (connection.Connected) {
		return true
	}

	connection.ClientId = 100000 + rand.Intn(900000)
	connection.QueryId = 1
	query := NewClientQuery(connection, "A")
	query.AddAnsi(connection.Username).NewLine()
	query.AddAnsi(connection.Password)
	response := connection.Execute(query)
	if response.GetReturnCode() < 0 {
		return false
	}

	connection.Connected = true

	// TODO implement

	fmt.Println("Connected")
	return true
}

func (connection *IrbisConnection) Disconnect() bool {
	if !connection.Connected {
		return true
	}

	query := NewClientQuery(connection, "B")
	query.AddAnsi(connection.Username)
	connection.Execute(query)
	connection.Connected = false

	fmt.Println("Disconnected")
	return true
}

func (connection *IrbisConnection) createDatabase(database string,
	description string, readerAccess bool) bool {
	if !connection.Connected {
		return false
	}

	query := NewClientQuery(connection, "T")
	query.AddAnsi(database).NewLine()
	query.AddAnsi(description).NewLine()
	value := 1
	if !readerAccess {
		value = 0
	}
	query.Add(value).NewLine()
	response := connection.Execute(query)
	return response.CheckReturnCode()
}

func (connection *IrbisConnection) createDictionary(database string) bool {
	if !connection.Connected {
		return false
	}

	query := NewClientQuery(connection, "Z")
	query.AddAnsi(database).NewLine()
	response := connection.Execute(query)
	return response.CheckReturnCode()
}

func (connection *IrbisConnection) DeleteDatabase(database string) bool {
	if !connection.Connected {
		return false
	}

	query := NewClientQuery(connection, "W")
	query.AddAnsi(database).NewLine()
	response := connection.Execute(query)
	return response.CheckReturnCode()
}

func (connection *IrbisConnection) Execute(query *ClientQuery) *ServerResponse {
	// TODO implement
	return NewServerResponse()
}

func (connection *IrbisConnection) GetMaxMfn(database string) int {
	if !connection.Connected {
		return 0
	}

	database = PeekOne(database, connection.Database)
	query := NewClientQuery(connection, "O")
	query.AddAnsi(database)
	response := connection.Execute(query)
	if !response.CheckReturnCode() {
		return 0
	}

	return response.ReturnCode
}

func (connection *IrbisConnection) ReadRawRecord(mfn int) *RawRecord {
	if !connection.Connected {
		return nil
	}

	query := NewClientQuery(connection, "C")
	query.AddAnsi(connection.Database).NewLine()
	query.Add(mfn).NewLine()
	response := connection.Execute(query)
	if !response.CheckReturnCode(-201, -600, -601, -602) {
		return nil
	}

	result := RawRecord{}
	lines := response.ReadRemainingUtfLines()
	result.Decode(lines)
	result.Database = connection.Database

	return &result
}

func (connection *IrbisConnection) ReadRecord(mfn int) *MarcRecord {
	// TODO implement
	return &MarcRecord{}
}

func (connection *IrbisConnection) WriteRawRecord(record *RawRecord) int {
	if !connection.Connected {
		return 0
	}

	database := PeekOne(record.Database, connection.Database)
	query := NewClientQuery(connection, "D")
	query.AddAnsi(database).NewLine()
	query.Add(0).NewLine()
	query.Add(1).NewLine()
	query.AddUtf(record.Encode("\x001F\x001E")).NewLine()
	response := connection.Execute(query)
	if !response.CheckReturnCode() {
		return 0
	}

	return response.ReturnCode
}

func (connection *IrbisConnection) WriteRecord(record *MarcRecord) int {
	// TODO implement
	return 0
}

func main() {
	connection := NewIrbisConnection()
	if !connection.Connect() {
		fmt.Println("Can't connect")
	}

	maxMfn := connection.GetMaxMfn("IBIS")
	fmt.Println("Max MFN", maxMfn)
	record := connection.ReadRecord(123);
	fmt.Println(record)
	connection.Disconnect()
}
