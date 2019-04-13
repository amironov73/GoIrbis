package main

import (
	"../src/irbis"
	"fmt"
	"strconv"
)

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

	// Записи будут помещаться в базу SANDBOX
	connection.Database = "SANDBOX"

	for i:=0; i < 10; i++ {
		// Создаём базу в памяти клиента
		record := irbis.NewMarcRecord()

		// Наполняем ее полями: первый автор (поле с подполями),
		record.Add(700, "").
			Add('a', "Миронов").
			Add('b', "А. В.").
			Add('g', "Алексей Владимирович")

		// заглавие (поле с подполями)
		record.Add(200, "").
			Add('a', "Работа с ИРБИС64: версия " +
				strconv.Itoa(i) ).
			Add('e', "руководство пользователя")

		// выходные данные (поле с подполями)
		record.Add(210, "").
			Add('a', "Иркутск").
			Add('c', "ИРНИТУ").
			Add('d', "2019")

		// рабочий лист (поле без подполей)
		record.Add(920, "PAZK")

		// Отсылаем запись на сервер.
		// Обратно приходит запись,
		// обработанная AUTOIN.GBL
		connection.WriteRecord(record)

		fmt.Println(record)
	}
}

