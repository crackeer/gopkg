// Package util...
//
// File:  file
//
// Description: file
//
// Date: 2021/4/11 下午9:21
package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FileExists @param filename
//
//	@param filename
//	@return bool
func FileExists(filename string) bool {
	_, err := os.Stat(filename)

	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// IsDir @return bool
//
//	@param filename
//	@return bool
//	@return error
func IsDir(filename string) (bool, error) {
	fd, err := os.Stat(filename)

	if err != nil {
		return false, err
	}

	fm := fd.Mode()

	return fm.IsDir(), nil
}

// GetFiles
//
//	@param folder
//	@return []string
func GetFiles(folder string) []string {
	files, _ := ioutil.ReadDir(folder)
	retData := []string{}
	for _, file := range files {
		if file.IsDir() {
			tmp := GetFiles(fmt.Sprintf("%s/%s", folder, file.Name()))
			retData = append(retData, tmp...)
		} else {
			retData = append(retData, fmt.Sprintf("%s/%s", folder, file.Name()))
		}
	}
	return retData
}

// GetDirFilesAsMap
//
//	@param fileDir
//	@return map
func GetDirFilesAsMap(fileDir string) map[string]string {
	fileList := GetFiles(fileDir)
	if len(fileList) < 1 {
		return nil
	}
	fileMap := map[string]string{}
	for _, file := range fileList {
		key := file[len(fileDir)+1:]
		fileMap[key] = file
	}
	return fileMap
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}

// ReadSubFiles //
//
//	@param folder
//	@return []string
func ReadSubFiles(folder string) []string {
	files, _ := ioutil.ReadDir(folder)
	retData := []string{}
	for _, file := range files {
		if !file.IsDir() {
			retData = append(retData, filepath.Join(folder, file.Name()))
		}
	}
	return retData
}
