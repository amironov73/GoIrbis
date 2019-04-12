@echo off

rem SET GOEXECUTABLE=%GOROOT%\bin\go.exe
rem SET GOPATH=%GOPATH%;%CD%;
SET SOURCE=%CD%\src
SET OUTPUT=%CD%\bin

go get -u github.com/bradleyfalzon/apicompat/cmd/apicompat

cd src\irbis

apicompat