package mapbuilder

import (
	"strings"

	"github.com/crackeer/gopkg/util"
	"github.com/tidwall/gjson"
)

// MutiGjsonGet
//
//	@param input
//	@param patterns
//	@return interface{}
func MutiGjsonGet(input []byte, patterns []string, pre string) interface{} {

	for _, k := range patterns {
		if val := GjsonGet(input, k, pre); val != nil {
			return val
		}
	}
	return nil
}

// GjsonGet get JSON value from response map
//
//	@param input
//	@param pattern
//	@return interface{}
func GjsonGet(input []byte, pattern string, pre string) interface{} {

	if len(pattern) < 1 {
		return nil
	}

	if !strings.HasPrefix(pattern, pre) {
		return pattern
	}

	if pattern == _defaultMap {
		return map[string]string{}
	}
	if pattern == _defaultString {
		return ""
	}
	if pattern == _defaultArray {
		return []interface{}{}
	}

	if pattern == _defaultNil {
		return nil
	}

	pattern1 := strings.TrimPrefix(pattern, pre)

	gr := gjson.GetBytes(input, pattern1)

	if !gr.Exists() {
		return nil
	}

	if gjson.String == gr.Type && gr.String() == nilString {
		return nil
	}

	//直接使用data.Value()有可能会丢失精度
	if gr.Type == gjson.JSON {
		var object interface{}
		if err := util.Unmarshal(gr.String(), &object); err == nil {
			return object
		}
	}

	if gr.Type == gjson.Number {
		if !strings.Contains(gr.String(), ".") {
			return gr.Int()
		}
	}

	return gr.Value()
}
