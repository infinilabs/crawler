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

package util

import (
	"fmt"
	log "github.com/cihub/seelog"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// JoinPath return joined file path
func JoinPath(filenames ...string) string {

	hasSlash := false
	result := ""
	for _, str := range filenames {
		currentHasSlash := false
		if len(result) > 0 {
			currentHasSlash = strings.HasPrefix(str, "/")
			if hasSlash && currentHasSlash {
				str = strings.TrimLeft(str, "/")
			}
			if !(hasSlash || currentHasSlash) {
				str = "/" + str
			}
		}
		hasSlash = strings.HasSuffix(str, "/")
		result += str
	}
	return result
}

// FileExists check if the path are exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// CopyFile copy file from src to dst
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

// FileMTime get file modified time
func FileMTime(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.ModTime().Unix(), nil
}

// FileSize get file size as how many bytes
func FileSize(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

// FileDelete delete file
func FileDelete(file string) error {
	return os.Remove(file)
}

// Rename handle file rename
func Rename(file string, to string) error {
	return os.Rename(file, to)
}

// FilePutContent put string to file
func FilePutContent(file string, content string) (int, error) {
	fs, e := os.Create(file)
	if e != nil {
		return 0, e
	}
	defer fs.Close()
	return fs.WriteString(content)
}

// FilePutContentWithByte put string to file
func FilePutContentWithByte(file string, content []byte) (int, error) {
	fs, e := os.Create(file)
	if e != nil {
		return 0, e
	}
	defer fs.Close()
	return fs.Write(content)
}

// FileAppendContentWithByte append bytes to the end of the file
func FileAppendContentWithByte(file string, content []byte) (int, error) {

	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	return f.Write(content)
}

// FileAppendNewLine append new line to the end of the file
func FileAppendNewLine(file string, content string) (int, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	return f.WriteString(content + "\n")
}

// FileAppendNewLineWithByte append bytes and break line(\n) to the end of the file
func FileAppendNewLineWithByte(file string, content []byte) (int, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	return f.WriteString(string(content) + "\n")
}

// FileGetContent get string from text file
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

// IsFile returns false when it's a directory or does not exist.
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

//CreateFile create file
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
