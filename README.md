# Gopa #

[狗爬],A High Performance Distributed Spider Written in GO.

[![Build Status](https://drone.io/github.com/medcl/gopa/status.png)](https://drone.io/github.com/medcl/gopa/latest)


## Building Gopa ##

linux: Run `chmod a+x build.sh &./build.sh` to build the Gopa

windows: Run `build.bat` to build the Gopa

a pre-compled gopa can be download from here:
https://bintray.com/medcl/gopa-release/gopa/



## Running Gopa ##

after building the project run `./gopa -h` for a list of commandline options

tips: you can download the pre-compiled [Gopa] from here: https://bintray.com/medcl/gopa-release/gopa

* -seed option : start a crawling, giving a seed url to Gopa. ie: `./gopa -seed=http://www.baidu.com`
* -log option : logging level,can be set to `trace`,`debug`,`info`,`warn`,`error` ,default is `info`

Gopa allow you to specify more detailed crawling-task settings by config a file called: `config.ini`
you can download the sample file,and changes the default value to your own setting,
BTW,do not forget to put this config stay with gopa,and to make it work..

## Stopping Gopa ##

it's safety to press `ctrl+c` stop the current running Gopa,Gopa will handle the rest,saving the checkpoint,
you may restore the job later,the world is still in your hand.


## Loving By Gopa ##

Gopa is standing on the shoulders of giants,thanks for the following goodies.

* https://github.com/zeebo/sbloom
* http://code.google.com/p/weed-fs
* https://github.com/robfig/config
* https://github.com/PuerkitoBio/purell
* https://github.com/cihub/seelog
* https://github.com/clarkduvall/hyperloglog
* https://github.com/PuerkitoBio/goquery
* https://github.com/pantsing/gograce

license
=======
    Copyright 2013 Medcl (m^medcl.net)

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
