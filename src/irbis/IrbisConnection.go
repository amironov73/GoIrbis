package irbis

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

type IrbisConnection struct {
	Host          string
	Port          int
	Username      string
	Password      string
	Database      string
	Workstation   string
	ClientId      int
	QueryId       int
	ServerVersion string
	Interval      int
	Connected     bool
	IniFile       *IniFile
}

//===================================================================

func NewConnection() *IrbisConnection {
	result := IrbisConnection{}
	result.Host = "127.0.0.1"
	result.Port = 6666
	result.Database = "IBIS"
	result.Workstation = "C"

	return &result
}

//===================================================================

func (connection *IrbisConnection) ActualizeRecord(database string, mfn int) bool {
	if !connection.Connected {
		return false
	}

	query := NewClientQuery(connection, "F")
	query.AddAnsi(database).NewLine()
	query.Add(mfn).NewLine()
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return false
	}

	return true
}

//===================================================================

func (connection *IrbisConnection) Connect() bool {
	if connection.Connected {
		return true
	}

AGAIN:
	connection.ClientId = 100000 + rand.Intn(900000)
	connection.QueryId = 1
	query := NewClientQuery(connection, "A")
	query.AddAnsi(connection.Username).NewLine()
	query.AddAnsi(connection.Password)
	response := connection.Execute(query)
	if response == nil {
		return false
	}

	if response.GetReturnCode() == -3337 {
		goto AGAIN
	}

	if response.ReturnCode < 0 {
		return false
	}

	connection.Connected = true
	connection.ServerVersion = response.ServerVersion
	connection.Interval = response.ReadInteger()
	lines := response.ReadRemainingAnsiLines()
	ini := NewIniFile()
	ini.Parse(lines)
	connection.IniFile = ini

	return true
}

//===================================================================

func (connection *IrbisConnection) CreateDatabase(database string,
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
	if response == nil || !response.CheckReturnCode() {
		return false
	}

	return true
}

//===================================================================

func (connection *IrbisConnection) CreateDictionary(database string) bool {
	if !connection.Connected {
		return false
	}

	query := NewClientQuery(connection, "Z")
	query.AddAnsi(database).NewLine()
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return false
	}

	return true
}

//===================================================================

func (connection *IrbisConnection) DeleteDatabase(database string) bool {
	if !connection.Connected {
		return false
	}

	query := NewClientQuery(connection, "W")
	query.AddAnsi(database).NewLine()
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return false
	}

	return true
}

//===================================================================

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

//===================================================================

func (connection *IrbisConnection) Execute(query *ClientQuery) *ServerResponse {
	address := connection.Host + ":" + strconv.Itoa(connection.Port)
	socket, err := net.Dial("tcp", address)
	if err != nil {
		return nil
	}

	defer func() { _ = socket.Close() }()

	buffer := query.Encode()
	_, err = socket.Write(buffer)
	if err != nil {
		return nil
	}

	result := NewServerResponse(socket)

	return result
}

//===================================================================

func (connection *IrbisConnection) FormatMfn(format string, mfn int) string {
	if !connection.Connected {
		return ""
	}

	query := NewClientQuery(connection, "G")
	query.AddAnsi(connection.Database).NewLine()
	prepared := prepareFormat(format)
	query.AddAnsi(prepared).NewLine()
	query.Add(1).NewLine()
	query.Add(mfn).NewLine()
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return ""
	}
	result := strings.TrimSpace(response.ReadRemainingUtfText())
	return result
}

//===================================================================

func (connection *IrbisConnection) FormatMfnUtf(format string, mfn int) string {
	if !connection.Connected {
		return ""
	}

	query := NewClientQuery(connection, "G")
	query.AddAnsi(connection.Database).NewLine()
	prepared := "!" + prepareFormat(format)
	query.AddUtf(prepared).NewLine()
	query.Add(1).NewLine()
	query.Add(mfn).NewLine()
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return ""
	}
	result := strings.TrimSpace(response.ReadRemainingUtfText())
	return result
}

