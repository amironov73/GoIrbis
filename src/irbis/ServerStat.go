package irbis

import (
	"strconv"
	"strings"
)

// ClientInfo Информация о клиенте, подключенном к серверу ИРБИС
// (не обязательно о текущем).
type ClientInfo struct {
	// Number порядковый номер.
	Number string

	// IPAddress Адрес клиента.
	IPAddress string

	// Port Порт клиента.
	Port string

	// Name Логин
	Name string

	// Id Идентификатор клиентской программы
	// (просто уникальное число).
	Id string

	// Workstation Клиентский АРМ.
	Workstation string

	// Registered Момент подключения к серверу.
	Registered string

	// Acknowledged Подследнее подтверждение, посланное серверу.
	Acknowledged string

	// LastCommand Последняя команда, посланная серверу.
	LastCommand string

	// CommandNumber Номер последней команлы.
	CommandNumber string
}

// Parse Разбор ответа сервера.
func (client *ClientInfo) Parse(lines []string) {
	client.Number = lines[0]
	client.IPAddress = lines[1]
	client.Port = lines[2]
	client.Name = lines[3]
	client.Id = lines[4]
	client.Workstation = lines[5]
	client.Registered = lines[6]
	client.Acknowledged = lines[7]
	client.LastCommand = lines[8]
	client.CommandNumber = lines[9]
}

func (client *ClientInfo) String() string {
	return client.IPAddress
}

// ServerStat Статистика работы ИРБИС-сервера.
type ServerStat struct {
	// RunningClients Подключенные клиенты.
	RunningClients []ClientInfo

	// ClientCount Число клиентов, подключенных в текущий момент.
	ClientCount int

	// TotalCommandCount Общее количество команд,
	// исполненных сервером с момента запуска.
	TotalCommandCount int
}

func (stat *ServerStat) Parse(lines []string) {
	stat.TotalCommandCount, _ = strconv.Atoi(lines[0])
	stat.ClientCount, _ = strconv.Atoi(lines[1])
	linesPerClient, err := strconv.Atoi(lines[2])
	if err != nil || linesPerClient < 8 {
		return
	}

	lines = lines[3:]

	for i := 0; i < stat.ClientCount; i++ {
		if len(lines) < linesPerClient {
			break
		}
		stat.RunningClients = append(stat.RunningClients, ClientInfo{})
		client := &stat.RunningClients[len(stat.RunningClients)-1]
		client.Parse(lines)
		lines = lines[linesPerClient+1:]
	}
}

func (stat *ServerStat) String() string {
	result := strings.Builder{}
	result.WriteString(strconv.Itoa(stat.TotalCommandCount))
	result.WriteString("\n")
	result.WriteString(strconv.Itoa(stat.ClientCount))
	result.WriteString("\n8\n")
	for _, client := range stat.RunningClients {
		result.WriteString(client.String())
		result.WriteString("\n")
	}

	return result.String()
}
