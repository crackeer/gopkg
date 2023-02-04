package router

import (
	"encoding/json"
	"fmt"
	"os"
)

// ParseRouterMetaByFile
//
//	@param fullPath
//	@return *RouterMeta
//	@return error
func ParseRouterMetaByFile(fullPath string) (*RouterMeta, error) {
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
