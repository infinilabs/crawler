<img width="200" alt="What a Spider!" src="https://raw.githubusercontent.com/infinitbyte/gopa/master/static/assets/img/logo.svg?sanitize=true">

GOPA, A Spider Written in Go.

[![Travis](https://travis-ci.org/infinitbyte/gopa.svg?branch=master)](https://travis-ci.org/infinitbyte/gopa)
[![Go Report Card](https://goreportcard.com/badge/github.com/infinitbyte/gopa)](https://goreportcard.com/report/github.com/infinitbyte/gopa)
[![Join the chat at https://gitter.im/infinitbyte/gopa](https://badges.gitter.im/infinitbyte/gopa.svg)](https://gitter.im/infinitbyte/gopa?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Finfinitbyte%2Fgopa.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Finfinitbyte%2Fgopa?ref=badge_shield)


## Goal

* Light weight, low footprint, memory requirement should < 100MB
* Easy to deploy, no runtime or dependency required
* Easy to use, no programming or scripts ability needed, out of box features


## Screenshoot

<img width="800" alt="What a Spider! GOPA Spider!" src="https://raw.githubusercontent.com/infinitbyte/gopa/master/docs/assets/img/screenshot/2017.10.20_v0.9.gif">


---


- [How to use](#how-to-use)
  - [Setup](#setup)
    - [Download Pre Built Package](#download-pre-built-package)
    - [Compile The Package Manually](#compile-the-package-manually)
  - [Optional Config](#optional-config)
  - [Start](#start)
  - [Stop](#stop)
- [Configuration](#configuration)
- [UI](#ui)
- [API](#api)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [License](#license)



## How to use

### Setup

First of all, get it, two opinions: download the pre-built package or compile it yourself.

#### Download Pre Built Package

Go to [Release](https://github.com/infinitbyte/gopa/releases) or [Snapshot](https://github.com/infinitbyte/gopa-snapshot/releases) page, download the right package for your platform.

_Note: Darwin is for Mac_

#### Compile The Package Manually

- Mac/Linux: Run `make build` to build the Gopa. <br/>
- Windows:  Checkout this wiki page - [How to build GOPA on windows](https://github.com/infinitbyte/gopa/wiki/How-to-build-GOPA-on-windows).

So far, we have:

> `gopa`, the main program, a single binary.<br/>
> `config/`, elasticsearch related scripts etc.<br/>
> `gopa.yml`, main configuration for gopa.<br/>


### Optional Config

By default, Gopa works well except indexing, if you want to use elasticsearch as indexing, follow these steps:

- Create a index in elasticsearch with script `config/elasticsearch/gopa-index-mapping.sh` (only work for default elasticsearch setting as follows)
<p><details>
  <summary>Example</summary>
  <pre>curl -XPUT "http://localhost:9200/gopa-index" -H 'Content-Type: application/json' -d'
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
       "ext": {
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
       "screenshot_id": {
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
       },
       "last_screenshot_id": {
       "type": "keyword",
       "ignore_above": 256
       }
       }
       }
       }
       }
       }
       }'
</pre>
</details></p>

_Note: Elasticsearch version should >= v5.3_

- Enable index module in `gopa.yml`, update the elasticsearch's setting:
```
  - module: index
    enabled: true
    ui:
      enabled: true
    elasticsearch:
      endpoint: http://localhost:9200
      index_prefix: gopa-
      username: elastic
      password: changeme
```
</details></p>


### Start

Gopa doesn't require any dependencies, simply run `./gopa` to start the program.

Gopa can be run as daemon(_Note: Only available on Linux and Mac_):
<p><details>
  <summary>Example</summary>
  <pre>
➜  gopa git:(master) ✗ ./bin/gopa --daemon
  ________ ________ __________  _____
 /  _____/ \_____  \\______   \/  _  \
/   \  ___  /   |   \|     ___/  /_\  \
\    \_\  \/    |    \    |  /    |    \
 \______  /\_______  /____|  \____|__  /
        \/         \/                \/
[gopa] 0.10.0_SNAPSHOT
///last commit: 99616a2, Fri Oct 20 14:04:54 2017 +0200, medcl, update version to 0.10.0 ///

[10-21 16:01:09] [INF] [instance.go:23] workspace: data/gopa/nodes/0
[gopa] started.</pre>
</details></p>

Also run `./gopa -h` to get the full list of command line options.
<p><details>
  <summary>Example</summary>
  <pre>
➜  gopa git:(master) ✗ ./bin/gopa -h
  ________ ________ __________  _____
 /  _____/ \_____  \\______   \/  _  \
/   \  ___  /   |   \|     ___/  /_\  \
\    \_\  \/    |    \    |  /    |    \
 \______  /\_______  /____|  \____|__  /
        \/         \/                \/
[gopa] 0.10.0_SNAPSHOT
///last commit: 99616a2, Fri Oct 20 14:04:54 2017 +0200, medcl, update version to 0.10.0 ///

Usage of ./bin/gopa:
  -config string
    	the location of config file (default "gopa.yml")
  -cpuprofile string
    	write cpu profile to this file
  -daemon
    	run in background as daemon
  -debug
    	run in debug mode, wi
  -log string
    	the log level,options:trace,debug,info,warn,error (default "info")
  -log_path string
    	the log path (default "log")
  -memprofile string
    	write memory profile to this file
  -pidfile string
    	pidfile path (only for daemon)
  -pprof string
    	enable and setup pprof/expvar service, eg: localhost:6060 , the endpoint will be: http://localhost:6060/debug/pprof/ and http://localhost:6060/debug/vars</pre>
</details></p>


### Stop

It's safety to press `ctrl+c` stop the current running Gopa, Gopa will handle the rest,saving the checkpoint,
you may restore the job later,the world is still in your hand.

If you are running `Gopa` as daemon, you may stop it like this:

```
 kill -QUIT `pgrep gopa`
```

## Configuration

## UI

* Search Console `http://127.0.0.1:9001/`
* Admin Console  `http://127.0.0.1:9001/admin/`

## API

* TBD

## Architecture

<img width="800" alt="What a Spider! GOPA Spider!" src="https://raw.githubusercontent.com/infinitbyte/gopa/master/docs/assets/img/architecture-v1.png">



## Contributing

You are sincerely and warmly welcomed to play with this project,
from UI style to core features,
or just a piece of document,
welcome! let's make it better.


License
=======
Released under the [Apache License, Version 2.0](https://github.com/infinitbyte/gopa/blob/master/LICENSE) .
