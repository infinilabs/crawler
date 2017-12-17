SHELL=/bin/bash

# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif

CURDIR := $(shell pwd)
OLDGOPATH:= $(GOPATH)
NEWGOPATH:= $(CURDIR):$(CURDIR)/vendor:$(GOPATH)

GO        := GO15VENDOREXPERIMENT="1" go
GOBUILD  := GOPATH=$(NEWGOPATH) CGO_ENABLED=1  $(GO) build -ldflags -s
GOTEST   := GOPATH=$(NEWGOPATH) CGO_ENABLED=1  $(GO) test -ldflags -s

ARCH      := "`uname -s`"
LINUX     := "Linux"
MAC       := "Darwin"
GO_FILES=$(find . -iname '*.go' | grep -v /vendor/)
PKGS=$(go list ./... | grep -v /vendor/)

.PHONY: all build update test clean

default: build

build: config
	@echo $(GOPATH)
	@echo $(NEWGOPATH)
	$(GOBUILD) -o bin/gopa

build-cluster-test: build
	cd bin && mkdir node1 node2 node3 && cp gopa node1 && cp gopa node2 && cp gopa node3

# used to build the binary for gdb debugging
build-race: clean config update-ui
	$(GOBUILD) -gcflags "-m -N -l" -race -o bin/gopa

update-ui:
	$(GO) get github.com/infinitbyte/esc
	(cd static&& esc -ignore="static.go|build_static.sh|.DS_Store" -o static.go -pkg static ../static )

update-template-ui:
	$(GO) get github.com/infinitbyte/ego/cmd/ego
	cd modules/ui/ && ego
	cd modules/index/ui/ && ego
	cd plugins/ && ego

tar: build
	cd bin && tar cfz ../bin/gopa.tar.gz gopa config gopa.yml

cross-build: clean config update-ui
	$(GO) test
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o bin/gopa-windows64.exe
	GOOS=darwin  GOARCH=amd64 $(GOBUILD) -o bin/gopa-darwin64
	GOOS=linux  GOARCH=amd64 $(GOBUILD) -o bin/gopa-linux64

build-win:
	CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64     $(GOBUILD) -o bin/gopa-windows64.exe
	CC=i686-w64-mingw32-gcc   CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386         $(GOBUILD) -o bin/gopa-windows32.exe

build-linux:
	GOOS=linux  GOARCH=amd64  $(GOBUILD) -o bin/gopa-linux64
	GOOS=linux  GOARCH=386    $(GOBUILD) -o bin/gopa-linux32

build-darwin:
	GOOS=darwin  GOARCH=amd64     $(GOBUILD) -o bin/gopa-darwin64
	GOOS=darwin  GOARCH=386       $(GOBUILD) -o bin/gopa-darwin32

all: clean config update-ui cross-build

all-platform: clean config update-ui cross-build-all-platform

cross-build-all-platform: clean config build-bsd build-linux build-darwin build-win

build-bsd:
	GOOS=freebsd  GOARCH=amd64    $(GOBUILD) -o bin/gopa-freebsd64
	GOOS=freebsd  GOARCH=386      $(GOBUILD) -o bin/gopa-freebsd32
	GOOS=netbsd  GOARCH=amd64     $(GOBUILD) -o bin/gopa-netbsd64
	GOOS=netbsd  GOARCH=386       $(GOBUILD) -o bin/gopa-netbsd32
	GOOS=openbsd  GOARCH=amd64    $(GOBUILD) -o bin/gopa-openbsd64
	GOOS=openbsd  GOARCH=386      $(GOBUILD) -o bin/gopa-openbsd32

format:
	gofmt -l -s -w .

clean_data:
	rm -rif dist
	rm -rif data
	rm -rif log

clean: clean_data
	rm -rif bin
	mkdir bin


update-commit-log:
	echo -e "package env\nconst lastCommitLog  =\""`git log -1 --pretty=format:"%h, %ad, %an, %s"` "\"\nconst buildDate  =\"`date`\"" > core/env/commit_log.go

config: update-commit-log update-ui update-template-ui
	@echo "init config"
	$(GO) env
	cp stop.sh bin/stop.sh
	cp gopa.yml bin/gopa.yml
	cp -r config bin


