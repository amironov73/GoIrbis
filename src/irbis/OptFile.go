package irbis

import (
	"regexp"
	"strconv"
	"strings"
)

// OptLine Строка OPT-файла.
type OptLine struct {
	// Pattern Паттерн.
	Pattern string

	// Worksheet Соответствующий рабочий лист.
	Worksheet string
}

func (opt *OptLine) Parse(text string) bool {
	rx := regexp.MustCompile(`\s+`)
	parts := rx.Split(text, 2)
	if len(parts) != 2 {
		return false
	}

	opt.Pattern = parts[0]
	opt.Worksheet = parts[1]
	return true
}

func (opt *OptLine) String() string {
	// TODO implement properly
	return opt.Pattern + " " + opt.Worksheet
}

// OptFile OPT-файл -- файл оптимизации рабочих листов и форматов отображения.
type OptFile struct {
	// WorksheetLength Длина рабочего листа.
	WorksheetLength int

	// WorksheetTag Метка поля рабочего листа.
	WorksheetTag int

	// Lines Строки с паттернами.
	Lines []OptLine
}

func NewOptFile() *OptFile {
	result := new(OptFile)
	result.WorksheetLength = 5
	result.WorksheetTag = 920

	return result
}

// GetWorksheet Поолучение рабочего листа записи.
func (opt *OptFile) GetWorksheet(record *MarcRecord) string {
	return record.FM(opt.WorksheetTag)
}

func (opt *OptFile) Parse(lines []string) {
	opt.WorksheetTag, _ = strconv.Atoi(lines[0])
	opt.WorksheetLength, _ = strconv.Atoi(lines[1])
	for _, line := range lines[2:] {
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "*****") {
			break
		}

		item := OptLine{}
		if item.Parse(line) {
			opt.Lines = append(opt.Lines, item)
		}
	}
}

func (opt *OptFile) String() string {
	result := strings.Builder{}
	result.WriteString(strconv.Itoa(opt.WorksheetTag))
	result.WriteString("\n")
	result.WriteString(strconv.Itoa(opt.WorksheetLength))
	result.WriteString("\n")

	for _, line := range opt.Lines {
		result.WriteString(line.String())
		result.WriteString("\n")
	}

	result.WriteString("*****\n")

	return result.String()
}
