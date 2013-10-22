set GOPATH=%CD%\..\;%CD%;%GOPATH%
echo %GOPATH%
mkdir bin
go get github.com/zeebo/sbloom
go get github.com/cihub/seelog
go get github.com/robfig/config
go get github.com/PuerkitoBio/purell
go build -o bin/gopa.exe