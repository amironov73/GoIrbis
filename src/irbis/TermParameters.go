package irbis

// TermParameters Параметры для запроса терминов с сервера.
type TermParameters struct {
	// Database Имя базы данных.
	Database string

	// NumberOfTerms Количество считываемых терминов.
	NumberOfTerms int

	// ReverseOrder Возвращать в обратном порядке?
	ReverseOrder bool

	// StartTerm Начальный термин.
	StartTerm string

	// Format Формат
	Format string
}
