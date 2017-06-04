SHELL=/bin/bash

# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif

CURDIR := $(shell pwd)
OLDGOPATH:= $(GOPATH)
NEWGOPATH:= $(GOPATH):$(CURDIR)/vendor

GO        := GO15VENDOREXPERIMENT="1" go
GOBUILD  := GOPATH=$(NEWGOPATH) $(GO) build -ldflags -s
GOTEST   := GOPATH=$(NEWGOPATH) $(GO) test -ldflags -s

ARCH      := "`uname -s`"
LINUX     := "Linux"
MAC       := "Darwin"
PACKAGES  := $$(go list ./...| grep -vE 'vendor')
FILES     := $$(find . -name '*.go' | grep -vE 'vendor')


.PHONY: all build update test clean

default: clean build

build: config update-ui update-template-ui
	@echo $(GOPATH)
	@echo $(NEWGOPATH)
	$(GOBUILD) -o bin/gopa


build-cluster-test: build
	cd bin && mkdir node1 node2 node3 && cp gopa node1 && cp gopa node2 && cp gopa node3

build-grace: clean config update-ui
	$(GOBUILD) -gcflags "-N -l" -race -o bin/gopa

update-ui:
	$(GO) get github.com/infinitbyte/esc
	(cd static&& esc -ignore="static.go|build_static.sh|.DS_Store" -o static.go -pkg static ../static )

update-template-ui:
	$(GO) get github.com/infinitbyte/ego/cmd/ego
	cd modules/ui/ && ego

tar: build
	cd bin && tar cfz ../bin/gopa.tar.gz gopa

cross-build: clean config update-ui
	$(GO) test
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o bin/windows64/gopa.exe
	GOOS=darwin  GOARCH=amd64 $(GOBUILD) -o bin/darwin64/gopa
	GOOS=linux  GOARCH=amd64 $(GOBUILD) -o bin/linux64/gopa

build-linux: clean config update-ui
	$(GO) test
	GOOS=linux  GOARCH=amd64 $(GOBUILD) -o bin/linux64/gopa

all: clean config update-ui cross-build

all-platform: clean config update-ui cross-build-all-platform

cross-build-all-platform: clean config build-bsd
	GOOS=windows GOARCH=amd64     $(GOBUILD) -o bin/windows64/gopa.exe
	GOOS=windows GOARCH=386       $(GOBUILD) -o bin/windows32/gopa.exe
	GOOS=darwin  GOARCH=amd64     $(GOBUILD) -o bin/darwin64/gopa
	GOOS=darwin  GOARCH=386       $(GOBUILD) -o bin/darwin32/gopa
	GOOS=linux  GOARCH=amd64      $(GOBUILD) -o bin/linux64/gopa
	GOOS=linux  GOARCH=386        $(GOBUILD) -o bin/linux32/gopa
	GOOS=linux  GOARCH=arm        $(GOBUILD) -o bin/linux_arm/gopa

build-bsd:
	GOOS=freebsd  GOARCH=amd64    $(GOBUILD) -o bin/freebsd64/gopa
	GOOS=freebsd  GOARCH=386      $(GOBUILD) -o bin/freebsd32/gopa
	GOOS=netbsd  GOARCH=amd64     $(GOBUILD) -o bin/netbsd64/gopa
	GOOS=netbsd  GOARCH=386       $(GOBUILD) -o bin/netbsd32/gopa
	GOOS=openbsd  GOARCH=amd64    $(GOBUILD) -o bin/openbsd64/gopa
	GOOS=openbsd  GOARCH=386      $(GOBUILD) -o bin/openbsd32/gopa

format:
	gofmt -l -s -w .

clean_data:
	rm -rif dist
	rm -rif data
	rm -rif log

clean: clean_data
	rm -rif bin
	mkdir bin
	mkdir bin/windows64
	mkdir bin/linux64
	mkdir bin/darwin64


update-commit-log:
	echo -e "package env\nconst LastCommitLog  =\""`git log -1 --pretty=format:"%h, %ad, %an, %s"` "\"\nconst BuildDate  =\"`date`\"" > core/env/commit_log.go

config: update-commit-log
	@echo "init config"
	$(GO) env

