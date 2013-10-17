package util

import (
	"os"
	"io"
log	"github.com/cihub/seelog"
	"io/ioutil"
	"fmt"
	"sort"
)

/**
 * User: Medcl
 * Date: 13-7-22
 * Time: 下午12:23
 */
func CheckFileExists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil { return true }
	if os.IsNotExist(err) { return false }
	return false
}

func CopyFile(src,dst string)(w int64,err error){
	srcFile,err := os.Open(src)
	if err!=nil{
		log.Error(err.Error())
		return
	}
	defer srcFile.Close()

	dstFile,err := os.Create(dst)

	if err!=nil{
		log.Error(err.Error())
		return
	}

	defer dstFile.Close()

	return io.Copy(dstFile,srcFile)
}


// get file modified time
func FileMTime(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.ModTime().Unix(), nil
}

// get file size as how many bytes
func FileSize(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

// delete file
func Unlink(file string) error {
	return os.Remove(file)
}

// rename file name
func Rename(file string, to string) error {
	return os.Rename(file, to)
}

// put string to file
func FilePutContent(file string, content string) (int, error) {
	fs, e := os.Create(file)
	if e != nil {
		return 0, e
	}
	defer fs.Close()
	return fs.WriteString(content)
}

// put string to file
func FilePutContentWithByte(file string, content []byte) (int, error) {
	fs, e := os.Create(file)
	if e != nil {
		return 0, e
	}
	defer fs.Close()
	return fs.Write(content)
}

//func FileAppendContentWithByte(file string,content []byte)(int,error){
//	if(!os.IsExist(file)){
//		fs,e :=os.Create(file)
//	}else{
//		fs,e:=os.OpenFile(file,os.O_APPEND)
//	}
//	f.Seek(stat.Size, 0)
//}

// get string from text file
func FileGetContent(file string) ([]byte, error) {
	if !IsFile(file) {
		return nil, os.ErrNotExist
	}
	b, e := ioutil.ReadFile(file)
	if e != nil {
		return nil, e
	}
	return b, nil
}

// it returns false when it's a directory or does not exist.
func IsFile(file string) bool {
	f, e := os.Stat(file)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// IsExist returns whether a file or directory exists.
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

//create file
func CreateFile(dir string, name string) (string, error) {
	src := dir + name + "/"
	if IsExist(src) {
		return src, nil
	}

	if err := os.MkdirAll(src, 0777); err != nil {
		if os.IsPermission(err) {
			fmt.Println("permission denied")
		}
		return "", err
	}

	return src, nil
}

type FileRepos []Repository

type Repository struct {
	Name     string
	FileTime int64
}

func (r FileRepos) Len() int {
	return len(r)
}

func (r FileRepos) Less(i, j int) bool {
	return r[i].FileTime < r[j].FileTime
}

func (r FileRepos) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// 获取所有文件
//如果文件达到最上限，按时间删除
func delFile(files []os.FileInfo, count int, fileDir string) {
	if len(files) <= count {
		return
	}

	result := new(FileRepos)

	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			*result = append(*result, Repository{Name: file.Name(), FileTime: file.ModTime().Unix()})
		}
	}

	sort.Sort(result)
	deleteNum := len(files) - count
	for k, v := range *result {
		if k+1 > deleteNum {
			break
		}
		Unlink(fileDir + v.Name)
	}

	return
}

//func main() {
//	//生成文件
//	dir := "a:/src/"
//	file, err := CreateFile(dir, "20130829")
//
//	if err != nil {
//		return
//	}
//
//	//写文件
//	content := "teststttttt"
//	l, e := FilePutContent(file+"1.txt", content)
//
//	if e != nil && l <= 0 {
//		return
//	}
//
//	//读文件
//	// str, _ := FileGetContent(file + "1.txt")
//	// fmt.Println("str", str)
//	// size, _ := FileSize(file + "1.txt")
//	// fmt.Println("size", size)
//	// ftime, _ := FileMTime(file + "1.txt")
//	// fmt.Println("ftime", ftime)
//
//	// 获取所有文件
//	//如果文件达到最上限，按时间删除
//	files, _ := ioutil.ReadDir(file)
//	delFile(files, 1, file)
//	fmt.Println("count", len(files))
//}
