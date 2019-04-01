package irbis

// PostingParameters
type PostingParameters struct {
	// Database База данных.
	Database string

	// FirstPosting Номер первого постинга. Отсчет от 1.
	FirstPosting int

	// Format Формат.
	Format string

	// NumberOfPostings Требуемое количество постингов
	NumberOfPostings int

	// Term Термин.
	Term string

	// ListOfTerms Список термов
	ListOfTerms []string
}

func NewPostingParameters() *PostingParameters {
	return &PostingParameters{FirstPosting: 1}
}
