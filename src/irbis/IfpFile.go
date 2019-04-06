package irbis

import "os"

const NodeLength = 2048
const MaxTermSize = 255

type TermLink struct {
	Mfn        int32
	Tag        int32
	Occurrence int32
	Index      int32
}

type IfpControlRecord struct {
	NextOffsetLow  int32
	NextOffsetHigh int32
	NodeBlockCount int32
	LeafBlockCount int32
	Reserved       int32
}

type IfpRecordLeader struct {
	LowOffset      int32
	HighOffset     int32
	TotalLinkCount int32
	BlockLinkCount int32
	Capacity       int32
}

type IfpRecord struct {
	Leader IfpRecordLeader
	Links  []TermLink
}

type IfpFile struct {
	ifpFile *os.File
	l01File *os.File
	n01File *os.File
}

// OpenIfpFile открывает файлы IFP, L01, N01
func OpenIfpFile(filename string) (result *IfpFile, err error) {
	var ifp *os.File
	ifp, err = os.Open(filename + ".ifp")
	if err != nil {
		return
	}

	var l01 *os.File
	l01, err = os.Open(filename + ".l01")
	if err != nil {
		_ = ifp.Close()
		return
	}

	var n01 *os.File
	n01, err = os.Open(filename + ".n01")
	if err != nil {
		_ = ifp.Close()
		_ = l01.Close()
		return
	}

	result = new(IfpFile)
	result.ifpFile = ifp
	result.l01File = l01
	result.n01File = n01

	return
}

// Close закрывает файлы IFP, L01, N01.
func (ifp *IfpFile) Close() {
	_ = ifp.ifpFile.Close()
	_ = ifp.l01File.Close()
	_ = ifp.n01File.Close()
}
