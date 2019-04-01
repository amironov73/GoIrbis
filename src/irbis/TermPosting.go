package irbis

import (
	"strconv"
	"strings"
)

// TermPosting Постинг термина в поисковом индексе.
type TermPosting struct {
	// Mfn MFN записи с искомым термином.
	Mfn int

	// Tag Метка поля с искомым термином.
	Tag int

	// Occurrence Номер повторения поля.
	Occurrence int

	// Count Количество повторений.
	Count int

	// Text Результат форматирования.
	Text string
}

func ParsePostings(lines []string) (result []TermPosting) {
	for _, line := range lines {
		parts := strings.SplitN(line, "#", 5)
		if len(parts) < 4 {
			continue
		}
		result = append(result, TermPosting{})
		posting := &result[len(result)-1]
		posting.Mfn, _ = strconv.Atoi(parts[0])
		posting.Tag, _ = strconv.Atoi(parts[1])
		posting.Occurrence, _ = strconv.Atoi(parts[2])
		posting.Count, _ = strconv.Atoi(parts[3])
		if len(parts) > 4 {
			posting.Text = parts[4]
		}
	}

	return
}

func (posting *TermPosting) String() string {
	return strconv.Itoa(posting.Mfn) + "#" +
		strconv.Itoa(posting.Tag) + "#" +
		strconv.Itoa(posting.Occurrence) + "#" +
		strconv.Itoa(posting.Count) + "#" +
		posting.Text
}
