curl --user elastic:changeme -XPUT "http://localhost:9200/gopa-host" -H 'Content-Type: application/json' -d'
{
"mappings": {
"doc": {
"properties": {
"created": {
"type": "date"
},
"host": {
"type": "keyword"
},
"links_count": {
"type": "long"
},
"updated": {
"type": "date"
}
}
}
}
}'
