package util

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/morysky/melody"
)

// GetStringFromSession Get a string param from websocket session
func GetStringFromSession(wsSession *melody.Session, key string, defaultValue string) string {
	var (
		ok       bool
		value    string
		rawValue interface{}
	)

	if wsSession == nil {
		return defaultValue
	}

	rawValue, ok = wsSession.Get(key)
	if !ok {
		return defaultValue
	}

	value, ok = rawValue.(string)
	if !ok {
		return defaultValue
	}

	return value
}

// GetInt64FromSession Get a int64 param from websocket session
func GetInt64FromSession(wsSession *melody.Session, key string) int64 {
	var (
		ok       bool
		err      error
		value    int64
		rawValue interface{}
		strValue string
	)

	if wsSession == nil {
		return value
	}

	rawValue, ok = wsSession.Get(key)
	if !ok {
		return value
	}

	if strValue, ok = rawValue.(string); ok {
		value, err = strconv.ParseInt(strValue, 10, 64)
		if err == nil {
			return value
		}
	}

	value, _ = rawValue.(int64)

	return value
}

// GetInt64FromMap Get int64 value from a flat map
func GetInt64FromMap(container map[string]interface{}, key string) int64 {
	var ret int64

	if v, exists := container[key]; exists {
		ret, _ = v.(int64)
	}

	return ret
}

//
func GetFloat64FromMap(container map[string]interface{}, key string) float64 {
	var ret float64

	if v, exists := container[key]; exists {
		ret, _ = v.(float64)
	}

	return ret
}

// GetUint32FromMap Get int64 value from a flat map
func GetUint32FromMap(container map[string]interface{}, key string) uint32 {
	var ret uint32

	if v, exists := container[key]; exists {
		ret, _ = v.(uint32)
	}

	return ret
}

// GetDurationFromMap ...
func GetDurationFromMap(container map[string]interface{}, key string) time.Duration {
	var ret time.Duration

	if v, exists := container[key]; exists {
		ret, _ = v.(time.Duration)
	}

	return ret
}

// GetIntFromMap Get int value from a flat map
func GetIntFromMap(container map[string]interface{}, key string) int {
	var ret int

	if container == nil {
		return ret
	}

	if v, exists := container[key]; exists {
		ret, _ = v.(int)
	}

	return ret
}

// GetStringFromMap Get string value from a flat map
func GetStringFromMap(container map[string]interface{}, key string) string {
	var ret string

	if v, exists := container[key]; exists {
		ret = fmt.Sprintf("%v", v)
	}

	return ret
}

// GetBoolFromMap Get boolean value from a flat map
func GetBoolFromMap(container map[string]interface{}, key string) bool {
	var ret bool

	if v, exists := container[key]; exists {
		ret, _ = v.(bool)
	}

	return ret
}

// GetInt64ValFromMap ...
func GetInt64ValFromMap(container map[string]interface{}, key string) int64 {

	val, ok := container[key]

	if !ok {
		return 0
	}

	if v, ok := val.(int64); ok {
		return v
	}

	if v, ok := val.(json.Number); ok {
		val, _ := v.Int64()
		return val
	}

	if v, ok := val.(int); ok {
		return int64(v)
	}

	if v, ok := val.(float64); ok {
		return int64(v)
	}

	if v, ok := val.(string); ok {
		if rv, err := strconv.Atoi(v); err == nil {
			return int64(rv)
		}
	}
	return 0
}

func GetStringValFromMap(container map[string]interface{}, key string) string {

	val, ok := container[key]

	if !ok {
		return ""
	}

	if v, ok := val.(json.Number); ok {
		val, _ := v.Int64()
		return fmt.Sprintf("%d", val)
	}

	if v, ok := val.(int64); ok {
		return fmt.Sprintf("%d", v)
	}

	if v, ok := val.(int); ok {
		return fmt.Sprintf("%d", v)
	}

	if v, ok := val.(float64); ok {
		return strconv.FormatFloat(v, 'f', -1, 64)
	}

	if v, ok := val.(float32); ok {
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	}

	if v, ok := val.(string); ok {
		return v
	}

	bs, _ := json.Marshal(val)
	return string(bs)
}

func ReplaceUrlParameter(urlStr string, params map[string]string) string {

	if !strings.Contains(urlStr, "{") || params == nil {
		return urlStr
	}
	for k, v := range params {
		urlStr = strings.Replace(urlStr, fmt.Sprintf("{%s}", k), v, -1)
	}

	return urlStr
}

func GetBooleanValFromMap(container map[string]interface{}, key string) bool {

	val, ok := container[key]
	if !ok || val == nil {
		return false
	}

	if v, ok := val.(bool); ok {
		return v
	}

	if v, ok := val.(json.Number); ok {
		val, _ := v.Int64()
		return val > 0
	}

	if v, ok := val.(int); ok {
		return v > 0
	}

	if v, ok := val.(float64); ok {

		return v > 0
	}

	if v, ok := val.(float32); ok {
		return v > 0
	}

	if v, ok := val.(string); ok {
		if intVal, err := strconv.Atoi(v); err == nil {
			return intVal > 0
		}
		return strings.ToLower(v) == "yes" || strings.ToLower(v) == "on"
	}

	return false
}

// GetMapValFromMap ...
func GetMapValFromMap(container map[string]interface{}, key string) map[string]interface{} {
	val, ok := container[key]
	if !ok || val == nil {
		return nil
	}

	if real, ok := val.(map[string]interface{}); ok {
		return real
	}

	return nil
}
