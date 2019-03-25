@echo off

SET GOEXECUTABLE=%GOROOT%\bin\go.exe
SET SOURCE=%CD%\Source\goirbis.go
SET OUTPUT=%CD%\Binaries\goirbis.exe

%GOEXECUTABLE% build -o %OUTPUT% -v -x -i %SOURCE%
