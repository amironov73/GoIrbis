package irbis

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

// IrbisConnection Подключение к серверу ИРБИС64.
type IrbisConnection struct {
	// Host Адрес сервера (можно задавать как my.domain.com,
	// так и 192.168.1.1).
	Host string

	// Port Порт сервера
	Port int

	// Username Логин пользователя. Регистр символов не учитывается.
	Username string

	// Password Пароль пользователя. Регистр символов учитывается.
	Password string

	// Database Имя текущей базы данных.
	Database string

	// Workstation Код АРМа.
	Workstation string

	// ClientId Идентификатор клиента. Задаётся автоматически
	// при подключении к серверу.
	ClientId int

	// QueryId Последовательный номер запроса к серверу.
	// Ведется автоматически
	QueryId int

	// ServerVersion Версия сервера
	// (становится доступна после подключения к нему).
	ServerVersion string

	// Interval Рекомендуемый интервал подключения, минуты.
	// Становится доступен после подключения к серверу.
	Interval int

	// Connected Признак подключения.
	Connected bool

	// Ini Серверный INI-файл (становится доступен после подключения).
	Ini *IniFile
}

//===================================================================

// NewConnection Конструктор, создает подключение с настройками по умолчанию.
func NewConnection() *IrbisConnection {
	result := IrbisConnection{}
	result.Host = "127.0.0.1"
	result.Port = 6666
	result.Database = "IBIS"
	result.Workstation = "C"

	return &result
}

//===================================================================

// ActualizeDatabase Актуализация всех неактуализированных записей
// в указанной базе данных.
func (connection *IrbisConnection) ActualizeDatabase(database string) bool {
	return connection.ActualizeRecord(database, 0)
}

//===================================================================

// ActualizeRecord Актуализация записи с указанным кодом.
// Если запись уже актуализирована, ничего не меняется.
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

// TODO возвращать ошибку

// Connect Подключение к серверу ИРБИС64.
// Если подключение уже установлено, ничего не меняется.
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
	connection.Ini = ini

	return true
}

//===================================================================

// CreateDatabase Создание базы данных.
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

// CreateDictionary Создание словаря в указанной базе данных.
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

// DeleteDatabase Удаление указанной базы данных.
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

// DeleteFile Удаление на сервере указанного файла.
func (connection *IrbisConnection) DeleteFile(fileName string) {
	connection.FormatMfn("&f('+9K"+fileName+"')", 1)
}

//===================================================================

// DeleteRecord Удаление записи по ее MFN.
func (connection *IrbisConnection) DeleteRecord(mfn int) {
	record := connection.ReadRecord(mfn)
	if record != nil && !record.IsDeleted() {
		record.Status |= LOGICALLY_DELETED
		connection.WriteRecord(record)
	}
}

//===================================================================

// Disconnect Отключение от сервера.
// Если подключение не установлено, ничего не меняется.
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

// Execute Отправка клиентского запроса на сервер
// и получение ответа от него.
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

// ExecuteAnyCommand Выполнение на сервере произвольной команды
// с опциональными параметрами в кодировке ANSI.
func (connection *IrbisConnection) ExecuteAnyCommand(command string, params ...string) bool {
	if !connection.Connected {
		return false
	}

	query := NewClientQuery(connection, command)
	for _, param := range params {
		query.AddAnsi(param).NewLine()
	}

	response := connection.Execute(query)
	if response == nil {
		return false
	}

	return true
}

//===================================================================

