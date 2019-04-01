package irbis

import (
	"strconv"
	"strings"
)

// ParFile PAR-файл -- содержит пути к файлам базы данных ИРБИС.
type ParFile struct {
	// Xrf Путь к файлу XRF.
	Xrf string

	// Mst Путь к файлу MST.
	Mst string

	// Cnt Путь к файлу CNT.
	Cnt string

	// N01 Путь к файлу N01.
	N01 string

	// N02 В ИРБИС64 не используется.
	N02 string

	// L01 Путь к файлу L01.
	L01 string

	// L02 В ИРБИС64 не используется.
	L02 string

	// Ifp Путь к файлу IFP.
	Ifp string

	// Any Путь к файлу ANY.
	Any string

	// Pft Путь к PFT-файлам.
	Pft string

	// Ext Расположение внешних объектов (поле 951).
	// Параметр появился в версии 2012.
	Ext string
}

func NewParFile(mst string) *ParFile {
	result := new(ParFile)
	result.Mst = mst
	result.Xrf = mst
	result.Cnt = mst
	result.N01 = mst
	result.N02 = mst
	result.L01 = mst
	result.L02 = mst
	result.Ifp = mst
	result.Any = mst
	result.Pft = mst
	result.Ext = mst

	return result
}

// Parse Разбор ответа сервера.
func (par *ParFile) Parse(lines []string) {
	m := make(map[int]string)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		value := strings.TrimSpace(parts[1])
		key, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			continue
		}

		m[key] = value
	}

	par.Xrf = m[1]
	par.Mst = m[2]
	par.Cnt = m[3]
	par.N01 = m[4]
	par.N02 = m[5]
	par.L01 = m[6]
	par.L02 = m[7]
	par.Ifp = m[8]
	par.Any = m[9]
	par.Pft = m[10]
	par.Ext = m[11]
}

func (par *ParFile) String() string {
	return "1=" + par.Xrf + "\n" +
		"2=" + par.Mst + "\n" +
		"3=" + par.Cnt + "\n" +
		"4=" + par.N01 + "\n" +
		"5=" + par.N02 + "\n" +
		"6=" + par.L01 + "\n" +
		"7=" + par.L02 + "\n" +
		"8=" + par.Ifp + "\n" +
		"9=" + par.Any + "\n" +
		"10=" + par.Pft + "\n" +
		"11=" + par.Ext + "\n"
}
