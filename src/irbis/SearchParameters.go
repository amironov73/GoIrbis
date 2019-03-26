package irbis

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
