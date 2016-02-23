/**
 * Created with IntelliJ IDEA.
 * User: medcl
 * Date: 13-10-28
 * Time: 上午11:44
 */
package tasks


import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"container/list"
)
var outputFileName string = "filesName.csv"

func CheckErr(err error) {
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
				listStr.PushBack(name)
			}

			return nil
		})

	OutputFilesName(listStr)
}

func ConvertToSlice(listStr *list.List)[]string{
	sli := []string{}
	for el:= listStr.Front(); nil != el; el= el.Next(){
		sli = append(sli, el.Value.(string))
	}

	return sli
}

func OutputFilesName(listStr *list.List) {
	files := ConvertToSlice(listStr)
	//sort.StringSlice(files).Sort()// sort

	f, err := os.Create(outputFileName)
	//f, err := os.OpenFile(outputFileName, os.O_APPEND | os.O_CREATE, os.ModeAppend)
	CheckErr(err)
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF")
	writer := csv.NewWriter(f)

	length := len(files)
	for i:= 0; i < length; i++{
		writer.Write([]string{files[i]})
	}

	writer.Flush()
}

func main() {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path, _ = os.Getwd()
	}
	PrintFilesName(path)

	fmt.Println("done!")
}
