package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io/ioutil"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"unicode"
)

const IrbisDelimiter = "\x1F\x1E"

func contains(s []int, item int) bool {
	for _, one := range s {
		if one == item {
			return true
		}
	}

	return false
}

func toWin1251(text string) []byte {
	encoder := charmap.Windows1251.NewEncoder()
	result, _ := encoder.Bytes([]byte(text))
	return result
}

func fromWin1251(buffer []byte) string {
	decoder := charmap.Windows1251.NewDecoder()
	temp, _ := decoder.Bytes(buffer)
	result := string(temp)
	return result
}

func toUtf8(text string) []byte {
	result := []byte(text)
	return result
}

func fromUtf8(buffer []byte) string {
	return string(buffer)
}

func IrbisToDos(text string) string {
	return strings.ReplaceAll(text, IrbisDelimiter, "\n")
}

func IrbisToLines(text string) []string {
	return strings.Split(text, IrbisDelimiter)
}

func PeekOne(lines ...string) string {
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

//=========================================================

type SubField struct {
	Code  rune
	Value string
}

func NewSubField(code rune, value string) *SubField {
	return &SubField{Code: code, Value: value}
}

func (subfield *SubField) Decode(text string) {
	runes := []rune(text)
	subfield.Code = runes[0]
	subfield.Value = text[1:]
}

func (subfield *SubField) Encode() string {
	return "^" + string(subfield.Code) + subfield.Value;
}

//=========================================================

type RecordField struct {
	Tag       int
	Value     string
	Subfields []SubField
}

func NewRecordField(tag int, value string) *RecordField {
	return &RecordField{Tag: tag, Value: value}
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

func (field *RecordField) Decode(text string) {
	parts := strings.SplitN(text, "#", 2)
	field.Tag, _ = strconv.Atoi(parts[0])
	body := parts[1]
	all := strings.Split(body, "^")
	if body[0] != '^' {
		field.Value = all[0]
		all = all[1:]
	}
	for _, one := range all {
		if len(one) != 0 {
			subfield := SubField{}
			subfield.Decode(one)
			field.Subfields = append(field.Subfields, subfield)
		}
	}
}

func (field *RecordField) Encode() string {
	result := strings.Builder{}
	result.WriteString(strconv.Itoa(field.Tag))
	result.WriteRune('#')
	result.WriteString(field.Value)
	for i := range field.Subfields {
		subfield := &field.Subfields[i]
		result.WriteString(subfield.Encode())
	}

	return result.String()
}

func (field *RecordField) GetFirstSubField(code rune) *SubField {
	for i := range field.Subfields {
		candidate := &field.Subfields[i]
		if SameRune(candidate.Code, code) {
			return candidate
		}
	}

	return nil
}

func (field *RecordField) GetFirstSubFieldValue(code rune) string {
	for i := range field.Subfields {
		candidate := &field.Subfields[i]
		if SameRune(candidate.Code, code) {
			return candidate.Value
		}
	}

	return ""
}

//=========================================================

type MarcRecord struct {
	Database string
	Mfn      int
	Version  int
	Status   int
	Fields   []RecordField
}

func NewMarcRecord() *MarcRecord {
	return &MarcRecord{}
}

func (record *MarcRecord) Add(tag int, value string) *RecordField {
	field := NewRecordField(tag, value)
	record.Fields = append(record.Fields, *field)
	return field; // ???
}

func (record *MarcRecord) Decode(lines []string) {
	firstLine := strings.Split(lines[0], "#")
	record.Mfn, _ = strconv.Atoi(firstLine[0])
	record.Status, _ = strconv.Atoi(firstLine[1])
	secondLine := strings.Split(lines[1], "#")
	record.Version, _ = strconv.Atoi(secondLine[1])
	length := len(lines)
	for i := 2; i < length; i++ {
		line := lines[i]
		if len(line) != 0 {
			field := RecordField{}
			field.Decode(line)
			record.Fields = append(record.Fields, field)
		}
	}
}

func (record *MarcRecord) Encode(delimiter string) string {
	result := strings.Builder{}
	result.WriteString(strconv.Itoa(record.Mfn))
	result.WriteRune('#')
	result.WriteString(strconv.Itoa(record.Status))
	result.WriteString(delimiter)
	result.WriteString("0#")
	result.WriteString(strconv.Itoa(record.Version))
	result.WriteString(delimiter)
	for i := range record.Fields {
		field := &record.Fields[i]
		result.WriteString(field.Encode())
		result.WriteString(delimiter)
	}

	return result.String()
}

func (record *MarcRecord) FM(tag int) string {
	for i := range record.Fields {
		field := &record.Fields[i]
		if field.Tag == tag {
			return field.Value
		}
	}

	return ""
}

func (record *MarcRecord) FMA(tag int) (result []string) {
	for i := range record.Fields {
		field := &record.Fields[i]
		if field.Tag == tag && field.Value != "" {
			result = append(result, field.Value)
		}
	}

	return
}

func (record *MarcRecord) GetField(tag, occurrence int) *RecordField {
	for i := range record.Fields {
		field := &record.Fields[i]
		if field.Tag == tag {
			if occurrence == 0 {
				return field
			}
			occurrence--
		}
	}

	return nil
}

func (record *MarcRecord) GetFields(tag int) (result []*RecordField) {
	for i := range record.Fields {
		field := &record.Fields[i]
		if field.Tag == tag {
			result = append(result, field)
		}
	}

	return
}

func (record *MarcRecord) IsDeleted() bool {
	return (record.Status & 3) != 0
}

//=========================================================

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
	result := strings.Builder{}
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

//=========================================================

type SearchParameters struct {
	Database        string
	FirstRecord     int
	Format          string
	MaxMfn          int
	MinMfn          int
	NumberOfRecords int
	Expression      string
	Sequential      string
	Filter          string
	IsUtf           bool
}

func NewSearchParameters() *SearchParameters {
	return &SearchParameters{
		FirstRecord: 1,
	}
}

//=========================================================

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

//=========================================================

type ServerResponse struct {
	Command       string
	ClientId      int
	QueryId       int
	AnswerSize    int
	ReturnCode    int
	ServerVersion string
	EOT           bool
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
	if response.EOT {
		return []byte{}
	}

	result := bytes.Buffer{}
	for response.reader.Len() != 0 {
		one, err := response.reader.ReadByte()
		if err != nil {
			break
		}
		if one == 13  {
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

	response.EOT = response.reader.Len() == 0

	return result.Bytes()
}

func (response *ServerResponse) GetReturnCode() int {
	response.ReturnCode = response.ReadInteger()
	return response.ReturnCode
}

func (response *ServerResponse) ReadAnsi() string {
	line := response.GetLine()
	result := fromWin1251(line)
	return result
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
	line := response.GetLine()
	result := fromUtf8(line)
	return result
}

//=========================================================

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
	if response == nil {
		return false
	}
	if response.GetReturnCode() < 0 {
		return false
	}

	connection.Connected = true

	// TODO implement

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
	address := connection.Host + ":" + strconv.Itoa(connection.Port)
	socket, err := net.Dial("tcp", address)
	if err != nil {
		return nil
	}

	defer socket.Close()

	buffer := query.Encode()
	_, err = socket.Write(buffer)
	if err != nil {
		return nil
	}

	result := NewServerResponse(socket)

	return result
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
	if !connection.Connected {
		return nil
	}

	query := NewClientQuery(connection, "C")
	query.AddAnsi(connection.Database).NewLine()
	query.Add(mfn).NewLine()
	response := connection.Execute(query)
	if response == nil {
		return nil
	}
	if !response.CheckReturnCode() {
		return nil
	}

	result := NewMarcRecord()
	lines := response.ReadRemainingUtfLines()
	result.Decode(lines)
	result.Database = connection.Database

	return result
}

func (connection *IrbisConnection) ToConnectionString() string {
	return "host=" + connection.Host +
		";port=" + strconv.Itoa(connection.Port) +
		";username=" + connection.Username +
		";password=" + connection.Password +
		";database=" + connection.Database +
		";arm=" + connection.Workstation + ";"

}

func (connection *IrbisConnection) TruncateDatabase(database string) bool {
	if !connection.Connected {
		return false
	}

	query := NewClientQuery(connection, "S")
	query.AddAnsi(database).NewLine()
	connection.Execute(query)

	return true
}

func (connection *IrbisConnection) UnlockDatabase(database string) bool {
	if !connection.Connected {
		return false
	}

	database = PeekOne(database, connection.Database)
	query := NewClientQuery(connection, "U")
	query.AddAnsi(database).NewLine()
	connection.Execute(query)

	return true
}

func (connection *IrbisConnection) UnlockRecords(database string,
	mfnList []int) bool {
	if !connection.Connected {
		return false
	}

	database = PeekOne(database, connection.Database)
	query := NewClientQuery(connection, "Q")
	query.AddAnsi(database).NewLine()
	for _, mfn := range mfnList {
		query.Add(mfn).NewLine()
	}
	connection.Execute(query)

	return true
}

func (connection *IrbisConnection) UpdateIniFile(lines []string) bool {
	if !connection.Connected {
		return false
	}

	if len(lines) == 0 {
		return true
	}

	query := NewClientQuery(connection, "8")
	for _, line := range lines {
		query.AddAnsi(line).NewLine()
	}
	connection.Execute(query)

	return true
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

//=========================================================

func main() {
	connection := NewIrbisConnection()
	connection.Username = "librarian"
	connection.Password = "secret"
	if !connection.Connect() {
		fmt.Println("Can't connect")
		return
	}

	defer connection.Disconnect()

	maxMfn := connection.GetMaxMfn("IBIS")
	fmt.Println("Max MFN", maxMfn)
	record := connection.ReadRecord(123);
	fmt.Println(record.Encode("\n"))
}