// FormatMfn Форматирование записи с указанным MFN.
func (connection *IrbisConnection) FormatMfn(format string, mfn int) string {
	if !connection.Connected {
		return ""
	}

	query := NewClientQuery(connection, "G")
	query.AddAnsi(connection.Database).NewLine()
	if !query.AddFormat(format) {
		return ""
	}

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

// FormatRecord Форматирование записи в клиентском представлении.
func (connection *IrbisConnection) FormatRecord(format string, record *MarcRecord) string {
	if !connection.Connected {
		return ""
	}
	database := PickOne(record.Database, connection.Database)
	query := NewClientQuery(connection, "G")
	query.AddAnsi(database).NewLine()
	if !query.AddFormat(format) {
		return ""
	}

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

// FormatRecords Расформатирование нескольких записей.
func (connection *IrbisConnection) FormatRecords(format string, list []int) (result []string) {
	if !connection.Connected || len(list) == 0 {
		return
	}

	query := NewClientQuery(connection, "G")
	query.AddAnsi(connection.Database).NewLine()
	if !query.AddFormat(format) {
		return make([]string, len(list))
	}

	query.Add(len(list)).NewLine()
	for _, mfn := range list {
		query.Add(mfn).NewLine()
	}

	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return
	}

	lines := response.ReadRemainingUtfLines()
	for _, line := range lines {
		parts := strings.SplitN(line, "#", 2)
		if len(parts) > 1 {
			result = append(result, parts[1])
		}
	}

	return
}

//===================================================================

// GetDatabaseInfo Получение информации об указанной базе данных.
func (connection *IrbisConnection) GetDatabaseInfo(database string) *DatabaseInfo {
	if !connection.Connected {
		return nil
	}

	query := NewClientQuery(connection, "0")
	query.AddAnsi(database)
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return nil
	}

	lines := response.ReadRemainingAnsiLines()
	result := new(DatabaseInfo)
	result.Parse(lines)
	result.Name = database
	return result
}

//===================================================================

// GetMaxMfn Получение максимального MFN для указанной базы данных.
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

// Получение статистики с сервера.
func (connection *IrbisConnection) GetServerStat() (result ServerStat) {
	if !connection.Connected {
		return
	}

	query := NewClientQuery(connection, "+1")
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return
	}

	lines := response.ReadRemainingAnsiLines()
	result.Parse(lines)

	return
}

//===================================================================

// GetServerVersion Получение версии сервера.
func (connection *IrbisConnection) GetServerVersion() (result VersionInfo) {
	if !connection.Connected {
		return
	}

	query := NewClientQuery(connection, "1")
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return
	}

	lines := response.ReadRemainingAnsiLines()
	result.Parse(lines)
	return
}

//===================================================================

// GetUserList Получение списка пользователей с сервера.
func (connection *IrbisConnection) GetUserList() (result []UserInfo) {
	if !connection.Connected {
		return
	}

	query := NewClientQuery(connection, "+9")
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return
	}

	lines := response.ReadRemainingAnsiLines()
	result = parseUsers(lines)
	return
}

//===================================================================

// GlobalCorrection Глобальная корректировка.
func (connection *IrbisConnection) GlobalCorrection(settings *GblSettings) (result []string) {
	if !connection.Connected {
		return
	}

	database := PickOne(settings.Database, connection.Database)
	query := NewClientQuery(connection, "5")
	query.AddAnsi(database).NewLine()
	query.Add(boolToInt(settings.Actualize)).NewLine()

	if len(settings.Filename) != 0 {
		query.AddAnsi("@" + settings.Filename).NewLine()
	} else {
		var encoded strings.Builder
		encoded.WriteString("!0")
		encoded.WriteString(IrbisDelimiter)
		for _, statement := range settings.Statements {
			encoded.WriteString(statement.Encode(IrbisDelimiter))
		}
		encoded.WriteString(IrbisDelimiter)
		query.AddUtf(encoded.String()).NewLine()
	}

	query.AddAnsi(settings.SearchExpression).NewLine()
	query.Add(settings.FirstRecord).NewLine()
	query.Add(settings.NumberOfRecords).NewLine()

	if len(settings.MfnList) != 0 {
		count := settings.MaxMfn - settings.MinMfn + 1
		query.Add(count).NewLine()
		for mfn := settings.MinMfn; mfn < settings.MaxMfn; mfn++ {
			query.Add(mfn).NewLine()
		}
	} else {
		query.Add(len(settings.MfnList)).NewLine()
		for _, mfn := range settings.MfnList {
			query.Add(mfn).NewLine()
		}
	}

	if !settings.FormalControl {
		query.AddAnsi("*").NewLine()
	}

	if !settings.Autoin {
		query.AddAnsi("&").NewLine()
	}

	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return
	}

	result = response.ReadRemainingAnsiLines()
	return
}

