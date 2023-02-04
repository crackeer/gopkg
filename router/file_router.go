package router

import (
	"fmt"
	"os"
	"sync"

	"github.com/crackeer/gopkg/util"
)

// FileRouter
type FileRouter struct {
	BasePath  string
	container *sync.Map
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
		BasePath:  basePath,
		container: new(sync.Map),
	}, nil
}

// GetRouterMeta
//
//	@receiver f
//	@param path
//	@return *RouterMeta
func (f *FileRouter) Get(uri string) *RouterMeta {
	value, ok := f.container.Load(uri)
	if !ok {
		return nil
	}

	return value.(*RouterMeta)
}

func (f *FileRouter) LoadAll() error {
	files := util.GetDirFilesAsMap(f.BasePath)
	for uri, file := range files {
		meta, err := ParseRouterMetaByFile(file)
		if err != nil {
			return fmt.Errorf("parse file `%s` as router `%s` error:%s", file, uri, err.Error())
		}
		f.container.Store(uri, meta)
	}
	return nil
}
