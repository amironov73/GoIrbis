package main

import "fmt"

type SubField struct {
	Code  string
	Value string
}

type RecordField struct {
	Tag       int
	Value     string
	Subfields []SubField
}

type MarcRecord struct {
	Database string
	Mfn      int
	Version  int
	Status   int
	Fields   []RecordField
}

type ClientQuery struct {
	Dummy int
}

type ServerResponse struct {
	Dummy int
}

type IrbisConnection struct {
	Host        string
	Port        int
	Username    string
	Password    string
	Database    string
	Workstation string
	ClientId    int
	QueryId     int
}

func (connection *IrbisConnection) Connect() {
	fmt.Println("Connecting")
}

func (connection *IrbisConnection) Disconnect() {
	fmt.Println("Disconnecting")
}

func (connection *IrbisConnection) Execute(query *ClientQuery) ServerResponse {
	return ServerResponse{}
}

func (connection *IrbisConnection) GetMaxMfn(database string) int {
	return 0
}

func (connection *IrbisConnection) ReadRecord(mfn int) MarcRecord {
	return MarcRecord{}
}

func main () {
	var connection = IrbisConnection{}
	connection.Connect()
	var maxMfn = connection.GetMaxMfn("IBIS")
	fmt.Println("Max MFN", maxMfn)
	var record = connection.ReadRecord(123);
	fmt.Println(record)
	connection.Disconnect()
}
