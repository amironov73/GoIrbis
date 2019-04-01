package irbis

// TableDefinition Данные для метода PrintTable
type TableDefinition struct {
	// Database Имя базы данных.
	Database string

	// Table Имя таблицы.
	Table string

	// Headers Заголовки таблицы.
	Headers []string

	// Mode Режим таблицы.
	Mode string

	// SearchQuery Поисковый запрос.
	SearchQuery string

	// MinMfn Минимальный MFN.
	MinMfn int

	// MaxMfn Максимальный MFN.
	MaxMfn int

	// SequentialQuery Запрос для последовательного поиска.
	SequentialQuery string

	// MfnList Список MFN, по которым строится таблица.
	MfnList []int
}

func (table *TableDefinition) String() string {
	return table.Table
}
