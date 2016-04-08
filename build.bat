set GOPATH=%CD%;%GOPATH%
echo %GOPATH%
mkdir bin
go get github.com/cihub/seelog
go get github.com/robfig/config
go get github.com/PuerkitoBio/purell
go get code.google.com/p/go.net/websocket
go get github.com/errplane/errplane-go
go get github.com/clarkduvall/hyperloglog
go get github.com/PuerkitoBio/goquery
go get github.com/syndtr/goleveldb/leveldb
go get github.com/dmuth/golang-stats
go get gopkg.in/yaml.v2
go get github.com/mjibson/esc

esc -o ui/static.go -pkg server ui

go build -o bin/gopa.exe