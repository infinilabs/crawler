
## Debugging Gopa

* Start gopa with pprof

    `./bin/gopa -pprof`

* HEAP

    `http://localhost:6060/debug/pprof/heap?debug=1`

    `go tool pprof --inuse_space http://localhost:6060/debug/pprof/heap`

    `go tool pprof --alloc_space http://localhost:6060/debug/pprof/heap`

    `go tool pprof --text http://localhost:6060/debug/pprof/heap`

    `go tool pprof --web http://localhost:6060/debug/pprof/heap`

* GC

    `go get -u -v github.com/davecheney/gcvis`

    `gcvis godoc -index -http=:6060`

    `env GODEBUG=gctrace=1 godoc -http=:6060`

* CPU

    `go tool pprof --text http://localhost:6060/debug/pprof/profile`

    `go tool pprof --web  http://localhost:6060/debug/pprof/profile`

    `go tool pprof --web --lines  http://localhost:6060/debug/pprof/profile`


* Reference

    https://software.intel.com/en-us/blogs/2014/05/10/debugging-performance-issues-in-go-programs