//===================================================================

// ListDatabases Получение списка баз данных с сервера.
func (connection *IrbisConnection) ListDatabases(specification string) (result []DatabaseInfo) {
	if !connection.Connected {
		return
	}

	if len(specification) == 0 {
		specification = "1..dbnam2.mnu"
	}

	menu := connection.ReadMenuFile(specification)
	if menu == nil {
		return
	}

	result = ParseMenu(menu)
	return
}

//===================================================================

// ListFiles Получение списка файлов на сервере.
func (connection *IrbisConnection) ListFiles(specifications ...string) (result []string) {
	if !connection.Connected {
		return
	}

	query := NewClientQuery(connection, "!")
	for _, specification := range specifications {
		query.AddAnsi(specification).NewLine()
	}

	response := connection.Execute(query)
	if response == nil {
		return
	}

	lines := response.ReadRemainingAnsiLines()
	for _, line := range lines {
		files := IrbisToLines(line)
		for _, file := range files {
			if len(file) != 0 {
				result = append(result, file)
			}
		}
	}

	return
}

//===================================================================

// ListProcesses Получение списка серверных процессов
func (connection *IrbisConnection) ListProcesses() (result []ProcessInfo) {
	if !connection.Connected {
		return
	}

	query := NewClientQuery(connection, "+3")
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return
	}

	lines := response.ReadRemainingAnsiLines()
	result = ParseProcesses(lines)
	return
}

//===================================================================

// ListTerms Получение списка терминов с указанным префиксом.
func (connection *IrbisConnection) ListTerms(prefix string) (result []string) {
	if !connection.Connected {
		return
	}

	prefixLength := len(prefix)
	startTerm := prefix
	lastTerm := startTerm
	flag := true
	for flag {
		terms := connection.ReadTerms(startTerm, 512)
		if len(terms) == 0 {
			break
		}
		for _, term := range terms {
			text := term.Text
			if !strings.HasPrefix(text, prefix) {
				flag = false
				break
			}
			if text != startTerm {
				lastTerm = text
				text = text[prefixLength:]
				result = append(result, text)
			}
		}

		startTerm = lastTerm
	}

	return
}

//===================================================================

// NoOp Пустая операция. Используется для периодического
// подтверждения подключения клиента.
func (connection *IrbisConnection) NoOp() bool {
	if !connection.Connected {
		return false
	}

	query := NewClientQuery(connection, "N")
	connection.Execute(query)

	return true
}

//===================================================================

// ParseConnectionString Разбор строки подключения.
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

func (connection *IrbisConnection) PrintTable(definition *TableDefinition) (result string) {
	if !connection.Connected {
		return
	}

	database := PickOne(definition.Database, connection.Database)
	query := NewClientQuery(connection, "7")
	query.AddAnsi(database).NewLine()
	query.AddAnsi(definition.Table).NewLine()
	query.AddAnsi(LinesToIrbis(definition.Headers)).NewLine()
	query.AddAnsi(definition.Mode).NewLine()
	query.AddAnsi(definition.SearchQuery).NewLine()
	query.Add(definition.MinMfn).NewLine()
	query.Add(definition.MaxMfn).NewLine()
	query.AddUtf(definition.SequentialQuery).NewLine()
	query.AddAnsi("") // Вместо перечня MFN
	response := connection.Execute(query)
	if response == nil {
		return
	}

	result = response.ReadRemainingUtfText()
	return
}

//===================================================================

// ReadBinaryFile Чтение двоичного файла с сервера.
func (connection *IrbisConnection) ReadBinaryFile(specification string) []byte {
	if !connection.Connected {
		return nil
	}

	// TODO implement

	return nil
}

//===================================================================

// ReadIniFile Чтение INI-файла с сервера.
func (connection *IrbisConnection) ReadIniFile(specification string) *IniFile {
	if !connection.Connected {
		return nil
	}

	lines := connection.ReadTextLines(specification)
	if len(lines) == 0 {
		return nil
	}

	result := new(IniFile)
	result.Parse(lines)

	return result
}

