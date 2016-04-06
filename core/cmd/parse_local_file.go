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
	"container/list"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

var outputFileName string = "filesName.csv"
var baseDomain string = "http://www.baidu.com/"

func checkErr(err error) {
	if nil != err {
		panic(err)
	}
}

func GetFullPath(path string) string {
	absolutePath, _ := filepath.Abs(path)
	return absolutePath
}

func PrintFilesName(path string) {
	fullPath := GetFullPath(path)

	listStr := list.New()

	filepath.Walk(fullPath, func(path string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		name := fi.Name()
		if outputFileName != name {
			listStr.PushBack(path)
		}

		return nil
	})

	outputFilesName(listStr)
}

func convertToSlice(listStr *list.List) []string {
	sli := []string{}
	for el := listStr.Front(); nil != el; el = el.Next() {
		sli = append(sli, el.Value.(string))
	}

	return sli
}

func outputFilesName(listStr *list.List) {
	files := convertToSlice(listStr)
	//sort.StringSlice(files).Sort()// sort

	f, err := os.Create(outputFileName)
	//	f, err := os.OpenFile(outputFileName, os.O_APPEND | os.O_CREATE, os.ModeAppend)
	checkErr(err)
	defer f.Close()

	writer := csv.NewWriter(f)

	length := len(files)
	for i := 0; i < length; i++ {
		writer.Write([]string{baseDomain + "|||" + files[i]})
	}

	writer.Flush()
}

func main() {
	var path string
	if len(os.Args) > 3 {
		baseDomain = os.Args[1]
		path = os.Args[2]
		outputFileName = os.Args[3]
	} else if len(os.Args) > 2 {
		baseDomain = os.Args[1]
		path = os.Args[2]
	} else if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path, _ = os.Getwd()
	}

	fmt.Println("baseDomain:" + baseDomain)
	fmt.Println("path:" + path)
	fmt.Println("output:" + outputFileName)

	PrintFilesName(path)
	fmt.Println("done!")
}
