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

	ini := connection.IniFile
	dbnnamecat := ini.GetValue("Main", "DBNNAMECAT", "???")
	fmt.Println("DBNNAMECAT:", dbnnamecat)

	fmt.Println("Server version: ", connection.ServerVersion)
	fmt.Println("Interval: ", connection.Interval)

	version := connection.GetServerVersion()
	fmt.Println("Organization:", version.Organization)

	connection.NoOp()

	users := connection.GetUserList()
	fmt.Println(users)

	maxMfn := connection.GetMaxMfn("IBIS")
	fmt.Println("Max MFN", maxMfn)

	formatted := connection.FormatMfn("@brief", 123)
	fmt.Println(formatted)

	mfn := 123
	format := "'Ἀριστοτέλης: ', v200^a"
	formatted = connection.FormatMfnUtf(format, mfn)
	fmt.Println(formatted)

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

	menu := connection.ReadMenuFile("3.IBIS.FORMATW.MNU")
	fmt.Println(menu.String())
}
