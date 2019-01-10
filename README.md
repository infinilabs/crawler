<img width="200" alt="What a Spider!" src="https://raw.githubusercontent.com/infinitbyte/gopa/master/docs/assets/img/logo.svg?sanitize=true">

GOPA, A Spider Written in Go.

[![Travis](https://travis-ci.org/infinitbyte/gopa.svg?branch=master)](https://travis-ci.org/infinitbyte/gopa)
[![Go Report Card](https://goreportcard.com/badge/github.com/infinitbyte/gopa)](https://goreportcard.com/report/github.com/infinitbyte/gopa)
[![Join the chat at https://gitter.im/infinitbyte/gopa](https://badges.gitter.im/infinitbyte/gopa.svg)](https://gitter.im/infinitbyte/gopa?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)


## Goal

* Light weight, low footprint, memory requirement should < 100MB
* Easy to deploy, no runtime or dependency required
* Easy to use, no programming or scripts ability needed, out of box features


## Screenshoot

<img width="800" alt="What a Spider! GOPA Spider!" src="https://raw.githubusercontent.com/infinitbyte/gopa/master/docs/assets/img/screenshot/2017.10.20_v0.9.gif">


---


- [How to use](#how-to-use)
  - [Requirements](#requirements)
  - [Setup](#setup)
    - [Download Pre Built Package](#download-pre-built-package)
    - [Compile The Package Manually](#compile-the-package-manually)
  - [Required Config](#required-config)
  - [Start](#start)
  - [Stop](#stop)
- [Configuration](#configuration)
- [UI](#ui)
- [API](#api)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [License](#license)



## How to use

### Requirements

* Elasticsearch v5.3+


### Setup

First of all, get it, two opinions: download the pre-built package or compile it yourself.

#### Download Pre Built Package

Go to [Release](https://github.com/infinitbyte/gopa/releases) page, download the right package for your platform.

_Note: Darwin is for Mac_

#### Compile The Package Manually

Requirements
* Golang 1.9+

Supported platform
- Mac/Linux: Run `make build` to build the Gopa. <br/>
- Windows:  Checkout this wiki page - [How to build GOPA on windows](https://github.com/infinitbyte/gopa/wiki/How-to-build-GOPA-on-windows).

So far, we have:

> `gopa`, the main program, a single binary.<br/>
> `gopa.yml`, main configuration for gopa.<br/>


### Required Config

_Note: Elasticsearch version should >= v5.3_

- Enable elastic module in `gopa.yml`, update the elasticsearch's setting:
```
- name: elastic
  enabled: true
  kv_enabled: true
  orm_enabled: true
  elasticsearch:
    endpoint: http://localhost:9200
    index_prefix: gopa-
    username: elastic
    password: changeme
```
</details></p>


### Start

Besides Elasticsearch, Gopa doesn't require any other dependencies, just simply run `./gopa` to start the program.

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
    	run in debug mode, gopa will quit with panic error
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
you may restore the job later, the world is still in your hand.

If you are running `Gopa` as daemon, you may stop it like this:

```
 kill -QUIT `pgrep gopa`
```

## Configuration

## UI

* Search Console `http://127.0.0.1:9000/`
* Admin Console  `http://127.0.0.1:9000/admin/`

## API

## Architecture

<img width="800" alt="What a Spider! GOPA Spider!" src="https://raw.githubusercontent.com/infinitbyte/gopa/master/docs/assets/img/architecture-v1.png">



## Contributing

You are sincerely and warmly welcomed to play with this project, from UI style to core features, or just a piece of document, welcome! let's make it better.


License
=======
Released under the [Apache License, Version 2.0](https://github.com/infinitbyte/gopa/blob/master/LICENSE) .
