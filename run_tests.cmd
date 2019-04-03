@echo off

SET GOEXECUTABLE=%GOROOT%\bin\go.exe
SET GOPATH=%GOPATH%;%CD%
SET SOURCE=%CD%\src
SET OUTPUT=%CD%\bin

go test -v irbis
