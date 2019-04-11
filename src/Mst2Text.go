package main

import (
	"./irbis"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		println("Usage: Mst2Text <mstfile> <textfile>")
		return
	}

	input := os.Args[1]
	output := os.Args[2]
	reader, err := irbis.OpenDatabase(input)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	writer, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer func() { _ = writer.Close() }()

	maxMfn := reader.GetMaxMfn()
	for mfn := 1; mfn < maxMfn; mfn++ {
		record, err := reader.ReadRecord(mfn)
		if err != nil || record == nil {
			continue
		}
		err = record.ExportPlainText(writer)
		if err != nil {
			panic(err)
		}
	}
}
