package main

import (
	"./irbis"
	"os"
	"sort"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		println("Usage: QuintEssence <mstfile> <textfile>")
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
	list := make([]string, 0, maxMfn)

	for mfn := 1; mfn < maxMfn; mfn++ {
		record, err := reader.ReadRecord(mfn)
		if err != nil || record == nil {
			continue
		}
		index := record.FM(903)
		countText := record.FM(999)
		if len(index) == 0 || len(countText) == 0 {
			continue
		}
		count, err := strconv.Atoi(countText)
		if count == 0 || err != nil {
			continue
		}
		line := index + "\t" + strconv.Itoa(count)
		list = append(list, line)
	}

	sort.Strings(list)
	for _, line := range list {
		_, err = writer.WriteString(line)
		if err != nil {
			panic(err)
		}
		_, err = writer.WriteString("\n")
		if err != nil {
			panic(err)
		}
	}
}
