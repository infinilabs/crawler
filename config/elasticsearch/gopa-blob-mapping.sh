curl --user elastic:changeme -XPUT "http://localhost:9200/gopa-blob/" -H 'Content-Type: application/json' -d'
{
"mappings": {
"doc": {
"properties": {
"content": {
"type": "binary",
"doc_values":false
}
}
}
}
}'