fetch-depends:
	@echo "get Dependencies"
	$(GO) get github.com/cihub/seelog
	$(GO) get github.com/robfig/config
	$(GO) get github.com/PuerkitoBio/purell
	$(GO) get github.com/clarkduvall/hyperloglog
	$(GO) get github.com/PuerkitoBio/goquery
	$(GO) get github.com/syndtr/goleveldb/leveldb
	$(GO) get github.com/jmoiron/jsonq
	$(GO) get github.com/gorilla/websocket
	$(GO) get github.com/boltdb/bolt/...
	$(GO) get github.com/alash3al/goemitter
	$(GO) get github.com/bkaradzic/go-lz4
	$(GO) get github.com/elgs/gojq
	$(GO) get github.com/kardianos/osext
	$(GO) get github.com/zeebo/sbloom
	$(GO) get github.com/asdine/storm
	$(GO) get github.com/julienschmidt/httprouter
	$(GO) get github.com/rs/xid
	$(GO) get github.com/seiflotfy/cuckoofilter
	$(GO) get github.com/hashicorp/raft
	$(GO) get github.com/hashicorp/raft-boltdb
	$(GO) get github.com/jaytaylor/html2text
	$(GO) get github.com/asdine/storm/codec/protobuf
	$(GO) get github.com/ryanuber/go-glob
	$(GO) get github.com/gorilla/sessions
	$(GO) get github.com/mattn/go-sqlite3
	$(GO) get github.com/jinzhu/gorm
	$(GO) get github.com/stretchr/testify/assert
	$(GO) get github.com/spf13/viper
	$(GO) get -t github.com/RoaringBitmap/roaring
	$(GO) get github.com/elastic/go-ucfg
	$(GO) get github.com/blevesearch/bleve

dist: cross-build package

dist-major-platform: all package

dist-all-platform: all-platform package-all-platform

package:
	@echo "Packaging"
	cd bin && tar cfz ../bin/darwin64.tar.gz darwin64
	cd bin && tar cfz ../bin/linux64.tar.gz linux64
	cd bin && tar cfz ../bin/windows64.tar.gz windows64

package-all-platform:
	@echo "Packaging"
	cd bin && tar cfz ../bin/windows64.tar.gz   windows64/gopa.exe
	cd bin && tar cfz ../bin/windows32.tar.gz   windows32/gopa.exe
	cd bin && tar cfz ../bin/darwin64.tar.gz      darwin64/gopa
	cd bin && tar cfz ../bin/darwin32.tar.gz      darwin32/gopa
	cd bin && tar cfz ../bin/linux64.tar.gz     linux64/gopa
	cd bin && tar cfz ../bin/linux32.tar.gz     linux32/gopa
	cd bin && tar cfz ../bin/linux_arm.tar.gz     linux_arm/gopa
	cd bin && tar cfz ../bin/freebsd64.tar.gz    freebsd64/gopa
	cd bin && tar cfz ../bin/freebsd32.tar.gz    freebsd32/gopa
	cd bin && tar cfz ../bin/netbsd64.tar.gz    netbsd64/gopa
	cd bin && tar cfz ../bin/netbsd32.tar.gz     netbsd32/gopa
	cd bin && tar cfz ../bin/openbsd64.tar.gz     openbsd64/gopa
	cd bin && tar cfz ../bin/openbsd32.tar.gz     openbsd32/gopa

test:
	go get -u github.com/kardianos/govendor
	govendor test +local
	#$(GO) test -timeout 60s ./... --ignore ./vendor
	#GORACE="halt_on_error=1" go test ./... -race -timeout 120s  --ignore ./vendor

check:
	$(GO) get github.com/golang/lint/golint

	@echo "vet"
	@ go tool vet $(FILES) 2>&1 | awk '{print} END{if(NR>0) {exit 1}}'
	@echo "vet --shadow"
	@ go tool vet --shadow $(FILES) 2>&1 | awk '{print} END{if(NR>0) {exit 1}}'
	@echo "golint"
	@ golint ./... 2>&1 | grep -vE 'context\.Context|LastInsertId|NewLexer|\.pb\.go' | awk '{print} END{if(NR>0) {exit 1}}'
	@echo "gofmt (simplify)"
	@ gofmt -s -l -w $(FILES) 2>&1 | awk '{print} END{if(NR>0) {exit 1}}'

errcheck:
	go get github.com/kisielk/errcheck
	errcheck -blank $(PACKAGES)
