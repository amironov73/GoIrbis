package irbis

// GblStatement Оператор глобальной корректировки с параметрами.
type GblStatement struct {
	// Command Команда, например, ADD или DEL.
	Command string

	// Parameter1 Первый параметр, как правило,
	// спецификация поля/подполя.
	Parameter1 string

	// Parameter2 Второй параметр, как правило,
	// спецификация повторения.
	Parameter2 string

	// Format1 Первый формат, например, выражение для замены.
	Format1 string

	// Format2 Второй формат, например, заменяющее выражение.
	Format2 string
}

func (statement *GblStatement) Encode(delimiter string) string {
	return statement.Command + delimiter +
		statement.Parameter1 + delimiter +
		statement.Parameter2 + delimiter +
		statement.Format1 + delimiter +
		statement.Format2 + delimiter
}

func (statement *GblStatement) String() string {
	return statement.Encode("\n")
}

// GblSettings Установки для глобальной корректировки.
type GblSettings struct {
	// Actualize Актуализировать записи?
	Actualize bool

	// Autoin Запускать autoing.gbl?
	Autoin bool

	// Database Имя базы данных.
	Database string

	// Filename Имя файла.
	Filename string

	// FirstRecord MFN первой записи.
	FirstRecord int

	// FormalControl Применять формальный контроль?
	FormalControl bool

	// MaxMfn Максимальный MFN.
	MaxMfn int

	// MfnList Список MFN для обработки.
	MfnList []int

	// MinMfn Минимальный MFN.
	MinMfn int

	// NumberOfRecords Число обрабатываемых записей.
	NumberOfRecords int

	// SearchExpression Поисковое выражение.
	SearchExpression string

	// Statements Список операторов
	Statements []GblStatement
}
