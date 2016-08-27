
## Gopa Storage

    
* Metadata
    *  Domain
    *  Task
    *  Stats
   
   
    key: task   
   
* Payload
   * Web content
      *  Html page
      *  CSS
      *  JS
      *  Img
      *  Others


    key: snapshot:[domain][uuid(url_path)], value: [string,raw html webpage content]
    
    key: metadata:[domain][url_path]:snapshot: value: [UUID, uuid(url_path), key of webpage snapshot]
    key: metadata:[domain][url_path]:size: value: [int, data length]
    key: metadata:[domain][url_path]:last_check: value: [unix timestamp, last check time]
    key: metadata:[domain][url_path]:created: value: [unix timestamp, create datetime]
    key: metadata:[domain][url_path]:updated: value: [unix timestamp, update datetime]
    key: metadata:[domain][url_path]:index_summary: value: [extracted JSON metadata from page, ready for indexing]
    
    key: stats:[domain]: total_size: value: [int, total file size in this domain]
    key: stats:[domain]: file_count: value: [int, total file count in this domain]
