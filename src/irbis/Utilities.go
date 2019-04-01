package irbis

import (
	//	"golang.org/x/text/encoding/charmap"
	"strings"
	"unicode"
)

const IrbisDelimiter = "\x1F\x1E"
const ShortDelimiter = "\x1E"

func contains(s []int, item int) bool {
	for _, one := range s {
		if one == item {
			return true
		}
	}

	return false
}

func toWin1251(text string) []byte {
	return cp1251FromUnicode(text)

	//encoder := charmap.Windows1251.NewEncoder()
	//result, _ := encoder.Bytes([]byte(text))
	//return result
}

func fromWin1251(buffer []byte) string {
	return cp1251ToUnicode(buffer)

	//decoder := charmap.Windows1251.NewDecoder()
	//temp, _ := decoder.Bytes(buffer)
	//result := string(temp)
	//return result
}

func toUtf8(text string) []byte {
	result := []byte(text)
	return result
}

func fromUtf8(buffer []byte) string {
	return string(buffer)
}

func removeComments(text string) string {
	if len(text) == 0 || !strings.Contains(text, "/*") {
		return text
	}

	result := strings.Builder{}
	state := '\x00'
	chars := []rune(text)
	index := 0
	length := len(chars)
	result.Grow(length)

	for index < length {
		c := chars[index]

		switch state {
		case '\'', '"', '|':
			if c == state {
				state = '\x00'
			}
			result.WriteRune(c)

		default:
			if c == '/' {
				if (index+1 < length) && (chars[index+1] == '*') {
					for index < length {
						c = chars[index]
						if (c == '\r') || (c == '\n') {
							result.WriteRune(c)
							break
						}
						index++
					}
				} else {
					result.WriteRune(c)
				}
			} else if (c == '\'') || (c == '"') || (c == '|') {
				state = c
				result.WriteRune(c)
			} else {
				result.WriteRune(c)
			}
		}

		index++
	}

	return result.String()
}

func prepareFormat(text string) string {
	text = removeComments(text)
	length := len(text)
	if length == 0 {
		return text
	}

	flag := false
	chars := []rune(text)
	for i := range chars {
		if chars[i] < ' ' {
			flag = true
			break
		}
	}

	if !flag {
		return text
	}

	result := strings.Builder{}
	result.Grow(length)
	for i := range chars {
		c := chars[i]
		if c >= ' ' {
			result.WriteRune(c)
		}
	}

	return result.String()
}

func IrbisToDos(text string) string {
	return strings.ReplaceAll(text, IrbisDelimiter, "\n")
}

func IrbisToLines(text string) []string {
	return strings.Split(text, IrbisDelimiter)
}

func LinesToIrbis(lines []string) string {
	result := strings.Builder{}
	for _, line := range lines {
		result.WriteString(line)
		result.WriteString(IrbisDelimiter)
	}

	return result.String()
}

func PickOne(lines ...string) string {
	for _, line := range lines {
		if len(line) != 0 {
			return line
		}
	}

	return ""
}

func SameRune(left, right rune) bool {
	return unicode.ToUpper(left) == unicode.ToUpper(right)
}

func SameString(left, right string) bool {
	return strings.EqualFold(left, right)
}

func SplitLines(text string) []string {
	// TODO implement properly
	return strings.Split(text, "\n")
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func DescribeError(code int) string {
	if code >= 0 {
		return "Нет ошибки"
	}

	switch code {
	case -100:
		return "Заданный MFN вне пределов БД"
	case -101:
		return "Ошибочный размер полки"
	case -102:
		return "Ошибочный номер полки"
	case -140:
		return "MFN вне пределов БД"
	case -141:
		return "Ошибка чтения"
	case -200:
		return "Указанное поле отсутствует"
	case -201:
		return "Предыдущая версия записи отсутствует"
	case -202:
		return "Заданный термин не найден (термин не существует)"
	case -203:
		return "Последний термин в списке"
	case -204:
		return "Первый термин в списке"
	case -300:
		return "База данных монопольно заблокирована"
	case -301:
		return "База данных монопольно заблокирована"
	case -400:
		return "Ошибка при открытии файлов MST или XRF (ошибка файла данных)"
	case -401:
		return "Ошибка при открытии файлов IFP (ошибка файла индекса)"
	case -402:
		return "Ошибка при записи"
	case -403:
		return "Ошибка при актуализации"
	case -600:
		return "Запись логически удалена"
	case -601:
		return "Запись физически удалена"
	case -602:
		return "Запись заблокирована на ввод"
	case -603:
		return "Запись логически удалена"
	case -605:
		return "Запись физически удалена"
	case -607:
		return "Ошибка autoin.gbl"
	case -608:
		return "Ошибка версии записи"
	case -700:
		return "Ошибка создания резервной копии"
	case -701:
		return "Ошибка восстановления из резервной копии"
	case -702:
		return "Ошибка сортировки"
	case -703:
		return "Ошибочный термин"
	case -704:
		return "Ошибка создания словаря"
	case -705:
		return "Ошибка загрузки словаря"
	case -800:
		return "Ошибка в параметрах глобальной корректировки"
	case -801:
		return "ERR_GBL_REP"
	case -802:
		return "ERR_GBL_MET"
	case -1111:
		return "Ошибка исполнения сервера (SERVER_EXECUTE_ERROR)"
	case -2222:
		return "Ошибка в протоколе (WRONG_PROTOCOL)"
	case -3333:
		return "Незарегистрированный клиент (ошибка входа на сервер) (клиент не в списке)"
	case -3334:
		return "Клиент не выполнил вход на сервер (клиент не используется)"
	case -3335:
		return "Неправильный уникальный идентификатор клиента"
	case -3336:
		return "Нет доступа к командам АРМ"
	case -3337:
		return "Клиент уже зарегистрирован"
	case -3338:
		return "Недопустимый клиент"
	case -4444:
		return "Неверный пароль"
	case -5555:
		return "Файл не существует"
	case -6666:
		return "Сервер перегружен. Достигнуто максимальное число потоков обработки"
	case -7777:
		return "Не удалось запустить/прервать поток администратора (ошибка процесса)"
	case -8888:
		return "Общая ошибка"
	}

	return "Неизвестная ошибка"
}
