<img width="200" alt="What a Spider!" src="https://raw.githubusercontent.com/infinitbyte/gopa/master/static/assets/img/logo.svg?sanitize=true">

GOPA, A Spider Written in Go.

[![Travis](https://travis-ci.org/infinitbyte/gopa.svg?branch=master)](https://travis-ci.org/infinitbyte/gopa)
[![Go Report Card](https://goreportcard.com/badge/github.com/infinitbyte/gopa)](https://goreportcard.com/report/github.com/infinitbyte/gopa)
[![Coverage Status](https://coveralls.io/repos/github/infinitbyte/gopa/badge.svg?branch=master)](https://coveralls.io/github/infinitbyte/gopa?branch=master)
[![Join the chat at https://gitter.im/infinitbyte/gopa](https://badges.gitter.im/infinitbyte/gopa.svg)](https://gitter.im/infinitbyte/gopa?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)


## Goal

* Light weight, low footprint, memory requirement should < 100MB
* Easy to deploy, no runtime or dependency required
* Easy to use, no programming or scripts ability needed, out of box features


## Build

Mac/Linux: Run `make build` to build the Gopa

Windows:  Checkout this wiki page - [How to build GOPA on windows](https://github.com/infinitbyte/gopa/wiki/How-to-build-GOPA-on-windows)


## Download

[Release](https://github.com/infinitbyte/gopa/releases)


## Start

After download/build the binary file, run `./gopa` to start the Gopa 

Run `./gopa -h` to get the full list of command line options

```
Usage of ./bin/gopa:
  -config string
        the location of config file, default: gopa.yml (default "gopa.yml")
  -cpuprofile string
        write cpu profile to this file
  -daemon
        run in background as daemon
  -debug
        enable debug
  -log string
        the log level,options:trace,debug,info,warn,error, default: info (default "info")
  -log_path string
        the log path, default: log (default "log")
  -memprofile string
        write memory profile to this file
  -pidfile string
        pidfile path (only for daemon)
  -pprof string
        enable and setup pprof/expvar service, eg: localhost:6060 , the endpoint will be: http://localhost:6060/debug/pprof/ and http://localhost:6060/debug/vars
```


## Stop

It's safety to press `ctrl+c` stop the current running Gopa, Gopa will handle the rest,saving the checkpoint,
you may restore the job later,the world is still in your hand.

If you are running `Gopa` as daemon, you may stop it like this:

```
 kill -QUIT `pgrep gopa`
```

## UI

* Search Console `http://127.0.0.1:9001/`
* Admin Console  `http://127.0.0.1:9001/admin/`

## API

* TBD

## Contribution

You are sincerely and warmly welcomed to play with this project,
from UI style to core features,
or just a piece of document,
welcome! let's make it better.


License
=======
Released under the [Apache License, Version 2.0](https://github.com/infinitbyte/gopa/blob/master/LICENSE) .
