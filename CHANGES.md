# Gopa #
[狗爬],A High Performance Distributed  Spider Written in GO.

CHANGES

v0.4
1.align storage path,store paged webpage
2.disable pprof by default
3.fix some critical bug and performance optimized
5.save paged webpage

v0.3
1.adding golang pprof,http://localhost:6060/debug/pprof/
    go tool pprof http://localhost:6060/debug/pprof/heap
    go tool pprof http://localhost:6060/debug/pprof/profile
    go tool pprof http://localhost:6060/debug/pprof/block

2.integrate with kafka to make task controllable and recoverable
3.parameters configable
4.goroutine canbe controlled now


v0.2
1.bloom-filter persistence
2.building script works

v0.1
1.just up and run.


