package irbis

// DirectAccess осуществляет прямой доступ к базам данных.
type DirectAccess struct {
	mst      *MstFile
	xrf      *XrfFile
	ifp      *IfpFile
	filename string
}

// OpenDatabase открывает базу данных для чтения.
func OpenDatabase(filename string) (result *DirectAccess, err error) {
	var mst *MstFile
	mst, err = OpenMstFile(filename + ".mst")
	if err != nil {
		return
	}

	var xrf *XrfFile
	xrf, err = OpenXrfFile(filename + ".xrf")
	if err != nil {
		mst.Close()
		return
	}

	var ifp *IfpFile
	ifp, err = OpenIfpFile(filename)
	if err != nil {
		mst.Close()
		xrf.Close()
		return
	}

	result = new(DirectAccess)
	result.filename = filename
	result.mst = mst
	result.xrf = xrf
	result.ifp = ifp

	return
}

// Close закрывает базу данных.
func (access *DirectAccess) Close() {
	access.mst.Close()
	access.xrf.Close()
	access.ifp.Close()
}

// GetMaxMfn получает максимальный MFN для данной базы.
func (access *DirectAccess) GetMaxMfn() int {
	return int(access.mst.Control.NextMfn - 1)
}

// ReadRawRecord считывает запись в сыром виде
func (access *DirectAccess) ReadRawRecord(mfn int) (result *MstRecord, err error) {
	var xrf XrfRecord
	xrf, err = access.xrf.ReadRecord(mfn)
	if err != nil {
		return
	}

	result, err = access.mst.ReadRecord(xrf.Offset())
	if err != nil {
		return
	}

	return
}

// ReadRecord считывает и декодирует запись.
func (access *DirectAccess) ReadRecord(mfn int) (result *MarcRecord, err error) {
	var mst *MstRecord
	mst, err = access.ReadRawRecord(mfn)
	if err != nil {
		return
	}

	result = mst.Decode()
	return
}
