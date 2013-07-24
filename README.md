# Gopa #
[狗爬],A High Performance Distributed  Spider Written in GO.


## Building Gopa ##

linux: Run `make build` to build the Gopa
windows: Run `build` to build the Gopa

## Required By Gopa ##

Gopa using kafka or jafka to store urls,so you need start a kafka server,
a pre-compled jafka can download from here:
https://github.com/medcl/gopa-release/tree/master/kafka

<pre>
wget https://github.com/medcl/gopa-release/raw/master/kafka/jafka-1.3.0-SNAPSHOT-all.tar.gz
tar vxzf jafka-1.3.0-SNAPSHOT-all.tar.gz
cd jafka-1.3.0-SNAPSHOT/bin
./jafka
#run [jafka.exe] if you are running windows
</pre>

## Running Gopa ##

after building the project run `./gopa -h` for a list of commandline options
tips: you can download the pre-compiled [Gopa] from here: https://github.com/medcl/gopa-release

* -seed option : start a crawling, giving a seed url to Gopa. ie: `./gopa -seed=http://www.baidu.com`
* -log option : logging level,can be set to `trace`,`debug`,`info`,`warn`,`error` ,default is `info`

Gopa to support more specified crawling task settings,there is a config file called: `config.ini`
download the sample file,and changes the default value to your own setting,BTW put this config stay with gopa to make it work..

## Stopping Gopa ##

it's safety to press `ctrl+c` stop the current running Gopa,Gopa will handle this,saving the checkpoint,
you may restore the job later,the world is still in your hand.


## Loving By Gopa ##

Gopa is standing on the shoulders of giants,thanks for the following goodies.

* https://github.com/cihub/seelog
* https://github.com/zeebo/sbloom
* http://code.google.com/p/weed-fs
* https://github.com/robfig/config
* https://github.com/PuerkitoBio/purell
* https://github.com/jdamick/kafka.go


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