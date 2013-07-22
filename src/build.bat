set GOPATH=%CD%\..\;%CD%;%GOPATH%
echo %GOPATH%
mkdir bin
go build -o bin/gopa.exe