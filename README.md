# Gopa #

[狗爬], A Spider Written in Go.

[![Travis](https://travis-ci.org/medcl/gopa.svg?branch=master)](https://travis-ci.org/medcl/gopa)
[![Build Status](https://drone.io/github.com/medcl/gopa/status.png)](https://drone.io/github.com/medcl/gopa/latest)


## Building Gopa ##

Mac/Linux: Run `chmod a+x build.sh &./build.sh` to build the Gopa

Windows: Run `build.bat` to build the Gopa


## Download ##

[Release](https://github.com/medcl/gopa/releases)


## Running Gopa ##

After building the project run `./gopa -h` for a list of commandline options

* -seed option : start a crawling, giving a seed url to Gopa. ie: `./gopa -seed=http://www.baidu.com`
* -log option : logging level,can be set to `trace`,`debug`,`info`,`warn`,`error` ,default is `info`
* -daemon option : run in background as daemon
* -pprof option : start pprof service, endpoint: http://localhost:6060/debug/pprof/
* -cpuprofile option : write cpu profile to this file
* -memprofile option : write memory profile to this file


## Stopping Gopa ##

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
    Copyright 2016 Medcl (m^medcl.net)

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
