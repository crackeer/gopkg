package util

import "encoding/json"

func CloneMap(src map[string]interface{}) map[string]interface{} {

	if bs, err := json.Marshal(src); err == nil {
		retData := map[string]interface{}{}
		if err1 := json.Unmarshal(bs, &retData); err1 == nil {
			return retData
		}
	}

	return nil
}

func ToInterfaceMap(src map[string]string) map[string]interface{} {
	retData := map[string]interface{}{}
	for k, v := range src {
		retData[k] = v
	}
	return retData
}

// Convert2Map
//  @param src
//  @return map
func Convert2Map(src interface{}) map[string]interface{} {
	bytes, _ := Marshal(src)
	retData := map[string]interface{}{}
	Unmarshal(string(bytes), &retData)
	return retData
}

// Map2MapString
//  @param src
//  @return map
func Map2MapString(src map[string]interface{}) map[string]string {
	retData := map[string]string{}
	for k, v := range src {
		retData[k] = ToString(v)
	}
	return retData
}
