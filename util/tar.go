package util

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// TarOption
type TarOption struct {
	IncludeRootDir bool
	EnableGZip     bool
}

// Tar  ...
//
//	@param fileName
//	@param dest
//	@param includeRootDir
//	@return error
func Tar(fileName string, dest string, opt *TarOption) error {
	if opt == nil {
		opt = &TarOption{}
	}
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

	var tarWriter *tar.Writer

	if opt.EnableGZip {
		gzipWriter := gzip.NewWriter(destFile)
		defer gzipWriter.Close()
		tarWriter = tar.NewWriter(gzipWriter)
	} else {
		tarWriter = tar.NewWriter(destFile)
	}
	defer tarWriter.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("file stat error: %v", err)
	}

	prefix := ""
	fileMap := map[string]string{}
	if fileStat.IsDir() {
		fileMap = GetDirFilesAsMap(fileName)
		if opt.IncludeRootDir {
			prefix = fileStat.Name() + "/"
		}
	} else {
		_, name := filepath.Split(fileName)
		fileMap[name] = fileName
	}

	return TarWriteFiles(tarWriter, fileMap, prefix)
}

// TarWriteFile ...
//
//	@param tw
//	@param tarFilePath
//	@param localFilePath
//	@return error
func TarWriteFile(tw *tar.Writer, tarFilePath, localFilePath string) error {

	if tw == nil {
		return fmt.Errorf("tar writer is nil")
	}

	file, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("open file `%s` error: %s", localFilePath, err.Error())
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("get file stat `%s` error: %s", localFilePath, err.Error())
	}
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return fmt.Errorf("tar.FileInfoHeader `%s` error: %s", localFilePath, err.Error())
	}
	header.Name = tarFilePath

	if err := tw.WriteHeader(header); err != nil {
		return fmt.Errorf("tar.WriteHeader `%s` error: %s", localFilePath, err.Error())
	}

	if _, err = io.Copy(tw, file); err != nil {
		return fmt.Errorf("copy file `%s` error: %s", localFilePath, err.Error())
	}
	return nil
}

// TarWriteFiles ...
//
//	@param tw
//	@param fileMap
//	@param prefix
//	@return error
func TarWriteFiles(tw *tar.Writer, fileMap map[string]string, prefix string) error {
	for key, value := range fileMap {
		if err := TarWriteFile(tw, prefix+key, value); err != nil {
			return err
		}
	}
	return nil
}

// UnTar
//
//	@param fullPath
//	@param dest
//	@param opt
//	@return error
func UnTar(fullPath string, dest string, opt *TarOption) error {

	var (
		gzipReader *gzip.Reader
		err        error
		osFile     *os.File
		tarReader  *tar.Reader
	)

	osFile, err = os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("open %s error: %s", fullPath, err.Error())
	}
	defer osFile.Close()

	if opt != nil && opt.EnableGZip {
		gzipReader, err = gzip.NewReader(osFile)
		if err != nil {
			return fmt.Errorf("new gzip reader error [%s] : %s", fullPath, err.Error())
		}
	}

	if gzipReader != nil {
		tarReader = tar.NewReader(gzipReader)
	} else {
		tarReader = tar.NewReader(osFile)
	}

	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return fmt.Errorf("error occur:%s", err.Error())
			}
		}

		curFile := header.FileInfo()
		if curFile.IsDir() {
			continue
		}
		tmpFile, err := createFile(filepath.Join(dest, header.Name))
		if err != nil {
			return fmt.Errorf("create file %s error: %s", filepath.Join(dest, header.Name), err.Error())
		}
		io.Copy(tmpFile, tarReader)
	}

	return nil
}
