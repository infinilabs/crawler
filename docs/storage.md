
## Gopa Storage

    
* Task
    *  Task metadata
   
   
    key: [domain][url_path]:snapshot: value: [UUID, uuid(url_path), key of webpage snapshot]
    key: [domain][url_path]:size: value: [int, data length]
    key: [domain][url_path]:last_check: value: [unix timestamp, last check time]
    key: [domain][url_path]:created: value: [unix timestamp, create datetime]
    key: [domain][url_path]:updated: value: [unix timestamp, update datetime]
    key: [domain][url_path]:index_summary: value: [extracted JSON metadata from page, ready for indexing]
    
* Stats
   
   
    key: [domain]: total_size: value: [int, total file size in this domain]
    key: [domain]: file_count: value: [int, total file count in this domain]
   
   
* Snapshot
   * Web content
      *  Html page
      *  CSS
      *  JS
      *  Img
      *  Others


    key: [domain][uuid(url_path)], value: [string,raw html webpage content]
    