fetch-depends:
	@echo "get Dependencies"
	$(GO) get github.com/cihub/seelog
	$(GO) get github.com/PuerkitoBio/purell
	$(GO) get github.com/clarkduvall/hyperloglog
	$(GO) get github.com/PuerkitoBio/goquery
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
	$(GO) get github.com/jasonlvhit/gocron
	$(GO) get github.com/quipo/statsd
	$(GO) get github.com/go-sql-driver/mysql
	$(GO) get github.com/jbowles/cld2_nlpt


dist: cross-build package

dist-major-platform: all package

dist-all-platform: all-platform package-all-platform

package:
	@echo "Packaging"
	cd bin && tar cfz ../bin/darwin64.tar.gz darwin64 config gopa.yml stop.sh
	cd bin && tar cfz ../bin/linux64.tar.gz linux64 config gopa.yml stop.sh
	cd bin && tar cfz ../bin/windows64.tar.gz windows64 config gopa.yml stop.sh

package-all-platform: package-darwin-platform package-linux-platform package-windows-platform
	@echo "Packaging all"
	cd bin && tar cfz ../bin/freebsd64.tar.gz     gopa-freebsd64 config gopa.yml stop.sh
	cd bin && tar cfz ../bin/freebsd32.tar.gz     gopa-freebsd32 config gopa.yml stop.sh
	cd bin && tar cfz ../bin/netbsd64.tar.gz      gopa-netbsd64 config gopa.yml stop.sh
	cd bin && tar cfz ../bin/netbsd32.tar.gz      gopa-netbsd32 config gopa.yml stop.sh
	cd bin && tar cfz ../bin/openbsd64.tar.gz     gopa-openbsd64 config gopa.yml stop.sh
	cd bin && tar cfz ../bin/openbsd32.tar.gz     gopa-openbsd32 config gopa.yml stop.sh


package-darwin-platform:
	@echo "Packaging Darwin"
	cd bin && tar cfz ../bin/darwin64.tar.gz      gopa-darwin64 config gopa.yml stop.sh
	cd bin && tar cfz ../bin/darwin32.tar.gz      gopa-darwin32 config gopa.yml stop.sh

package-linux-platform:
	@echo "Packaging Linux"
	cd bin && tar cfz ../bin/linux64.tar.gz     gopa-linux64 config gopa.yml stop.sh
	cd bin && tar cfz ../bin/linux32.tar.gz     gopa-linux32 config gopa.yml stop.sh

package-windows-platform:
	@echo "Packaging Windows"
	cd bin && tar cfz ../bin/windows64.tar.gz   gopa-windows64.exe config gopa.yml stop.sh
	cd bin && tar cfz ../bin/windows32.tar.gz   gopa-windows32.exe config gopa.yml stop.sh

test:
	go get -u github.com/kardianos/govendor
	go get github.com/stretchr/testify/assert
	govendor test +local
	#$(GO) test -timeout 60s ./... --ignore ./vendor
	#GORACE="halt_on_error=1" go test ./... -race -timeout 120s  --ignore ./vendor
	#go test -bench=. -benchmem

check:
	$(GO)  get github.com/golang/lint/golint
	$(GO)  get honnef.co/go/tools/cmd/megacheck
	test -z $(gofmt -s -l $GO_FILES)    # Fail if a .go file hasn't been formatted with gofmt
	$(GO) test -v -race $(PKGS)            # Run all the tests with the race detector enabled
	$(GO) vet $(PKGS)                      # go vet is the official Go static analyzer
	@echo "go tool vet"
	go tool vet main.go
	go tool vet core
	go tool vet modules
	megacheck $(PKGS)                      # "go vet on steroids" + linter
	golint -set_exit_status $(PKGS)    # one last linter

errcheck:
	go get github.com/kisielk/errcheck
	errcheck -blank $(PKGS)

cover:
	go get github.com/mattn/goveralls
	go test -v -cover -race -coverprofile=data/coverage.out
	goveralls -coverprofile=data/coverage.out -service=travis-ci -repotoken=$COVERALLS_TOKEN

cyclo:
	go get -u github.com/fzipp/gocyclo
	gocyclo -top 10 -over 12 $$(ls -d */ | grep -v vendor)

benchmarks:
	go test github.com/infinitbyte/gopa/core/util -benchtime=1s -bench ^Benchmark -run ^$
	go test github.com/infinitbyte/gopa//modules/crawler/pipe -benchtime=1s -bench  ^Benchmark -run ^$
