# Gopa #
[狗爬],A High Performance Distributed  Spider Written in GO.

CHANGES

v0.6
breaking change:
1.remove bloom, use leveldb to store urls
feature:
1.crawling speed control
2.cookie supported
3.brief logging format
bugfix:
1.shutdown nil exception
2.wrong relative link in parse phrase


v0.5
feature:
1.ruled fetch
2.fetch/parse offset can be persisted and reloadable
3.http console

v0.4
improve:
1.refactor storage interface,data path are now configable
2.disable pprof by default
3.use local storage instead of kafka,kafka will be removed later
5.check local file's exists first before fetch the remote page
bugfix:
resolve memory leak caused by sbloom filter
feature:
1.download by url template
2.list page download

v0.3
1.adding golang pprof,http://localhost:6060/debug/pprof/
    go tool pprof http://localhost:6060/debug/pprof/heap
    go tool pprof http://localhost:6060/debug/pprof/profile
    go tool pprof http://localhost:6060/debug/pprof/block

2.integrate with kafka to make task controllable and recoverable
3.parameters configable
4.goroutine can be controlled now


v0.2
1.bloom-filter persistence
2.building script works

v0.1
1.just up and run.


