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

package main

import (
	"flag"
	"fmt"
	"github.com/infinitbyte/framework/core/util"
	"github.com/infinitbyte/gopa/model"
	"os"
	"path"
	"strings"
	"time"
)

var backupOutput string
var scope string
var host string

func main() {
	//Backup and restore tasks
	flag.StringVar(&host, "host", "localhost:8001", "the host,eg: localhost:8001")
	flag.StringVar(&backupOutput, "out", "data/backup/", "the output path,eg: data/backup/")
	flag.StringVar(&scope, "scope", "tasks,snapshots,hosts,projects", "the scope to do the snapshot,eg:tasks")
	flag.Parse()

	objs := strings.Split(scope, ",")
	for _, x := range objs {
		execute(x)
	}

}

func execute(x string) {
	url := "http://%s/%s/?from=%v&size=%v"
	from := 0
	size := 100
	os.MkdirAll(backupOutput, 0777)

	output := path.Join(backupOutput, "/", x+"_"+util.FormatTimeForFileName(time.Now())+".json")
begin:

	result, err := util.HttpGet(fmt.Sprintf(url, host, x, from, size))
	if err != nil {
		panic(err)
	}
	v := struct {
		Total  int64
		Result []model.Task
	}{}

	util.FromJSONBytes(result.Body, &v)
	resultSize := len(v.Result)

	fmt.Printf("%s,%v,%v,%v,%v\n", x, from, size, v.Total, resultSize)

	for _, o := range v.Result {
		_, err := util.FileAppendNewLine(output, util.ToJson(o, false))
		if err != nil {
			panic(err)
		}
	}

	if len(v.Result) >= size {
		//continue
		from = from + size
		goto begin
	}
}
