set GOPATH=%CD%;%GOPATH%
echo %GOPATH%
mkdir bin

rem Install TDM-GCC first!  http://tdm-gcc.tdragon.net/download

echo package env > core/env/commit_log.go
echo const LastCommitLog ="N/A" >> core/env/commit_log.go
echo const BuildDate  ="N/A" >> core/env/commit_log.go

go build -o bin/gopa.exe