//===================================================================

// ReadMenuFile Чтение MNU-файла с сервера.
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

// ReadOptFile Чтение OPT-файла с сервера.
func (connection *IrbisConnection) ReadOptFile(specification string) (result *OptFile) {
	lines := connection.ReadTextLines(specification)
	if len(lines) == 0 {
		return
	}

	result = NewOptFile()
	result.Parse(lines)

	return result
}

//===================================================================

// ReadParFile Чтение PAR-файла с сервера
func (connection *IrbisConnection) ReadParFile(specification string) (result *ParFile) {
	lines := connection.ReadTextLines(specification)
	if len(lines) == 0 {
		return
	}

	result = NewParFile("")
	result.Parse(lines)
	return
}

//===================================================================

// ReadPostings Считывание постингов из поискового индекса.
func (connection *IrbisConnection) ReadPostings(parameters *PostingParameters) (result []TermPosting) {
	if !connection.Connected {
		return
	}

	database := PickOne(parameters.Database, connection.Database)
	query := NewClientQuery(connection, "I")
	query.AddAnsi(database).NewLine()
	query.Add(parameters.NumberOfPostings).NewLine()
	query.Add(parameters.FirstPosting).NewLine()
	prepared := prepareFormat(parameters.Format)
	query.AddUtf("!" + prepared).NewLine()
	if len(parameters.ListOfTerms) == 0 {
		query.AddUtf(parameters.Term).NewLine()
	} else {
		for _, term := range parameters.ListOfTerms {
			query.AddUtf(term).NewLine()
		}
	}

	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return
	}

	lines := response.ReadRemainingUtfLines()
	result = ParsePostings(lines)

	return
}

//===================================================================

// ReadRawRecord Чтение указанной записи в "сыром" виде.
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

// ReadRecord Чтение записи по ее MFN.
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

// ReadRecordVersion Чтение указанной версии записи.
func (connection *IrbisConnection) ReadRecordVersion(mfn, version int) *MarcRecord {
	if !connection.Connected {
		return nil
	}

	query := NewClientQuery(connection, "C")
	query.AddAnsi(connection.Database).NewLine()
	query.Add(mfn).NewLine()
	query.Add(version)
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

// ReadRecords Чтение с сервера нескольких записей.
func (connection *IrbisConnection) ReadRecords(mfnList []int) (result []MarcRecord) {
	if !connection.Connected || len(mfnList) == 0 {
		return
	}

	if len(mfnList) == 1 {
		record := connection.ReadRecord(mfnList[0])
		if record == nil {
			return
		}
		result = append(result, *record)
		return
	}

	query := NewClientQuery(connection, "G")
	query.AddAnsi(connection.Database).NewLine()
	query.AddAnsi(ALL_FORMAT).NewLine()
	query.Add(len(mfnList)).NewLine()
	for _, mfn := range mfnList {
		query.Add(mfn).NewLine()
	}

	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode() {
		return
	}

	lines := response.ReadRemainingUtfLines()
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		parts := strings.SplitN(line, "#", 2)
		if len(parts) != 2 {
			continue
		}
		parts = strings.Split(parts[1], "\x1F")[1:]
		record := NewMarcRecord()
		record.Decode(lines)
		record.Database = connection.Database
		result = append(result, *record)
	}

	return
}

//===================================================================

// ReadSearchScenario Загрузка сценариев поиска с сервера.
func (connection *IrbisConnection) ReadSearchScenario(specification string) (result []SearchScenario) {
	ini := connection.ReadIniFile(specification)
	if ini == nil || len(ini.Sections) == 0 {
		return
	}

	result = ParseScenarios(ini)
	return
}

//===================================================================

// ReadTerms Простое получение терминов поискового словаря.
func (connection *IrbisConnection) ReadTerms(startTerm string, number int) []TermInfo {
	parameters := TermParameters{StartTerm: startTerm, NumberOfTerms: number}
	return connection.ReadTermsEx(&parameters)
}

//===================================================================

