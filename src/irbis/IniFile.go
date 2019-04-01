package irbis

import (
	"strings"
)

type IniLine struct {
	Key   string
	Value string
}

func (line *IniLine) String() string {
	return line.Key + "=" + line.Value
}

type IniSection struct {
	Name  string
	Lines []IniLine
}

func (section *IniSection) Find(key string) *IniLine {
	for i := range section.Lines {
		if SameString(section.Lines[i].Key, key) {
			return &section.Lines[i]
		}
	}
	return nil
}

func (section *IniSection) GetValue(key, defaultValue string) string {
	line := section.Find(key)
	if line == nil {
		return defaultValue
	}
	return line.Value
}

func (section *IniSection) Remove(key string) {
	// TODO implement
}

func (section *IniSection) SetValue(key, value string) {
	if len(value) == 0 {
		section.Remove(key)
	} else {
		item := section.Find(key)
		if item != nil {
			item.Value = value
		} else {
			item = new(IniLine)
			section.Lines = append(section.Lines, *item)
			item = &section.Lines[len(section.Lines)-1]
			item.Key = key
			item.Value = value
		}
	}
}

func (section *IniSection) String() string {
	var result strings.Builder
	if len(section.Name) != 0 {
		result.WriteString("[")
		result.WriteString(section.Name)
		result.WriteString("]")
		result.WriteString("\n")
	}
	for i := range section.Lines {
		result.WriteString(section.Lines[i].String())
		result.WriteString("\n")
	}
	return result.String()
}

type IniFile struct {
	Sections []IniSection
}

func NewIniFile() *IniFile {
	return &IniFile{}
}

func (file *IniFile) FindSection(name string) *IniSection {
	for i := range file.Sections {
		section := &file.Sections[i]
		if SameString(section.Name, name) {
			return section
		}
	}
	return nil
}

func (file *IniFile) GetOrCreateSection(name string) *IniSection {
	result := file.FindSection(name)
	if result == nil {
		result = new(IniSection)
		file.Sections = append(file.Sections, *result)
		result = &file.Sections[len(file.Sections)-1]
		result.Name = name
	}
	return result
}

func (file *IniFile) GetValue(sectionName, key, defaultValue string) string {
	section := file.FindSection(sectionName)
	if section != nil {
		return section.GetValue(key, defaultValue)
	}
	return defaultValue
}

func (file *IniFile) Parse(lines []string) {
	var section *IniSection = nil
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if line[0] == '[' {
			name := line[1 : len(line)-1]
			section = file.GetOrCreateSection(name)
		} else if section != nil {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				item := IniLine{parts[0], parts[1]}
				section.Lines = append(section.Lines, item)
			}
		}
	}
}

func (file *IniFile) SetValue(sectionName, key, value string) {
	section := file.GetOrCreateSection(sectionName)
	section.SetValue(key, value)
}

func (file *IniFile) String() string {
	result := strings.Builder{}
	first := true
	for i := range file.Sections {
		if !first {
			result.WriteString("\n")
		}
		result.WriteString(file.Sections[i].String())
		first = false
	}
	return result.String()
}
