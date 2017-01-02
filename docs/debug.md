./bin/gopa -pprof

** HEAP **

go tool pprof --alloc_space http://localhost:6060/debug/pprof/heap

** GC **

go get -u -v github.com/davecheney/gcvis
gcvis godoc -index -http=:6060

env GODEBUG=gctrace=1 godoc -http=:6060
