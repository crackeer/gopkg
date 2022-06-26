package util

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// MapContainer
type MapContainer struct {
	value map[string]interface{}
}

// LoadMap
//  @param data
//  @return *MapContainer
func LoadMap(data map[string]interface{}) *MapContainer {
	return &MapContainer{
		value: data,
	}
}

// GetInt64 ...
//  @receiver container
//  @param key
//  @return int64
func (container *MapContainer) GetInt64(key string, defaultValue int64) int64 {
	val, ok := container.value[key]

	if !ok {
		return defaultValue
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
	return defaultValue
}

// GetString
//  @receiver container
//  @param key
//  @return string
func (container *MapContainer) GetString(key string, defaultValue string) string {

	val, ok := container.value[key]

	if !ok {
		return defaultValue
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

// GetBool
//  @receiver container
//  @param key
//  @return bool
func (container *MapContainer) GetBool(key string, defaultValue bool) bool {

	val, ok := container.value[key]
	if !ok || val == nil {
		return defaultValue
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

	return defaultValue
}

// GetMap ValFromMap ...
//  @receiver container
//  @param key
//  @return map
func (container *MapContainer) GetMap(key string) map[string]interface{} {
	val, ok := container.value[key]
	if !ok || val == nil {
		return nil
	}

	if real, ok := val.(map[string]interface{}); ok {
		return real
	}

	return nil
}
