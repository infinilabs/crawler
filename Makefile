SHELL=/bin/bash
CWD=$(shell pwd)
OLDGOPATH=${GOPATH}
NEWGOPATH:=${CWD}:${OLDGOPATH}
export GOPATH=$(NEWGOPATH)


build: clean config update_ui
	go build  -o bin/gopa

update_ui:
	go get github.com/mjibson/esc
	(cd ui&& esc -ignore="static.go|build_static.sh|.DS_Store" -o static.go -pkg server ../ui )

tar: build
	tar cfz bin/gopa.tar.gz bin/gopa

cross-build: clean config update_ui
	go test
	GOOS=windows GOARCH=amd64 go build -o bin/windows64/gopa.exe
	GOOS=darwin  GOARCH=amd64 go build -o bin/darwin64/gopa
	GOOS=linux  GOARCH=amd64 go build -o bin/linux64/gopa

build-linux: clean config update_ui
	go test
	GOOS=linux  GOARCH=amd64 go build -o bin/linux64/gopa

all: clean config update_ui cross-build

all-platform: clean config update_ui cross-build-all-platform

cross-build-all-platform: clean config
	go test
	GOOS=windows GOARCH=amd64     go build -o bin/windows64/gopa.exe
	GOOS=windows GOARCH=386       go build -o bin/windows32/gopa.exe
	GOOS=darwin  GOARCH=amd64     go build -o bin/darwin64/gopa
	GOOS=darwin  GOARCH=386       go build -o bin/darwin32/gopa
	GOOS=linux  GOARCH=amd64      go build -o bin/linux64/gopa
	GOOS=linux  GOARCH=386        go build -o bin/linux32/gopa
	GOOS=linux  GOARCH=arm        go build -o bin/linux_arm/gopa
	GOOS=freebsd  GOARCH=amd64    go build -o bin/freebsd64/gopa
	GOOS=freebsd  GOARCH=386      go build -o bin/freebsd32/gopa
	GOOS=netbsd  GOARCH=amd64     go build -o bin/netbsd64/gopa
	GOOS=netbsd  GOARCH=386       go build -o bin/netbsd32/gopa
	GOOS=openbsd  GOARCH=amd64    go build -o bin/openbsd64/gopa
	GOOS=openbsd  GOARCH=386      go build -o bin/openbsd32/gopa


format:
	gofmt -s -w -tabs=false -tabwidth=4 gopa.go

clean:
	rm -rif bin
	mkdir bin
	mkdir bin/windows64
	mkdir bin/linux64
	mkdir bin/darwin64

config:
	@echo "get Dependencies"
	go env
	go get github.com/cihub/seelog
	go get github.com/robfig/config
	go get github.com/PuerkitoBio/purell
	go get github.com/clarkduvall/hyperloglog
	go get github.com/PuerkitoBio/goquery
	go get github.com/syndtr/goleveldb/leveldb
	go get gopkg.in/yaml.v2
	go get github.com/jmoiron/jsonq
	go get github.com/gorilla/websocket
	go get github.com/boltdb/bolt/...



dist: cross-build package

dist-all: all package

dist-all-platform: all-platform package-all-platform

package:
	@echo "Packaging"
	tar cfz bin/darwin64.tar.gz bin/darwin64
	tar cfz bin/linux64.tar.gz bin/linux64
	tar cfz bin/windows64.tar.gz bin/windows64

package-all-platform:
	@echo "Packaging"
	tar cfz 	 bin/windows64.tar.gz   bin/windows64/gopa.exe
	tar cfz 	 bin/windows32.tar.gz   bin/windows32/gopa.exe
	tar cfz 	 bin/darwin64.tar.gz      bin/darwin64/gopa
	tar cfz 	 bin/darwin32.tar.gz      bin/darwin32/gopa
	tar cfz 	 bin/linux64.tar.gz     bin/linux64/gopa
	tar cfz 	 bin/linux32.tar.gz     bin/linux32/gopa
	tar cfz 	 bin/linux_arm.tar.gz     bin/linux_arm/gopa
	tar cfz 	 bin/freebsd64.tar.gz    bin/freebsd64/gopa
	tar cfz 	 bin/freebsd32.tar.gz    bin/freebsd32/gopa
	tar cfz 	 bin/netbsd64.tar.gz    bin/netbsd64/gopa
	tar cfz 	 bin/netbsd32.tar.gz     bin/netbsd32/gopa
	tar cfz 	 bin/openbsd64.tar.gz     bin/openbsd64/gopa
	tar cfz 	 bin/openbsd32.tar.gz     bin/openbsd32/gopa
