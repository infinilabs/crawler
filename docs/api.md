
## Gopa API


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
curl -XGET http://127.0.0.1:8001/domains
```


* Get tasks

```
curl -XGET http://127.0.0.1:8001/tasks?from=100&size=10&domain=elasticsearch.cn

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

``` http://localhost:8001/snapshot/?url=http://xxx.com ```


* Get cluster
``` http://127.0.0.1:8001/cluster/info ```

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
