package irbis

import "strconv"

// ProcessInfo Информация о запущенном на ИРБИС-сервере процессе.
type ProcessInfo struct {
	// Number Просто порядковый номер в списке.
	Number string

	// IPAddress С каким клиентом взаимодействует.
	IPAddress string

	// Name Логин оператора.
	Name string

	// ClientId Идентификатор клиента.
	ClientId string

	// Workstation Тип АРМ.
	Workstation string

	// Started Время запуска
	Started string

	// LastCommand Последняя выполненная (или выполняемая) команда.
	LastCommand string

	// CommandNumber Порядковый номер последней команды.
	CommandNumber string

	// ProcessId Идентификатор процесса.
	ProcessId string

	// State Состояние.
	State string
}

func ParseProcesses(lines []string) (result []ProcessInfo) {
	processCount, err := strconv.Atoi(lines[0])
	if err != nil {
		return
	}
	linesPerProcess, err := strconv.Atoi(lines[1])
	if err != nil {
		return
	}
	if processCount == 0 || linesPerProcess < 9 {
		return
	}

	lines = lines[2:]
	for i := 0; i < processCount; i++ {
		process := ProcessInfo{
			Number:        lines[0],
			IPAddress:     lines[1],
			Name:          lines[2],
			ClientId:      lines[3],
			Workstation:   lines[4],
			Started:       lines[5],
			LastCommand:   lines[6],
			CommandNumber: lines[7],
			ProcessId:     lines[8],
			State:         lines[9],
		}
		result = append(result, process)
		lines = lines[linesPerProcess:]
	}

	return
}

func (process *ProcessInfo) String() string {
	return process.ProcessId + " " + process.IPAddress + " " + process.Number
}
