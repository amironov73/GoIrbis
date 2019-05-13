package main

import (
	"./irbis"
	"os"
	"time"
)

/*
    Простая программа, удаляющая из базы данных RQST все выполненные заказы
    (для уменьшения нагрузки на сеть и сервер со стороны АРМ "Книговыдачи").
*/

func main() {
	start := time.Now()

	if len(os.Args) != 2 {
		println("USAGE: RqstShrink <connectionString>")
		return
	}

	connectionString := os.Args[1]
	connection := irbis.NewConnection()
	connection.ParseConnectionString(connectionString)
	if !connection.Connect() {
		println("Can't connect")
		return
	}

	defer connection.Disconnect()

	if connection.Workstation != irbis.ADMINISTRATOR {
		println("Not administrator! Exting")
		return
	}

	maxMfn := connection.GetMaxMfn(connection.Database)
	expression := `"I=0" + "I=2"` // Невыполненные и зарезервированные
	found := connection.SearchAll(expression)
	if len(found) == maxMfn {
		println("No truncation needed, exiting")
	}

	// TODO ReadAllRecords
	goodRecords := connection.ReadRecords(found)
	println("Good records loaded:", len(goodRecords))
	for i := range goodRecords {
		record := &goodRecords[i]
		record.Reset()
		record.Database = connection.Database
	}

	connection.TruncateDatabase(connection.Database)
	if connection.GetMaxMfn(connection.Database) > 1 {
		println("Error while truncating the database, exiting")
		return
	}

	// TODO WriteAllRecords
	connection.WriteRecords(goodRecords)
	println("Good records restored")

	elapsed := time.Since(start)
	println("Elapsed time:", elapsed.String())
}
