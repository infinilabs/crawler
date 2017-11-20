
* Create pipeline
```
curl -XPOST http://127.0.0.1:8001/pipeline/config/ -d'
{
 "name": "test_pipe_line",
 "start": {
  "joint": "init_task","enabled": true

 },
 "process": [
   {
   "joint": "url_normalization","enabled": true
  },
 {
   "joint": "fetch","enabled": true
  },{
   "joint": "save_snapshot_fs","enabled": true
  }
 ],
 "end": {
  "joint": "save_task","enabled": true
 }
}'
```


* Get Pipeline tasks
```
curl -XGET http://127.0.0.1:8001/pipeline/tasks/
{
"tasks":[
{
    "crawler":{ }
}
]
}
```

* Assign a pipeline to host
```
curl -XPOST http://127.0.0.1:8001/host/pipeline_config/
{
 "host":"baidu.com",
 "url_pattern":".*",
 "sort_order":1,
 "runner":"fetch"
 "pipeline_id":"123123123"
}
```


