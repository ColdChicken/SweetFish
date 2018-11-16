package common

import (
	"io/ioutil"
	"os"
	"path"
)

func StringInSlice(s string, ss []string) bool {
	for idx := range ss {
		if s == ss[idx] {
			return true
		}
	}
	return false
}

func Int64InSlice(i int64, ii []int64) bool {
	for idx := range ii {
		if i == ii[idx] {
			return true
		}
	}
	return false
}

func Int64SliceEqual(a []int64, b []int64) bool {
	if len(a) != len(b) {
		return false
	}

	for idx := range a {
		av := a[idx]
		if Int64InSlice(av, b) == false {
			return false
		}
	}

	return true
}

func Mkdir(fullPath string) error {
	err := os.MkdirAll(fullPath, 0777)
	return err
}

func GetDirAllFiles(dir string) ([]string, error) {
	files := []string{}
	directory, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range directory {
		if f.IsDir() {
			files_, err_ := GetDirAllFiles(path.Join(dir, f.Name()))
			if err_ != nil {
				return nil, err_
			}
			files = append(files, files_...)
		} else {
			files = append(files, f.Name())
		}
	}

	return files, nil
}
