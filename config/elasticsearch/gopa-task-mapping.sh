curl -XPUT "http://localhost:9200/gopa-task" -H 'Content-Type: application/json' -d'
{
"mappings": {
"doc": {
"properties": {
"breadth": {
"type": "long"
},
"created": {
"type": "date"
},
"depth": {
"type": "long"
},
"host": {
"type": "keyword"
},
"id": {
"type": "keyword"
},
"last_check": {
"type": "date"
},
"last_fetch": {
"type": "date"
},
"next_check": {
"type": "date"
},
"original_url": {
"type": "keyword"
},
"reference_url": {
"type": "keyword"
},
"schema": {
"type": "keyword"
},
"snapshot_created": {
"type": "date"
},
"snapshot_hash": {
"type": "keyword"
},
"snapshot_id": {
"type": "keyword"
},
"last_screenshot_id": {
"type": "keyword"
},
"snapshot_simhash": {
"type": "keyword"
},
"snapshot_version": {
"type": "long"
},
"status": {
"type": "long"
},
"updated": {
"type": "date"
},
"url": {
"type": "keyword"
}
}
}
}
}'
