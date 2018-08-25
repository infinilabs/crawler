# Gopa #
[狗爬], [aims to be] A high performance distributed and lightweight spider written in GO .

## CHANGES


#### v0.12

##### breaking changes

##### features

##### improvement

##### bugfix


#### v0.11

##### breaking changes
1. extract common codebase to another repo: https://github.com/infinitbyte/framework
2. sqlite retired, elasticsearch is the first citizen

##### features
1. add a new cmd `static_fs` to support load static files from folder
2. auto generate elasticsearch mapping and template, no need to manual create mapping first

##### improvement
1. optimize sql, speed up task list
2. enable cross domain requests

##### bugfix
1. fix mysql as database option
2. update update_check_time, fix init next_fetch time


#### v0.10

##### breaking changes
1. refactor domain to host, api and mapping has changed
2. refactor module, update yml settings: module->name

##### features
1. dynamic create pipelines
2. init plugin architecture
3. support extract tags by css path
4. add chrome fetch joint, via chrome debug protocol
5. add auto-completion to search ui
6. search ui support mobile
7. support access control by github oauth

##### improvement
1. remove goleveldb due to memory leak
2. update logo
3. remove hard coded version
4. update task UI, support filter by status and host
5. clean offset_canvas menu
##### bugfix


#### v0.9

##### breaking changes
1. move repo to infinitbyte/gopa, for better collaboration, namespace changed as well  
2. separate API and UI, listen on different port
3. add mysql as database option
3. add elasticsearch as database option
4. add elasticsearch as blob(snapshot) datastore

##### features
1. task fetch and update with stepped delay
2. add hash joint to crawler pipeline
3. dispatch tasks and auto update tasks
4. add proxy to fetch joint
5. filter url before push to checker
6. add rules config to url filter 
7. support elasticsearch as database store
8. add task_deduplication in the check phrase
9. add content hash check to detect duplication
10. refactor webhunter, support basic auth
11. add pipeline joint to detect the language of webpage
12. add search ui

##### improvement
1. multi instance support on local machine
2. streamline clustering on local machine
3. modules and pipelines dynamic config ready
4. pipeline and context refactored to support dynamic parameters
5. save snapshot to KV store and update task management
6. optimize shutdown logic, reduce half of goroutines
7. add a wiki about how to build gopa on windows
8. remove timeout in queue by default
9. improve statsd performance with buffered client
10. refine log level, enable pprof to config listen address
11. update task ui, limit length of name
12. detect dead process, re-place lock file
13. persist auto-incremented id sequence to disk
14. simplified joint register
15. add high performance tolowercase and touppercase func
16. add queue stats api

##### bugfix
1. remove simhash due to poor performance and memory leak
2. fix wrong relative url by using unicode index
3. fix statsd no data was send out
4. fix poor string merge performance
5. fix http goroutine leak


#### v0.8

##### features
1. raft clustering
2. dynamic change logging setting from the console, can be filter log by level, message, file and function name
3. dynamic create pipeline
4. add tls to security api and websocket
5. add proxy to crawler pipeline

##### improvement
1. use template engine, UI refactoring
2. add a logo

##### bugfix
1. fix incorrect stats number, incorrect task filter
2. fix incorrect redirect handler, url ignored

#### v0.7
##### features:
1. add stats api to expose the task info, http://localhost:8001/stats
2. add websocket and simple ui to interact with Gopa, http://localhost:8001/ui/
3. add task api to accept seed
4. dynamic change the seelog config via api, [GET/POST] http://localhost:8001/setting/seelog/
5. follow 301/302 redirect, and continue fetch
6. add boltdb status page, http://localhost:8001/ui/boltdb
7. add pipeline framework to create crawler
8. add command to dynamic change logging level and add seed url
8. export metrics to statsD
9. support daemon mode in linux and darwin
10. add task management api

##### improvement:
1. add update_ui setup to Makefile in order to build static ui
2. add git commit log and build_date to gopa binary
3. console ui support websocket reconnect

#### v0.6
##### breaking change:
1. remove bloom, use leveldb to store urls

##### feature:
1. crawling speed control
2. cookie supported
3. brief logging format

##### bugfix:
1. shutdown nil exception
2. wrong relative link in parse phrase


#### v0.5
##### feature:
1. ruled fetch
2. fetch/parse offset can be persisted and reloadable
3. http console

#### v0.4
##### improvement:
1. refactor storage interface,data path are now configable
2. disable pprof by default
3. use local storage instead of kafka,kafka will be removed later
5. check local file's exists first before fetch the remote page

##### bugfix:
1. resolve memory leak caused by sbloom filter

##### feature:
1. download by url template
2. list page download

#### v0.3
1. adding golang pprof, http://localhost:6060/debug/pprof/
  - go tool pprof http://localhost:6060/debug/pprof/heap
  - go tool pprof http://localhost:6060/debug/pprof/profile
  - go tool pprof http://localhost:6060/debug/pprof/block

2. integrate with kafka to make task controllable and recoverable
3. parameters configable
4. goroutine can be controlled now


#### v0.2
1. bloom-filter persistence
2. building script works

#### v0.1
1. just up and run.


