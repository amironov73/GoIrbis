package main

import (
	"fmt"
	"irbis"
)

func readAndShowXrfRecord(xrf *irbis.XrfFile, mfn int) {
	record, err := xrf.ReadRecord(mfn)
	if err != nil {
		panic(err)
	}
	fmt.Println(mfn, record, record.Offset())
}

func main() {
	xrf, err := irbis.OpenXrfFile("data/irbis64/datai/ibis/ibis.xrf")
	if err != nil {
		panic(err)
	}
	defer xrf.Close()

	for mfn := 1; mfn < 10; mfn++ {
		readAndShowXrfRecord(xrf, mfn)
	}
}
