@echo off

SET GOEXECUTABLE=%GOROOT%\bin\go.exe
SET GOPATH=%GOPATH%;%CD%;
SET SOURCE=%CD%\src
SET OUTPUT=%CD%\bin

rem %GOEXECUTABLE% get -u -v golang.org/x/text/encoding/charmap
call :COMPILE SafeExperiments
call :COMPILE DirectExperiments

goto :DONE

:COMPILE

%GOEXECUTABLE% build -o %OUTPUT%\%1.exe -v -x -i %SOURCE%\%1.go

:DONE
