# Gopa #

![https://github.com/medcl/gopa](https://raw.githubusercontent.com/medcl/gopa/master/static/assets/img/logo.png)

[狗爬], A Spider Written in Go.

[![Travis](https://travis-ci.org/medcl/gopa.svg?branch=master)](https://travis-ci.org/medcl/gopa)
[![Go Report Card](https://goreportcard.com/badge/github.com/medcl/gopa)](https://goreportcard.com/report/github.com/medcl/gopa)


## Build Gopa ##

Mac/Linux: Run `make build` to build the Gopa


## Download ##

[Release](https://github.com/medcl/gopa/releases)


## Start Gopa ##

After download/build the binary file, run `./gopa` to start the Gopa 

Run `./gopa -h` to get the full list of commandline options

* -log option : logging level,can be set to `trace`,`debug`,`info`,`warn`,`error` ,default is `info`
* -daemon option : run in background as daemon
* -pprof option : start pprof service, endpoint: http://localhost:6060/debug/pprof/
* -cpuprofile option : write cpu profile to this file
* -memprofile option : write memory profile to this file


## Stop Gopa ##

It's safety to press `ctrl+c` stop the current running Gopa, Gopa will handle the rest,saving the checkpoint,
you may restore the job later,the world is still in your hand.

If you are running `Gopa` as daemon, you can stop it like this:

```
 kill -QUIT `pgrep gopa`
```

## UI

Visit `http://127.0.0.1:8001/ui/` for more details.


License
=======
Released under the [Apache License, Version 2.0](https://github.com/medcl/gopa/blob/master/LICENSE) .
