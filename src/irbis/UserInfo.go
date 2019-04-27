package irbis

import "strconv"

// Информация о зарегистрированном пользователе системы
// (по данным client_m.mnu).
type UserInfo struct {
	Number        string // номер по порядку в списке.
	Name          string // логин.
	Password      string // пароль.
	Cataloger     string // доступность АРМ "Каталогизатор".
	Reader        string // доступность АРМ "Читатель".
	Circulation   string // доступность АРМ "Книговыдача".
	Acquisitions  string // доступность АРМ "Комплектатор".
	Provision     string // доступность АРМ "Книгообеспеченность".
	Administrator string // доступность АРМ "Администратор".
}

func formatPair(prefix, value, defaultValue string) string {
	if SameString(value, defaultValue) {
		return ""
	}
	return prefix + "=" + value + ";"
}

// Encode формирует строковое представление пользователя.
func (user *UserInfo) Encode() string {
	return user.Name + "\r\n" +
		user.Password + "\r\n" +
		formatPair("C", user.Cataloger, "irbisc.ini") +
		formatPair("R", user.Reader, "irbisr.ini") +
		formatPair("B", user.Circulation, "irbisb.ini") +
		formatPair("M", user.Acquisitions, "irbism.ini") +
		formatPair("K", user.Provision, "irbisk.ini") +
		formatPair("A", user.Administrator, "irbisa.ini")
}

// Разбор ответа сервера.
func parseUsers(lines []string) (result []UserInfo) {
	userCount, _ := strconv.Atoi(lines[0])
	linesPerUser, _ := strconv.Atoi(lines[1])
	if userCount == 0 || linesPerUser == 0 {
		return
	}
	lines = lines[2:]
	for i := 0; i < userCount; i++ {
		if len(lines) < linesPerUser {
			break
		}
		user := UserInfo{
			Number:        lines[0],
			Name:          lines[1],
			Password:      lines[2],
			Cataloger:     lines[3],
			Reader:        lines[4],
			Circulation:   lines[5],
			Acquisitions:  lines[6],
			Provision:     lines[7],
			Administrator: lines[8],
		}
		result = append(result, user)
		lines = lines[linesPerUser+1:]
	}
	return
}

func (user *UserInfo) String() string {
	return user.Name
}
