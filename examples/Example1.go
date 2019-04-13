package main

import "../src/irbis"

func main() {
	// Подключаемся к серверу
	connection := irbis.NewConnection()
	connection.Host = "localhost"
	connection.Username = "librarian"
	connection.Password = "secret"
	if !connection.Connect() {
		println("Не удалось подключиться")
		return
	}

	// По выходу из функции произойдет отключение от сервера
	defer connection.Disconnect()

	// Общие сведения о сервере
	println("Версия сервера:", connection.ServerVersion)
	println("Интервал:", connection.Interval)

	// Из INI-файла можно получить настройки клиента
	ini := connection.Ini
	dbnnamecat := ini.GetValue("Main", "DBNNAMECAT", "???")
	println("DBNNAMECAT:", dbnnamecat)

	// Находим записи с автором "Пушкин"
	found := connection.Search("\"A=Пушкин$\"")
	println("Найдено:", len(found))

	// Ограничиваемся первыми 10 записями
	found = found[:10]

	for _, mfn := range found {
		// Считываем запись с сервера
		record := connection.ReadRecord(mfn)

		// Получаем значение поля/подполя
		title := record.FSM(200, 'a')
		println("Заглавие:", title)

		// Расформатируем запись на сервере
		description := connection.FormatMfn("@brief", mfn)
		println("Биб. описание:", description)
	}
}
