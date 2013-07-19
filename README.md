# Gopa #
[狗爬],A High Performance Distributed  Spider Written in GO.


## Building ##
Run `make build` to build the project

## Running the project ##

after building the project run `./gopa -h` for a list of commandline options

* -seed option : start a crawling. begin with "http://" , ie: `./gopa -seed=http://www.baidu.com`



## Blocks ##

Standing on the shoulders of giants,thanks for these goodies.

* https://github.com/pmylund/go-bloom
* https://github.com/pmylund/go-bitset
* http://code.google.com/p/weed-fs
* https://github.com/Unknwon/goconfig


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