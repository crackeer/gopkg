package mapbuilder

import (
	"encoding/json"
	"strings"
)

const (
	DefaultSeperator = "?"
	DefaultPrefix    = "@"
)

// ResponseBuildClient response struct build client
type MapBuilder struct {
	Src       interface{}
	bytes     []byte
	seperator string
	prefix    string
}

func MapBuilderFrom(src interface{}) (*MapBuilder, error) {
	retData := &MapBuilder{
		Src:       src,
		seperator: DefaultSeperator,
		prefix:    DefaultPrefix,
	}
	bytes, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	retData.bytes = bytes
	return retData, nil
}

// UseSeperator
//
//	@receiver builder
//	@param seq
//	@return *MapBuilder
func (builder *MapBuilder) UseSeperator(seq string) *MapBuilder {
	builder.seperator = seq
	return builder
}

// UsePre
//
//	@receiver builder
//	@param pre
//	@return *MapBuilder
func (builder *MapBuilder) UsePre(pre string) *MapBuilder {
	builder.prefix = pre
	return builder
}

// Build
//
//	@receiver builder
//	@param dest
//	@return interface{}
func (builder *MapBuilder) Build(dest interface{}) interface{} {
	return Build(builder.bytes, dest, builder.prefix, builder.seperator)
}

func Build(input []byte, structObj interface{}, pre string, seperator string) interface{} {

	if key, ok := structObj.(string); ok {
		if !strings.HasPrefix(key, pre) {
			return key
		}
		return MutiGjsonGet(input, strings.Split(key, seperator), pre)
	}

	retData := map[string]interface{}{}
	mapStruct, ok := structObj.(map[string]interface{})
	if !ok {
		return structObj
	}

	for k, v := range mapStruct {
		if val := Build(input, v, pre, seperator); val != nil {
			retData[k] = val
		}
	}
	return retData
}
