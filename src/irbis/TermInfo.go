package irbis

import (
	"strconv"
	"strings"
)

// TermInfo Информация о термине поискового словаря.
type TermInfo struct {
	// Count Количество ссылок.
	Count int

	// Text Поисковый термин.
	Text string
}

// ParseTerms Разбор ответа сервера, содержащего массив терминов.
func ParseTerms(lines []string) (result []TermInfo) {
	for _, line := range lines {
		if len(line) != 0 {
			parts := strings.SplitN(line, "#", 2)
			if len(parts) == 2 {
				result = append(result, TermInfo{})
				term := &result[len(result)-1]
				term.Count, _ = strconv.Atoi(parts[0])
				term.Text = parts[1]
			}
		}
	}
	return
}

func (term *TermInfo) String() string {
	return strconv.Itoa(term.Count) + "#" + term.Text
}
