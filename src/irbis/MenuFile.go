package irbis

import "strings"

type MenuEntry struct {
	Code    string
	Comment string
}

func (entry *MenuEntry) String() string {
	return entry.Code + " - " + entry.Comment
}

type MenuFile struct {
	Entries []MenuEntry
}

func (menu *MenuFile) Add(code, comment string) *MenuFile {
	entry := MenuEntry{code, comment}
	menu.Entries = append(menu.Entries, entry)
	return menu
}

func (menu *MenuFile) Clear() *MenuFile {
	menu.Entries = []MenuEntry{}
	return menu
}

func (menu *MenuFile) GetEntry(code string) *MenuEntry {
	if menu.Entries == nil {
		return nil
	}

	for i := range menu.Entries {
		entry := &menu.Entries[i]
		if SameString(entry.Code, code) {
			return entry
		}
	}

	code = strings.TrimSpace(code)
	for i := range menu.Entries {
		entry := &menu.Entries[i]
		if SameString(entry.Code, code) {
			return entry
		}
	}

	code = strings.Trim(code, "-=:")
	for i := range menu.Entries {
		entry := &menu.Entries[i]
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
	menu.Entries = make([]MenuEntry, 0, len(lines)/2)
	length := len(lines)
	for i := 0; i < length; i += 2 {
		code := lines[i]
		if len(code) == 0 || strings.HasPrefix(code, "*****") {
			break
		}
		comment := lines[i+1]
		entry := MenuEntry{code, comment}
		menu.Entries = append(menu.Entries, entry)
	}
}

func (menu *MenuFile) String() string {
	var result strings.Builder
	for _, entry := range menu.Entries {
		result.WriteString(entry.String() + "\n")
	}
	result.WriteString("*****\n")

	return result.String()
}
