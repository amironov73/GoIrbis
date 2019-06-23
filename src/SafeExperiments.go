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

	fmt.Println("NoOp")
	connection.NoOp()

	fmt.Println("ListProcesses")
	processes := connection.ListProcesses()
	fmt.Println("Processes:", processes)

	fmt.Println("ListDatabases")
	databases := connection.ListDatabases("")
	fmt.Println("Databases:", databases)

	fmt.Println("GetDatabaseInfo")
	dbInfo := connection.GetDatabaseInfo("IBIS")
	fmt.Println("Deleted records:", len(dbInfo.LogicallyDeletedRecords)+
		len(dbInfo.PhysicallyDeletedRecords))

	fmt.Println("GetUserList")
	users := connection.GetUserList()
	fmt.Println(users)

	fmt.Println("GetMaxMfn")
	maxMfn := connection.GetMaxMfn("IBIS")
	fmt.Println("Max MFN", maxMfn)

	fmt.Println("FormatMfn")
	formatted := connection.FormatMfn("@brief", 123)
	fmt.Println(formatted)

	fmt.Println("FormatMfn")
	mfn := 123
	format := "'Ἀριστοτέλης: ', v200^a"
	formatted = connection.FormatMfn(format, mfn)
	fmt.Println(formatted)

	fmt.Println("FormatRecords")
	manyFormatted := connection.FormatRecords("@brief", []int{1, 2, 3})
	fmt.Println(manyFormatted)

	fmt.Println("ReadRecord")
	record := connection.ReadRecord(123)
	fmt.Println(record.Encode("\n"))

	fmt.Println("ReadTextFile")
	content := connection.ReadTextFile("3.IBIS.WS.OPT")
	fmt.Println(content)

	fmt.Println("SearchCount")
	count := connection.SearchCount("\"A=ПУШКИН$\"")
	fmt.Println("COUNT:", count)

	fmt.Println("Search")
	found := connection.Search("\"A=ПУШКИН$\"")
	for _, mfn := range found {
		fmt.Print(",", mfn)
	}
	fmt.Println()

	fmt.Println("ReadMenuFile")
	menu := connection.ReadMenuFile("3.IBIS.FORMATW.MNU")
	fmt.Println(menu.String())

	fmt.Println("ListTerms")
	languages := connection.ListTerms("J=")
	fmt.Println(languages)

	fmt.Println("ReadPostings")
	postingParameters := irbis.NewPostingParameters()
	postingParameters.Term = "J=CHI"
	postingParameters.NumberOfPostings = 100
	postings := connection.ReadPostings(postingParameters)
	fmt.Println(postings)

	fmt.Println("GetRecordPostings")
	postings = connection.GetRecordPostings(2, "A=$")
	fmt.Println(postings)

	fmt.Println("GetServerStat")
	stat := connection.GetServerStat()
	fmt.Println(stat)

	fmt.Println("ListFiles")
	files := connection.ListFiles("3.IBIS.brief.*", "3.IBIS.a*.pft")
	fmt.Println(files)

	fmt.Println("ReadParFile")
	parFile := connection.ReadParFile("1..IBIS.PAR")
	fmt.Println(parFile)

	fmt.Println("ReadOptFile")
	optFile := connection.ReadOptFile("3.IBIS.WS31.OPT")
	fmt.Println(optFile)

	fmt.Println("ReadRecords")
	records := connection.ReadRecords([]int{1, 2, 3})
	fmt.Println(records)

	fmt.Println("ReadTreeFile")
	tree := connection.ReadTreeFile("3.IBIS.II.TRE")
	fmt.Println(tree)

	fmt.Println("SearchSingleRecord")
	single := connection.SearchSingleRecord(`"I=65.304.13-772296"`)
	fmt.Println(single)

	fmt.Println("THAT'S ALL FOLKS!")
}