//===================================================================

func (connection *IrbisConnection) FormatRecord(format string, record *MarcRecord) string {
	if !connection.Connected {
		return ""
	}
	database := PickOne(record.Database, connection.Database)
	query := NewClientQuery(connection, "G")
	query.AddAnsi(database).NewLine()
	prepared := prepareFormat(format)
	query.AddAnsi(prepared).NewLine()
	query.Add(-2).NewLine()
	query.AddUtf(record.Encode(IrbisDelimiter))
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return ""
	}
	result := response.ReadRemainingUtfText()
	return result
}

//===================================================================

func (connection *IrbisConnection) GetMaxMfn(database string) int {
	if !connection.Connected {
		return 0
	}

	database = PickOne(database, connection.Database)
	query := NewClientQuery(connection, "O")
	query.AddAnsi(database)
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return 0
	}

	return response.ReturnCode
}

//===================================================================

func (connection *IrbisConnection) NoOp() bool {
	if !connection.Connected {
		return false
	}

	query := NewClientQuery(connection, "N")
	connection.Execute(query)

	return true
}

//===================================================================

func (connection *IrbisConnection) ParseConnectionString(connectionString string) {
	items := strings.Split(connectionString, ";")
	for _, item := range items {
		if len(item) == 0 {
			continue
		}

		parts := strings.SplitN(item, "=", 2)
		if len(parts) != 2 {
			continue
		}

		name := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])

		switch name {
		case "host", "server", "address":
			connection.Host = value

		case "port":
			connection.Port, _ = strconv.Atoi(value)

		case "user", "username", "name", "login":
			connection.Username = value

		case "pwd", "password":
			connection.Password = value

		case "db", "database", "catalog":
			connection.Database = value

		case "arm", "workstation":
			connection.Workstation = value

		default:
			log.Println("Unknown connection key", name)

		}
	}
}

//===================================================================

func (connection *IrbisConnection) ReadMenuFile(specification string) *MenuFile {
	if !connection.Connected {
		return nil
	}

	lines := connection.ReadTextLines(specification)
	if lines == nil || len(lines) == 0 {
		return nil
	}

	result := new(MenuFile)
	result.Parse(lines)

	return result
}

//===================================================================

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

//===================================================================

func (connection *IrbisConnection) ReadRecord(mfn int) *MarcRecord {
	if !connection.Connected {
		return nil
	}

	query := NewClientQuery(connection, "C")
	query.AddAnsi(connection.Database).NewLine()
	query.Add(mfn).NewLine()
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return nil
	}

	result := NewMarcRecord()
	lines := response.ReadRemainingUtfLines()
	result.Decode(lines)
	result.Database = connection.Database

	return result
}

//===================================================================

func (connection *IrbisConnection) ReadTextLines(specification string) []string {
	if !connection.Connected {
		return []string{}
	}

	query := NewClientQuery(connection, "L")
	query.AddAnsi(specification).NewLine()
	response := connection.Execute(query)
	if response == nil {
		return []string{}
	}

	text := response.ReadAnsi()
	result := IrbisToLines(text)

	return result
}

//===================================================================

func (connection *IrbisConnection) ReadTextFile(specification string) string {
	if !connection.Connected {
		return ""
	}

	query := NewClientQuery(connection, "L")
	query.AddAnsi(specification).NewLine()
	response := connection.Execute(query)
	if response == nil {
		return ""
	}

	result := response.ReadAnsi()
	result = IrbisToDos(result)

	return result
}

//===================================================================

func (connection *IrbisConnection) Search(expression string) []int {
	if !connection.Connected {
		return []int{}
	}

	query := NewClientQuery(connection, "K")
	query.AddAnsi(connection.Database).NewLine()
	query.AddUtf(expression).NewLine()
	query.Add(0).NewLine()
	query.Add(1).NewLine()
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return []int{}
	}

	_ = response.ReadInteger() // Число найденных записей
	lines := response.ReadRemainingUtfLines()
	result := parseFoundMfn(lines)
	return result
}

