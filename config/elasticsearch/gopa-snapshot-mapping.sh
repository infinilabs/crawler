curl -XPUT "http://localhost:9200/gopa-snapshot" -H 'Content-Type: application/json' -d'
{
"mappings": {
"doc": {
"properties": {
"content_type": {
"type": "keyword"
},
"created": {
"type": "date"
},
"file": {
"type": "keyword"
},
"h2": {
"type": "text"
},
"hash": {
"type": "keyword"
},
"lang": {
"type": "keyword"
},
"id": {
"type": "keyword"
},
"images": {
"type": "object"
},
"links": {
"properties": {
"internal": {
"properties": {
"label": {
"type": "keyword"
},
"url": {
"type": "keyword"
},
"screenshot_id": {
"type": "keyword"
}
}
}
}
},
"path": {
"type": "keyword"
},
"size": {
"type": "long"
},
"task_id": {
"type": "keyword"
},
"text": {
"type": "text"
},
"title": {
"type": "text",
"fields": {
"keyword": {
"type": "keyword"
}
}
},
"url": {
"type": "keyword"
},
"version": {
"type": "long"
}
}
}
}
}'
