curl -XPUT "http://localhost:9200/gopa-index" -H 'Content-Type: application/json' -d'
{
"mappings": {
"doc": {
"properties": {
"host": {
"type": "keyword",
"ignore_above": 256
},
"snapshot": {
"properties": {
"bold": {
"type": "text"
},
"url": {
"type": "keyword",
"ignore_above": 256
},
"content_type": {
"type": "keyword",
"ignore_above": 256
},
"file": {
"type": "keyword",
"ignore_above": 256
},
"h1": {
"type": "text"
},
"h2": {
"type": "text"
},
"h3": {
"type": "text"
},
"h4": {
"type": "text"
},
"hash": {
"type": "keyword",
"ignore_above": 256
},
"id": {
"type": "keyword",
"ignore_above": 256
},
"images": {
"properties": {
"external": {
"properties": {
"label": {
"type": "text"
},
"url": {
"type": "keyword",
"ignore_above": 256
}
}
},
"internal": {
"properties": {
"label": {
"type": "text"
},
"url": {
"type": "keyword",
"ignore_above": 256
}
}
}
}
},
"italic": {
"type": "text"
},
"links": {
"properties": {
"external": {
"properties": {
"label": {
"type": "text"
},
"url": {
"type": "keyword",
"ignore_above": 256
}
}
},
"internal": {
"properties": {
"label": {
"type": "text"
},
"url": {
"type": "keyword",
"ignore_above": 256
}
}
}
}
},
"path": {
"type": "keyword",
"ignore_above": 256
},
"sim_hash": {
"type": "keyword",
"ignore_above": 256
},
"lang": {
"type": "keyword",
"ignore_above": 256
},
"size": {
"type": "long"
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
"version": {
"type": "long"
}
}
},
"task": {
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
"id": {
"type": "keyword",
"ignore_above": 256
},
"original_url": {
"type": "keyword",
"ignore_above": 256
},
"reference_url": {
"type": "keyword",
"ignore_above": 256
},
"schema": {
"type": "keyword",
"ignore_above": 256
},
"status": {
"type": "integer"
},
"updated": {
"type": "date"
},
"url": {
"type": "keyword",
"ignore_above": 256
}
}
}
}
}
}
}'
