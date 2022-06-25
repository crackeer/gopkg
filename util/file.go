// Package util...
//
// File:  file
//
// Description: file
//
// Date: 2021/4/11 下午9:21
package util

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// FileExists 文件或者文件夹是否存在
//
// Author : fuhaixu@ke.com<付海旭>
//
// Date : 下午9:22 2021/4/11
func FileExists(filename string) bool {
	_, err := os.Stat(filename)

	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// IsDir 文件是否
//
// Author : fuhaixu@ke.com<付海旭>
//
// Date : 下午9:31 2021/4/11
func IsDir(filename string) (bool, error) {
	fd, err := os.Stat(filename)

	if err != nil {
		return false, err
	}

	fm := fd.Mode()

	return fm.IsDir(), nil
}

func Unzip(fullPath string, dest string) error {
	// 解压需要使用tar.NewReader方法, 这个方法接收一个io.Reader对象
	// 那边怎么从源文件得到io.Reader对象呢？
	// 这边通过os.Open打开文件,会得到一个os.File对象，
	// 因为他实现了io.Reader的Read方法，所有可以直接传递给tar.NewReader
	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()
	gr, err := gzip.NewReader(file)

	if err != nil {
		return err
	}
	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		curFile := hdr.FileInfo()
		if curFile.IsDir() {
			continue
		}
		filename := dest + hdr.Name
		tmpFile, err := createFile(filename)
		if err != nil {
			return err
		}
		io.Copy(tmpFile, tr)
	}

	return nil
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}

// Compress ...
func Compress(fileName string, dest string) error {

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()
	gw := gzip.NewWriter(destFile)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	err = compress(file, "", tw)
	if err != nil {
		return err
	}
	return nil
}

func compress(file *os.File, prefix string, tw *tar.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, tw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(info, "")
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetFiles
//  @param folder
//  @return []string
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
