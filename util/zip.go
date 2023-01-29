package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Zip ZipArchive  ...
//
//	@param fileName
//	@param dest
//	@param includeRootDir
//	@return error
func Zip(fileName string, dest string, includeRootDir bool) error {

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	zipWriter := zip.NewWriter(destFile)
	defer zipWriter.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("file stat error: %v", err)
	}

	prefix := ""
	fileMap := map[string]string{}
	if fileStat.IsDir() {
		fileMap = GetDirFilesAsMap(fileName)
		if includeRootDir {
			prefix = fileStat.Name() + "/"
		}
	} else {
		_, name := filepath.Split(fileName)
		fileMap[name] = fileName
	}

	return ZipWriteFiles(zipWriter, fileMap, prefix)
}

// ZipWriteFile
//
//	@param writer
//	@param zipFileName
//	@param localFilePath
//	@return error
func ZipWriteFile(writer *zip.Writer, zipFileName, localFilePath string) error {

	if writer == nil {
		return fmt.Errorf("tar writer is nil")
	}

	file, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("open file `%s` error: %s", localFilePath, err.Error())
	}
	defer file.Close()
	tmpWriter, err := writer.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("zipWriter create zip file `%s` error: %s", localFilePath, err.Error())
	}
	if _, err = io.Copy(tmpWriter, file); err != nil {
		return fmt.Errorf("copy file `%s` error: %s", localFilePath, err.Error())
	}
	return nil
}

// ZipWriteFiles
//
//	@param writer
//	@param fileMap
//	@param prefix
//	@return error
func ZipWriteFiles(writer *zip.Writer, fileMap map[string]string, prefix string) error {
	for key, value := range fileMap {
		if err := ZipWriteFile(writer, prefix+key, value); err != nil {
			return err
		}
	}
	return nil
}
