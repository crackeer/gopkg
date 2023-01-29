package router

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// FileRouter
type FileRouter struct {
	BasePath string
}

// NewFileRouter
//
//	@param basePath
//	@return *FileRouter
//	@return error
func NewFileRouter(basePath string) (*FileRouter, error) {
	file, err := os.Stat(basePath)
	if err != nil {
		return nil, err
	}
	if !file.IsDir() {
		return nil, fmt.Errorf("`%s` is not a directory", basePath)
	}
	return &FileRouter{
		BasePath: basePath,
	}, nil
}

// GetRouterMeta
//
//	@receiver f
//	@param path
//	@return *RouterMeta
func (f *FileRouter) GetRouterMeta(path string) (*RouterMeta, error) {
	fullPath := filepath.Join(f.BasePath, path)
	bytes, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("read file `%s` error:%s", fullPath, err.Error())
	}
	retData := &RouterMeta{}
	if err := json.Unmarshal(bytes, retData); err != nil {
		return nil, fmt.Errorf("json unmarshal `%s` content error:%s", fullPath, err.Error())
	}

	return retData, nil
}
