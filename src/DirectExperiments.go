package main

import (
	"./irbis"
	"fmt"
)

func readAndShowRecord(access *irbis.DirectAccess, mfn int) {
	record, err := access.ReadRecord(mfn)
	if err != nil {
		panic(err)
	}

	fmt.Println(record.String())
}

func main() {
	access, err := irbis.OpenDatabase("data/irbis64/datai/ibis/ibis")
	if err != nil {
		panic(err)
	}
	defer access.Close()

	fmt.Println("Max MFN=", access.GetMaxMfn())

	for mfn := 1; mfn <= 10; mfn++ {
		readAndShowRecord(access, mfn)
	}
}
