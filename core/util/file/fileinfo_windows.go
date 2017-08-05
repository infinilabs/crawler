package file

//https://github.com/elastic/beats/blob/master/LICENSE
//Licensed under the Apache License, Version 2.0 (the "License");
import (
	"os"
)

func stat(name string, statFunc func(name string) (os.FileInfo, error)) (FileInfo, error) {
	info, err := statFunc(name)
	if err != nil {
		return nil, err
	}

	return fileInfo{FileInfo: info}, nil
}
