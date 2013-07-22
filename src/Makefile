build:
    export GOPATH="`pwd`:$GOPATH"
    echo $GOPATH
	go build gopa.go
format:
	gofmt -s -w -tabs=false -tabwidth=4 gopa.go
clean:
	rm -f gopa