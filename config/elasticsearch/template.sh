curl --user elastic:changeme -XDELETE "http://localhost:9200/_template/gopa"

curl --user elastic:changeme -XPUT "http://localhost:9200/_template/gopa" -H 'Content-Type: application/json' -d'
{
"index_patterns": "gopa-*",
"settings": {
    "number_of_shards": 1,
    "index.max_result_window":10000000
  },
  "mappings": {
    "doc": {
      "dynamic_templates": [
        {
          "strings": {
            "match_mapping_type": "string",
            "mapping": {
              "type": "keyword",
              "ignore_above": 256
            }
          }
        }
      ]
    }
  }
}'
