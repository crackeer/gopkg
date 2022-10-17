package util

import (
	"bytes"
	"encoding/json"
	"reflect"

	jsoniter "github.com/json-iterator/go"
)

var jsonX = jsoniter.ConfigCompatibleWithStandardLibrary

// GetStructJsonTag ...
//  @param i
//  @return []string
func GetStructJsonTag(i interface{}) []string {

	iType := reflect.TypeOf(i)

	jsonTags := []string{}
	idx := 0
	for idx < iType.NumField() {
		jsonTags = append(jsonTags, iType.Field(idx).Tag.Get("json"))
		idx++
	}

	return jsonTags
}

// StructToMapViaJson ...
//  @param data
//  @return map
func StructToMapViaJson(data interface{}) map[string]interface{} {
	m := make(map[string]interface{})

	j, _ := json.Marshal(data)

	json.Unmarshal(j, &m)

	return m
}

// Unmarshal JSON UnmarshalJSON a string of float64 safe
//  @param input
//  @param output
//  @return error
func Unmarshal(input string, output interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader([]byte(input)))
	decoder.UseNumber()

	return decoder.Decode(output)
}

// Marshal
//  @param dest
//  @return []byte
//  @return error
func Marshal(dest interface{}) ([]byte, error) {
	return json.Marshal(dest)
}

// Marshalx Marshal ...
//  @param dest
//  @return []byte
//  @return error
func Marshalx(dest interface{}) ([]byte, error) {
	return jsonX.Marshal(dest)
}

// Unmarshalx Unmarshal ...
//  @param input
//  @param data
//  @return error
func Unmarshalx(input []byte, data interface{}) error {
	return jsonX.Unmarshal(input, &data)
}

// MarshalEscapeHtml ...
//  @param data
//  @return string
//  @return error
func MarshalEscapeHtml(data interface{}) (string, error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(true)
	if err := jsonEncoder.Encode(data); err != nil {
		return "", err
	}
	return bf.String(), nil
}
