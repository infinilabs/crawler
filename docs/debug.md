
## Debugging Gopa

* Start gopa with pprof

    `./bin/gopa -pprof`
    
    
* HEAP

    `http://localhost:6060/debug/pprof/heap?debug=1`

    `go tool pprof --inuse_space http://localhost:6060/debug/pprof/heap`

    `go tool pprof --alloc_space http://localhost:6060/debug/pprof/heap`

    `go tool pprof --text http://localhost:6060/debug/pprof/heap`

    `go tool pprof --web http://localhost:6060/debug/pprof/heap`


* Diff two heap profiles

    `curl -s http://localhost:6060/debug/pprof/heap >1.heap`
    
    `curl -s http://localhost:6060/debug/pprof/heap >2.heap`
    
    `go tool pprof -inuse_objects  --base 1.heap ~/go/src/github.com/infinitbyte/gopa/bin/gopa  2.heap`
    
    use `top` to find top functions, and then use `list func_name` to view the source code.


* GC

    `go get -u -v github.com/davecheney/gcvis`

    `gcvis godoc -index -http=:6060`

    `env GODEBUG=gctrace=1 godoc -http=:6060`

* CPU

    `go tool pprof --text http://localhost:6060/debug/pprof/profile`

    `go tool pprof --web  http://localhost:6060/debug/pprof/profile`

    `go tool pprof --web --lines  http://localhost:6060/debug/pprof/profile`

* Builtin Web

  ```
    ./bin/gopa -pprof=localhost:6060
    go tool pprof -http :9090 http://localhost:6060/debug/pprof/heap
   ```

* Go-torch analysis CPU cycles

    ```
    go get github.com/uber/go-torch
    git clone git@github.com:brendangregg/FlameGraph.git
    export PATH-$PATH:/path/to/FlameGraph
    go-torch --file "torch.svg" --url http://localhost:6060
    ```
    
* Core dump analysis, https://golang.org/doc/gdb

    `gdb path/to/the/binary path/to/the/core`
    `(gdb) where`
    `(gdb) bt full`


* Reference

    https://software.intel.com/en-us/blogs/2014/05/10/debugging-performance-issues-in-go-programs
    http://goog-perftools.sourceforge.net/doc/heap_profiler.html
    

* Chrome

    chrome --headless --disable-gpu --remote-debugging-port=9223
