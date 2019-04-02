package irbis

import "strconv"

// SearchScenario Сценарий поиска.
type SearchScenario struct {
	// Name Название поискового атрибута
	// (автор, инвентарный номер и т. д.).
	Name string

	// Prefix Префикс соответствующих терминов
	// в словаре (может быть пустым).
	Prefix string

	// DictionaryType Тип словаря для соответствующего поиска.
	DictionaryType int

	// MenuName Имя файла справочника.
	MenuName string

	// OldFormat Имя формата (без расширения).
	OldFormat string

	// Correction Способ корректировки по словарю.
	Correction string

	// Truncation Исходное положение переключателя "Усечение".
	Truncation string

	// Hint Текст подсказки/предупреждения.
	Hint string

	// ModByDicAuto Параметр пока не задействован.
	ModByDicAuto string

	// Logic Применимые логические операторы.
	Logic string

	// Advance Правила автоматического расширения поиска
	// на основе авторитетного файла или тезауруса.
	Advance string

	// Format Имя формата показа документов.
	Format string
}

func (section *IniSection) get(name string, index int) string {
	fullName := "Item" + name + strconv.Itoa(index)
	return section.GetValue(fullName, "")
}

func (section *IniSection) getInt(name string, index int) int {
	value := section.get(name, index)
	result, _ := strconv.Atoi(value)
	return result
}

func ParseScenarios(ini *IniFile) (result []SearchScenario) {
	section := ini.FindSection("SEARCH")
	if section == nil {
		return
	}

	count, _ := strconv.Atoi(section.GetValue("ItemNumb", "0"))
	for i := 0; i < count; i++ {
		scenario := SearchScenario{}
		scenario.Name = section.get("Name", i)
		scenario.Prefix = section.get("Pref", i)
		scenario.DictionaryType = section.getInt("DictionType", i)
		scenario.MenuName = section.get("Menu", i)
		scenario.OldFormat = ""
		scenario.Correction = section.get("ModByDic", i)
		scenario.Truncation = section.get("Tranc", i)
		scenario.Hint = section.get("Hint", i)
		scenario.ModByDicAuto = section.get("ModByDicAuto", i)
		scenario.Logic = section.get("Logic", i)
		scenario.Advance = section.get("Adv", i)
		scenario.Format = section.get("Pft", i)
		result = append(result, scenario)
	}

	return
}
