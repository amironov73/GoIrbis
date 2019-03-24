@echo off

SET MYGOROOT=C:\Go
SET MYGOBIN=%MYGOROOT%\bin
SET MYGOEXE=%MYGOBIN%\go.exe
SET SOURCE=%CD%\Source\goirbis.go
SET OUTPUT=%CD%\Binaries\goirbis.exe

%MYGOEXE% build -o %OUTPUT% %SOURCE%
