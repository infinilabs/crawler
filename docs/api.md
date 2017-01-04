
## Gopa API


* Get task status

```
curl -XGET http://localhost:8001/stats
```

* Send seed to Gopa

    ```
    curl -X POST "http://localhost:8001/task/" -d '{
    "seed":"http://elasticsearch.cn"
    }' 
    ```
    
* Update logging config on the fly (visit https://github.com/cihub/seelog/wiki for more details)
    ```
    curl -X POST "http://localhost:8001/setting/seelog/" -d '
    <seelog type="asynctimer" asyncinterval="5000000" minlevel="debug" maxlevel="error">
        ... ...
    </seelog>
    ' 
    ```

* Get web snapshot

``` http://localhost:8001/snapshot/?url=http://xxx.com ```


* Get Cluster
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