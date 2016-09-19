
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
