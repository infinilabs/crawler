* [Home](/home) 
* [API](/api)
* [UI](/ui)


*  Working Process

Fetch Url -> ParseLinks -> Persist Links to Disk Queue -> Save Page


* Start GOPA Cluster

- start a leader

` cd node1&& ./gopa -http_bind=:8001 `

- start some followers

` cd node2&& ./gopa -http_bind=:8002  -cluster_seeds=127.0.0.1:8001,127.0.0.1:8002,127.0.0.1:8003 `
` cd node3&& ./gopa -http_bind=:8003  -cluster_seeds=127.0.0.1:8001 -debug `
` ./gopa -cluster_seeds=127.0.0.1:8001,127.0.0.1:8002,127.0.0.1:8003 -debug -data_path=data1 `
