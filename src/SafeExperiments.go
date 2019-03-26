package main

import (
	"./irbis"
	"fmt"
)

func main() {
	connection := irbis.NewConnection()
	connection.Username = "librarian"
	connection.Password = "secret"
	if !connection.Connect() {
		fmt.Println("Can't connect")
		return
	}

	defer connection.Disconnect()

	fmt.Println("Server version: ", connection.ServerVersion)
	fmt.Println("Interval: ", connection.Interval)

	connection.NoOp()

	maxMfn := connection.GetMaxMfn("IBIS")
	fmt.Println("Max MFN", maxMfn)
	record := connection.ReadRecord(123)
	fmt.Println(record.Encode("\n"))

	content := connection.ReadTextFile("3.IBIS.WS.OPT")
	fmt.Println(content)

	count := connection.SearchCount("\"A=ПУШКИН$\"")
	fmt.Println("COUNT:", count)

	found := connection.Search("\"A=ПУШКИН$\"")
	for _, mfn := range found {
		fmt.Print(",", mfn)
	}
	fmt.Println()
}
