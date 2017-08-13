# Gopa #

![https://github.com/infinitbyte/gopa](https://raw.githubusercontent.com/infinitbyte/gopa/master/static/assets/img/logo.png)

[狗爬], A Spider Written in Go.

[![Travis](https://travis-ci.org/infinitbyte/gopa.svg?branch=master)](https://travis-ci.org/infinitbyte/gopa)
[![Go Report Card](https://goreportcard.com/badge/github.com/infinitbyte/gopa)](https://goreportcard.com/report/github.com/infinitbyte/gopa)
[![Coverage Status](https://coveralls.io/repos/github/infinitbyte/gopa/badge.svg?branch=master)](https://coveralls.io/github/infinitbyte/gopa?branch=master)


## Goals of this project

* Light weight, low footprint, memory requirement should < 100MB
* Easy to deploy, no runtime or dependency required
* Easy to use, no programming or scripts ability needed, out of box features


## Build Gopa ##

Mac/Linux: Run `make build` to build the Gopa

Windows:  Checkout this wiki page - [How to build GOPA on windows](https://github.com/infinitbyte/gopa/wiki/How-to-build-GOPA-on-windows)


## Download ##

[Release](https://github.com/infinitbyte/gopa/releases)


## Start Gopa ##

After download/build the binary file, run `./gopa` to start the Gopa 

Run `./gopa -h` to get the full list of command line options

* -log option : logging level,can be set to `trace`,`debug`,`info`,`warn`,`error` ,default is `info`
* -daemon option : run in background as daemon
* -pprof option : enable and setup pprof/expvar service, eg: localhost:6060 , the endpoint will be: http://localhost:6060/debug/pprof/ and http://localhost:6060/debug/vars
* -cpuprofile option : write cpu profile to this file
* -memprofile option : write memory profile to this file


## Stop Gopa ##

It's safety to press `ctrl+c` stop the current running Gopa, Gopa will handle the rest,saving the checkpoint,
you may restore the job later,the world is still in your hand.

If you are running `Gopa` as daemon, you may stop it like this:

```
 kill -QUIT `pgrep gopa`
```

## UI

Visit `http://127.0.0.1:9001/` for more details.


License
=======
Released under the [Apache License, Version 2.0](https://github.com/infinitbyte/gopa/blob/master/LICENSE) .
