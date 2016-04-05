/**
 * User: Medcl
 * Date: 13-7-22
 * Time: 下午12:23
 */
package util

import (
	"bufio"
	"fmt"
	log "github.com/cihub/seelog"
	"io"
	"io/ioutil"
	"os"
	"sort"
)

func CheckFileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func CopyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)

	if err != nil {
		log.Error(err.Error())
		return
	}

	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
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

func FileAppendContentWithByte(file string, content []byte) (int, error) {

	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	return f.Write(content)
}

func FileAppendNewLine(file string, content string) (int, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	return f.WriteString(content + "\n")
}
func FileAppendNewLineWithByte(file string, content []byte) (int, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	return f.WriteString(string(content) + "\n")
}

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

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func ReadAllLines(file string) []string {
	lines := []string{}
	f, err := os.Open(file)
	if err != nil {
		log.Error("error opening file,", file, " ", err)
		os.Exit(1)
	}

	r := bufio.NewReader(f)
	s, e := Readln(r)
	lines = append(lines, s)
	for e == nil {
		s, e = Readln(r)
		if s != "" {
			lines = append(lines, s)
		}
	}

	return lines
}
