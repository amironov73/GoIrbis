@echo off

SET GOEXECUTABLE=%GOROOT%\bin\go.exe
SET SOURCE=%CD%\src
SET OUTPUT=%CD%\bin

%GOEXECUTABLE% get -u -v golang.org/x/text/encoding/charmap
call :COMPILE SafeExperiments

goto :DONE

:COMPILE

%GOEXECUTABLE% build -o %OUTPUT%\%1.exe -v -x -i %SOURCE%\%1.go

:DONE