// ReadTermsEx Получение терминов поискового словаря.
func (connection *IrbisConnection) ReadTermsEx(parameters *TermParameters) (result []TermInfo) {
	if !connection.Connected {
		return
	}

	command := "H"
	if parameters.ReverseOrder {
		command = "P"
	}

	database := PickOne(parameters.Database, connection.Database)
	query := NewClientQuery(connection, command)
	query.AddAnsi(database).NewLine()
	query.AddUtf(parameters.StartTerm).NewLine()
	query.Add(parameters.NumberOfTerms).NewLine()
	prepared := prepareFormat(parameters.Format)
	query.AddAnsi(prepared).NewLine()
	response := connection.Execute(query)
	if response == nil || !response.CheckReturnCode(-202, -203, -204) {
		return
	}

	lines := response.ReadRemainingUtfLines()
	result = ParseTerms(lines)

	return
}

//===================================================================

// ReadTextFile Чтение текстового файла с сервера.
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

// ReadTextLines Чтение текстового файла в виде слайса строк.
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

// ReadTreeFile Чтение TRE-файла с сервера.
func (connection *IrbisConnection) ReadTreeFile(specification string) (result *TreeFile) {
	lines := connection.ReadTextLines(specification)
	if len(lines) == 0 {
		return
	}

	result = new(TreeFile)
	result.Parse(lines)
	return
}

//===================================================================

// ReloadDictionary Пересоздание словаря для указанной базы данных.
func (connection *IrbisConnection) ReloadDictionary(database string) (result bool) {
	return connection.ExecuteAnyCommand("Y", database)
}

//===================================================================

// ReloadMasterFile Пересоздание мастер-файла для указанной базы данных.
func (connection *IrbisConnection) ReloadMasterFile(database string) (result bool) {
	return connection.ExecuteAnyCommand("X", database)
}

//===================================================================

// RestartServer Перезапуск сервера (без утери подключенных клиентов).
func (connection *IrbisConnection) RestartServer() (result bool) {
	return connection.ExecuteAnyCommand("+8")
}

//===================================================================

// Search Простой поиск записей (возвращается не более 32 тыс. записей).
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

// SearchAll Поиск всех записей (даже если их окажется больше 32 тыс.).
func (connection *IrbisConnection) SearchAll(expression string) (result []int) {
	if !connection.Connected {
		return
	}

	firstRecord := 1
	var totalCount int

	for {
		query := NewClientQuery(connection, "K")
		query.AddAnsi(connection.Database).NewLine()
		query.AddUtf(expression).NewLine()
		query.Add(10000).NewLine()
		query.Add(firstRecord).NewLine()

		response := connection.Execute(query)
		if response == nil || !response.CheckReturnCode() {
			break
		}

		if firstRecord == 1 {
			totalCount = response.ReadInteger()
			if totalCount == 0 {
				break
			}
		} else {
			_ = response.ReadInteger()
		}

		lines := response.ReadRemainingUtfLines()
		if len(lines) == 0 {
			break
		}
		for _, line := range lines {
			if len(line) != 0 {
				mfn, _ := strconv.Atoi(line)
				result = append(result, mfn)
				firstRecord++
			}
		}

		if firstRecord >= totalCount {
			break
		}
	}

	return
}

//===================================================================

// SearchCount Определение количества записей, соответствующих
// поисковому выражению.
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

// SearchEx Расширенный поиск записей.
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
	query.AddFormat(parameters.Format)
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

// SearchRead Поиск записей с их одновременным считыванием.
func (connection *IrbisConnection) SearchRead(expression string, limit int) (result []MarcRecord) {
	if !connection.Connected {
		return
	}

	parameters := NewSearchParameters()
	parameters.Expression = expression
	parameters.Format = ALL_FORMAT
	parameters.NumberOfRecords = limit
	found := connection.SearchEx(parameters)
	if len(found) == 0 {
		return
	}

	for _, item := range found {
		lines := strings.Split(item.Description, "\x1F")
		lines = lines[1:]
		record := MarcRecord{}
		record.Decode(lines)
		record.Database = connection.Database
		result = append(result, record)
	}

	return result
}

