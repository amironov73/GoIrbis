package irbis

// WsLine - одна строчка ввода в рабочем листе.
type WsLine struct {
	Tag                string // Числовая метка поля.
	Title              string // Наименование поля.
	Repeatable         string // Повторяемость.
	Help               string // Индекс контекстной помощи.
	EditMode           string // Режим ввода.
	InputInfo          string // Дополнительная информация для расширенных средств ввода
	FormalVerification string // Формально-логический контроль
	Hint               string // Подсказка - текст помощи (инструкции), сопровождающий ввод в поле
	DefaultValue       string // Значение по умолчанию при создании новой записи
	Reserved           string // Используется при определенных режимах ввода
}

// Parse разбирает элемент ввода.
func (ws *WsLine) Parse(lines []string) {
	ws.Tag = lines[0]
	ws.Title = lines[1]
	ws.Repeatable = lines[2]
	ws.Help = lines[3]
	ws.EditMode = lines[4]
	ws.InputInfo = lines[5]
	ws.FormalVerification = lines[6]
	ws.Hint = lines[7]
	ws.DefaultValue = lines[8]
	ws.Reserved = lines[9]
}
