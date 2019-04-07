package main

import (
	"./irbis"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("data/test1.iso")
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()

	for mfn := 1; mfn <= 81; mfn++ {
		record := irbis.ReadIsoRecord(file, irbis.FromAnsi)
		record.Mfn = mfn
		fmt.Println(record)
	}
}