//===================================================================

func (connection *IrbisConnection) SearchCount(expression string) int {
	if !connection.Connected {
		return 0
	}

	query := NewClientQuery(connection, "K")
	query.AddAnsi(connection.Database).NewLine()
	query.AddUtf(expression).NewLine()
	query.Add(0).NewLine()
	query.Add(0).NewLine()
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return 0
	}

	result := response.ReadInteger()
	return result
}

//===================================================================

func (connection *IrbisConnection) SearchEx(parameters *SearchParameters) []FoundLine {
	if !connection.Connected {
		return []FoundLine{}
	}

	database := PickOne(parameters.Database, connection.Database)
	query := NewClientQuery(connection, "K")
	query.AddAnsi(database).NewLine()
	query.AddUtf(parameters.Expression).NewLine()
	query.Add(parameters.NumberOfRecords).NewLine()
	query.Add(parameters.FirstRecord).NewLine()
	prepared := prepareFormat(parameters.Format)
	query.AddAnsi(prepared).NewLine()
	query.Add(parameters.MinMfn).NewLine()
	query.Add(parameters.MaxMfn).NewLine()
	query.AddAnsi(parameters.Sequential).NewLine()
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return []FoundLine{}
	}

	_ = response.ReadInteger() // Число найденных записей
	lines := response.ReadRemainingUtfLines()
	result := parseFoundLines(lines)
	return result
}

//===================================================================

func (connection *IrbisConnection) ToConnectionString() string {
	return "host=" + connection.Host +
		";port=" + strconv.Itoa(connection.Port) +
		";username=" + connection.Username +
		";password=" + connection.Password +
		";database=" + connection.Database +
		";arm=" + connection.Workstation + ";"

}

//===================================================================

func (connection *IrbisConnection) TruncateDatabase(database string) bool {
	if !connection.Connected {
		return false
	}

	query := NewClientQuery(connection, "S")
	query.AddAnsi(database).NewLine()
	connection.Execute(query)

	return true
}

//===================================================================

func (connection *IrbisConnection) UnlockDatabase(database string) bool {
	if !connection.Connected {
		return false
	}

	database = PickOne(database, connection.Database)
	query := NewClientQuery(connection, "U")
	query.AddAnsi(database).NewLine()
	connection.Execute(query)

	return true
}

//===================================================================

func (connection *IrbisConnection) UnlockRecords(database string,
	mfnList []int) bool {
	if !connection.Connected {
		return false
	}

	database = PickOne(database, connection.Database)
	query := NewClientQuery(connection, "Q")
	query.AddAnsi(database).NewLine()
	for _, mfn := range mfnList {
		query.Add(mfn).NewLine()
	}
	connection.Execute(query)

	return true
}

//===================================================================

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

//===================================================================

func (connection *IrbisConnection) WriteRawRecord(record *RawRecord) int {
	if !connection.Connected {
		return 0
	}

	database := PickOne(record.Database, connection.Database)
	query := NewClientQuery(connection, "D")
	query.AddAnsi(database).NewLine()
	query.Add(0).NewLine()
	query.Add(1).NewLine()
	query.AddUtf(record.Encode("\x001F\x001E")).NewLine()
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return 0
	}

	return response.ReturnCode
}

//===================================================================

func (connection *IrbisConnection) WriteRecord(record *MarcRecord) int {
	if !connection.Connected {
		return 0
	}

	database := PickOne(record.Database, connection.Database)
	query := NewClientQuery(connection, "D")
	query.AddAnsi(database).NewLine()
	query.Add(0).NewLine()
	query.Add(1).NewLine()
	query.AddUtf(record.Encode(IrbisDelimiter)).NewLine()
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return 0
	}

	// Decode the response
	temp := response.ReadRemainingUtfLines()
	if len(temp) != 0 {
		record.Clear()
		lines := append([]string{temp[0]}, strings.Split(temp[1], ShortDelimiter)...)
		record.Decode(lines)
		record.Database = connection.Database
	}

	return response.ReturnCode
}
