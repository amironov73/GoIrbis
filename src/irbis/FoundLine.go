package irbis

import (
	"strconv"
	"strings"
)

type FoundLine struct {
	Mfn         int
	Description string
}

func (line *FoundLine) parse(text string) {
	parts := strings.SplitN(text, "#", 2)
	line.Mfn, _ = strconv.Atoi(parts[0])
	if len(parts) > 1 {
		line.Description = parts[1]
	}
}

func parseFoundLines(lines []string) []FoundLine {
	result := []FoundLine{}
	for _, line := range lines {
		if len(line) != 0 {
			item := FoundLine{}
			item.parse(line)
			if item.Mfn != 0 {
				result = append(result, item)
			}
		}
	}

	return result
}

func parseFoundMfn(lines []string) []int {
	result := []int{}
	for _, line := range lines {
		if len(line) != 0 {
			mfn, _ := strconv.Atoi(line)
			if mfn != 0 {
				result = append(result, mfn)
			}
		}
	}

	return result
}

func parseFoundDescriptions(lines []string) []string {
	result := []string{}
	for _, line := range lines {
		if len(line) != 0 {
			parts := strings.SplitN(line, "#", 2)
			if len(parts) > 1 {
				text := parts[1]
				if len(text) != 0 {
					result = append(result, text)
				}
			}
		}
	}

	return result
}
