/**
 * Created with IntelliJ IDEA.
 * User: medcl
 * Date: 13-10-28
 * Time: 上午11:49
 */
package main


import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"container/list"
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
			if outputFileName != name{
				listStr.PushBack(path)
			}

			return nil
		})

	outputFilesName(listStr)
}

func convertToSlice(listStr *list.List)[]string{
	sli := []string{}
	for el:= listStr.Front(); nil != el; el= el.Next(){
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
	for i:= 0; i < length; i++{
		writer.Write([]string{baseDomain+"|||"+files[i]})
	}

	writer.Flush()
}

func main() {
	var path string
	if len(os.Args) > 3 {
		baseDomain = os.Args[1]
		path = os.Args[2]
		outputFileName = os.Args[3]
	}else 	if len(os.Args) > 2 {
		baseDomain = os.Args[1]
		path = os.Args[2]
	}else if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path, _ = os.Getwd()
	}

	fmt.Println("baseDomain:"+baseDomain)
	fmt.Println("path:"+path)
	fmt.Println("output:"+outputFileName)


	PrintFilesName(path)
	fmt.Println("done!")
}