//===================================================================

// SearchSingleRecord Поиск и считывание одной записи, соответствующей выражение.
// Если таких записей больше одной, то будет считана любая из них.
// Если таких записей нет, будет возвращен nil.
func (connection *IrbisConnection) SearchSingleRecord(expression string) *MarcRecord {
	found := connection.SearchRead(expression, 1)
	if len(found) != 0 {
		return &found[0]
	}

	return nil
}

//===================================================================

// ToConnectionString Выдача строки подключения для текущего соеденения
// (соединение не обязательно должно быть установлено).
func (connection *IrbisConnection) ToConnectionString() string {
	return "host=" + connection.Host +
		";port=" + strconv.Itoa(connection.Port) +
		";username=" + connection.Username +
		";password=" + connection.Password +
		";database=" + connection.Database +
		";arm=" + connection.Workstation + ";"

}

//===================================================================

// TruncateDatabase Опустошение указанной базы данных.
func (connection *IrbisConnection) TruncateDatabase(database string) bool {
	return connection.ExecuteAnyCommand("S", database)
}

//===================================================================

// UndeleteRecord Восстановление записи по ее MFN.
func (connection *IrbisConnection) UndeleteRecord(mfn int) *MarcRecord {
	if !connection.Connected {
		return nil
	}

	record := connection.ReadRecord(mfn)
	if record == nil {
		return nil
	}

	if record.IsDeleted() {
		record.Status &= 0xFFFE
		if connection.WriteRecord(record) == 0 {
			return nil
		}
	}

	return record
}

//===================================================================

// UnlockDatabase Разблокирование указанной базы данных.
func (connection *IrbisConnection) UnlockDatabase(database string) bool {
	return connection.ExecuteAnyCommand("U", database)
}

//===================================================================

// UnlockRecords Разблокирование перечисленных записей.
func (connection *IrbisConnection) UnlockRecords(database string,
	mfnList []int) bool {
	if !connection.Connected {
		return false
	}

	if len(mfnList) == 0 {
		return true
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

// UpdateIniFile Обновление строк серверного INI-файла
// для текущего пользователя.
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

// UpdateUserList Обновление списка пользователей на сервере.
func (connection *IrbisConnection) UpdateUserList(users []UserInfo) {
	// TODO implement
}

//===================================================================

// WriteRawRecord Сохранение на сервере "сырой" записи.
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

// WriteRecord Сохранение записи на сервере.
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

//===================================================================

// WriteRecords Сохранение нескольких записей на сервере
// (могут относиться к разным базам).
func (connection *IrbisConnection) WriteRecords(records []MarcRecord) bool {
	if !connection.Connected || len(records) == 0 {
		return false
	}

	if len(records) == 1 {
		return connection.WriteRecord(&records[0]) != 0
	}

	query := NewClientQuery(connection, "6")
	query.Add(0).NewLine()
	query.Add(1).NewLine()
	for i := range records {
		record := &records[i]
		database := PickOne(record.Database, connection.Database)
		query.AddUtf(database)
		query.AddUtf(IrbisDelimiter)
		query.AddUtf(record.Encode(IrbisDelimiter))
		query.NewLine()
	}

	response := connection.Execute(query)
	if response == nil {
		return false
	}

	response.GetReturnCode()

	lines := response.ReadRemainingUtfLines()
	for i := 0; i < len(records); i++ {
		text := lines[i]
		lines := IrbisToLines(text)
		record := &records[i]
		record.Clear()
		record.Decode(lines)
		record.Database = PickOne(record.Database, connection.Database)
	}

	return true
}

//===================================================================

// WriteTextFile Сохранение текстового файла на сервере.
func (connection *IrbisConnection) WriteTextFile(specification, text string) bool {
	if !connection.Connected {
		return false
	}

	query := NewClientQuery(connection, "L")
	query.AddAnsi("&").AddAnsi(specification).AddAnsi("&").AddAnsi(DosToIrbis(text)).NewLine()

	response := connection.Execute(query)
	if response == nil {
		return false
	}

	return true
}
