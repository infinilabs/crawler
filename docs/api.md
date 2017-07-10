
## Gopa API


* Get basic info

```
curl -XGET http://127.0.0.1:8001/
```

* Get task status

```
curl -XGET http://localhost:8001/stats
```

* Send seed to Gopa

```
curl -XPOST "http://localhost:8001/task/" -d '{
"seed":"http://elasticsearch.cn"
}' 
```

* Get domains

```
curl -XGET http://127.0.0.1:8001/domain
```


* Get tasks

```
curl -XGET http://127.0.0.1:8001/task?from=100&size=10&domain=elasticsearch.cn

```


* Get logging config

```
curl -XGET http://127.0.0.1:8001/setting/logger
```

```
{
 "realtime": false,
 "log_level": "info",
 "push_log_level": "info",
 "func_pattern": "*",
 "file_pattern": "*"
}
```

    
* Update logging config on the fly

```
curl -XPOST "http://localhost:8001/setting/logger/" -d '
{
"realtime": true,
"log_level": "info",
"push_log_level": "info",
"func_pattern": "*",
"file_pattern": "*"
}' 
```

* Get web snapshot

``` 
curl -XGET http://localhost:8001/snapshot/?url=http://xxx.com 
```


* Get cluster
``` 
curl -XGET http://127.0.0.1:8001/cluster/status 
```
```
{
	"addr": "Node at :13003 [Follower]",
	"leader": ":13000",
	"stats": {
		"applied_index": "21",
		"commit_index": "21",
		"fsm_pending": "0",
		"last_contact": "55.516082ms",
		"last_log_index": "21",
		"last_log_term": "408",
		"last_snapshot_index": "0",
		"last_snapshot_term": "0",
		"num_peers": "4",
		"state": "Follower",
		"term": "408"
	}
}
```

* Get joints
```
curl -XGET http://127.0.0.1:8001/joint/
```
```
{
 "fetch": {},
 "html2text": {
  "MergeWhitespace": false
 },
 "ignore_timeout": {
  "IgnoreTimeoutAfterCount": 0
 },
 "init_task": {
  "data": null,
  "Task": null
 },
 "load_metadata": {},
 "parse": {
  "DispatchLinks": false,
  "MaxDepth": 0
 },
 "save2db": {
  "CompressBody": false
 },
 "save2fs": {},
 "save_task": {
  "IsCreate": false
 },
 "url_checked_filter": {
  "data": null,
  "SkipPageParsePattern": null
 },
 "url_ext_filter": {
  "SkipPageParsePattern": null
 },
 "url_normalization": {
  "FollowSubDomain": false
 }
}
```

* Create pipeline
```
curl -XPOST http://127.0.0.1:8001/joint/ -d'
{
 "name": "test_pipe_line",
 "context": {
  "data": {
   "URL": "http://facebook.com",
   "HOST": "facebook.com",
   "DEPTH": 0,
   "BREADTH": 0
  },
  "phrase": 0
 },
 "start": {
  "joint": "empty",
  "parameters": {
   "key": "value"
  }
 },
 "process": [
 {
   "joint": "fetch",
   "parameters": {
   "proxy": "socks5://127.0.0.1:9742"
   }
  }
 ],
 "end": null
}'
```
