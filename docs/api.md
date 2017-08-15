
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
curl -XGET http://127.0.0.1:8001/_cluster/health 
```
```
{
	"cluster_name": "gopa",
	"raft": {
		"leader": "127.0.0.1:13001",
		"seeds": [
			"127.0.0.1:13002",
			"127.0.0.1:13003"
		],
		"stats": {
			"applied_index": "1",
			"commit_index": "1",
			"fsm_pending": "0",
			"last_contact": "never",
			"last_log_index": "1",
			"last_log_term": "54",
			"last_snapshot_index": "0",
			"last_snapshot_term": "0",
			"num_peers": "2",
			"state": "Leader",
			"term": "54"
		}
	},
	"status": "green"
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
