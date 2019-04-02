package irbis

import (
	"strings"
)

// MenuEntry представляет собой пару строк в MNU-файле.
type MenuEntry struct {
	Code    string // Условный код
	Comment string // Комментарий к коду
}

func NewMenuEntry(code, comment string) *MenuEntry {
	result := new(MenuEntry)
	result.Code = code
	result.Comment = comment
	return result
}

func (entry *MenuEntry) String() string {
	return entry.Code + " - " + entry.Comment
}

type MenuFile struct {
	Entries []*MenuEntry
}

func (menu *MenuFile) Add(code, comment string) *MenuFile {
	entry := NewMenuEntry(code, comment)
	menu.Entries = append(menu.Entries, entry)
	return menu
}

func (menu *MenuFile) Clear() *MenuFile {
	menu.Entries = []*MenuEntry{}
	return menu
}

func (menu *MenuFile) GetEntry(code string) *MenuEntry {
	if menu.Entries == nil {
		return nil
	}

	for _, entry := range menu.Entries {
		if SameString(entry.Code, code) {
			return entry
		}
	}

	code = strings.TrimSpace(code)
	for _, entry := range menu.Entries {
		if SameString(entry.Code, code) {
			return entry
		}
	}

	code = strings.Trim(code, "-=:")
	for _, entry := range menu.Entries {
		if SameString(entry.Code, code) {
			return entry
		}
	}

	return nil
}

func (menu *MenuFile) GetValue(code, defaultValue string) string {
	entry := menu.GetEntry(code)
	if entry == nil {
		return defaultValue
	}
	return entry.Comment
}

func (menu *MenuFile) Parse(lines []string) {
	menu.Entries = make([]*MenuEntry, 0, len(lines)/2)
	length := len(lines)
	for i := 0; i < length; i += 2 {
		code := lines[i]
		if len(code) == 0 || strings.HasPrefix(code, "*****") {
			break
		}
		comment := lines[i+1]
		entry := NewMenuEntry(code, comment)
		menu.Entries = append(menu.Entries, entry)
	}
}

func (menu *MenuFile) String() string {
	result := strings.Builder{}
	for _, entry := range menu.Entries {
		result.WriteString(entry.String() + "\n")
	}
	result.WriteString("*****\n")

	return result.String()
}
