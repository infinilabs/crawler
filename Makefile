SHELL=/bin/bash

build: clean config
	go build  -o bin/gopa

all: clean config  cross-compile
	go test
	GOOS=windows GOARCH=amd64 go build -o bin/windows64/gopa.exe
	GOOS=darwin  GOARCH=amd64 go build -o bin/darwin64/gopa
	GOOS=linux  GOARCH=amd64 go build -o bin/linux64/gopa

format:
	gofmt -s -w -tabs=false -tabwidth=4 gopa.go

clean:
	rm -rif bin
	mkdir bin

config:
	echo "Prepare Go Path"

	export GOPATH=`pwd`:$GOPATH
	go env

	echo "get Dependencies"
	go get github.com/zeebo/sbloom
	go get github.com/cihub/seelog
	go get github.com/robfig/config
	go get github.com/PuerkitoBio/purell

dist: all
	echo "Packaging"
	tar cfz bin/darwin64.tar.gz bin/darwin64
	tar cfz bin/linux64.tar.gz bin/linux64
	tar cfz bin/windows64.tar.gz bin/windows64

cross-compile:
	echo "Prepare Cross Compiling"
	CWD=`pwd`
	cd $(GOROOT)/src && GOOS=windows GOARCH=amd64 ./make.bash --no-clean 2> /dev/null 1> /dev/null
	cd $(GOROOT)/src && GOOS=darwin  GOARCH=amd64 ./make.bash --no-clean 2> /dev/null 1> /dev/null
	cd $(GOROOT)/src && GOOS=linux  GOARCH=amd64 ./make.bash --no-clean 2> /dev/null 1> /dev/null

	cd $(CWD)
	mkdir bin/windows64
	mkdir bin/linux64
	mkdir bin/darwin64