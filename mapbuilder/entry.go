package mapbuilder

import (
	"strings"

	"github.com/crackeer/gopkg/util"
)

// Build
//  @param sourceData
//  @param structData
//  @return map[string]interface{}
//  @return error
func Build(src interface{}, structData map[string]interface{}) (map[string]interface{}, error) {

	raws := util.ToString(src)
	bytes := []byte(raws)
	retData := map[string]interface{}{}
	for key, val := range structData {
		if rv := get(bytes, val); rv != nil {
			retData[key] = rv
		}
	}

	return retData, nil
}

func get(input []byte, structObj interface{}) interface{} {

	if key, ok := structObj.(string); ok {
		if len(key) > 0 && key[0] != '@' {
			return key
		}
		return MutiGjsonGet(input, strings.Split(key, _sep))
	}

	retData := map[string]interface{}{}
	mapStruct, ok := structObj.(map[string]interface{})

	if !ok {
		return structObj
	}

	for k, v := range mapStruct {
		if val := get(input, v); val != nil {
			retData[k] = val
		}
	}
	return retData
}
