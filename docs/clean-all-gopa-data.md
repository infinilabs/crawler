```
curl --user elastic:changeme -XDELETE "http://localhost:9200/gopa-*"

curl --user elastic:changeme -XPOST "http://localhost:9200/gopa-*/_delete_by_query" -H 'Content-Type: application/json' -d'
{
"query": {"match_all": {}},"size":10000
}'

curl --user elastic:changeme -XDELETE http://localhost:9200/_template/gopa
```