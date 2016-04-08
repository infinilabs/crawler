/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func Load() {

	//load external yaml config
	filename, _ := filepath.Abs("./gopa.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config RuntimeConfig

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	test, _ := yaml.Marshal(config)
	fmt.Printf("Value: %#v\n", config)
	fmt.Printf("Config: %#v\n", string(test))

	//override built-in config

}
