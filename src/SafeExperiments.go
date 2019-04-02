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

	ini := connection.Ini
	dbnnamecat := ini.GetValue("Main", "DBNNAMECAT", "???")
	fmt.Println("DBNNAMECAT:", dbnnamecat)

	fmt.Println("Server version: ", connection.ServerVersion)
	fmt.Println("Interval: ", connection.Interval)

	version := connection.GetServerVersion()
	fmt.Println("Organization:", version.Organization)

	connection.NoOp()

	processes := connection.ListProcesses()
	fmt.Println("Processes:", processes)

	databases := connection.ListDatabases("")
	fmt.Println("Databases:", databases)

	dbInfo := connection.GetDatabaseInfo("IBIS")
	fmt.Println("Deleted records:", len(dbInfo.LogicallyDeletedRecords)+
		len(dbInfo.PhysicallyDeletedRecords))

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

	manyFormatted := connection.FormatRecords("@brief", []int{1, 2, 3})
	fmt.Println(manyFormatted)

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

	languages := connection.ListTerms("J=")
	fmt.Println(languages)

	postingParameters := irbis.NewPostingParameters()
	postingParameters.Term = "J=CHI"
	postingParameters.NumberOfPostings = 100
	postings := connection.ReadPostings(postingParameters)
	fmt.Println(postings)

	stat := connection.GetServerStat()
	fmt.Println(stat)

	files := connection.ListFiles("3.IBIS.brief.*", "3.IBIS.a*.pft")
	fmt.Println(files)

	parFile := connection.ReadParFile("1..IBIS.PAR")
	fmt.Println(parFile)

	optFile := connection.ReadOptFile("3.IBIS.WS31.OPT")
	fmt.Println(optFile)

	records := connection.ReadRecords([]int{1, 2, 3})
	fmt.Println(records)

	tree := connection.ReadTreeFile("3.IBIS.II.TRE")
	fmt.Println(tree)

	single := connection.SearchSingleRecord(`"I=65.304.13-772296"`)
	fmt.Println(single)
}
