package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/crackeer/gopkg/util"
)

const (
	jsonExtension = ".json"
)

// JSONAPI
type JsonObject struct {
	ServiceMap     map[string]JsonServiceItem `json:"service_map"`
	SuccessCodeKey string                     `json:"success_code_key"`
	MessageKey     string                     `json:"message_key"`
	CodeKey        string                     `json:"code_key"`
	SuccessCode    string                     `json:"success_code"`
	DataKey        string                     `json:"data_key"`
	APIList        []JsonAPIItem              `json:"api_list"`
}

// JsonServiceItem
type JsonServiceItem struct {
	Host       string      `json:"host"`
	Name       string      `json:"name"`
	SignConfig *SignConfig `json:"sign_config"`
	Timeout    int64       `json:"timeout"`
}

// JsonAPIItem
type JsonAPIItem struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Method      string `json:"method"`
	ContentType string `json:"content_type"`
}

// JSONAPI
type JSONAPI struct {
	Path      string
	container *sync.Map
}

// NewYamlAPIMeta
//
//	@param prefix
//	@return *YamlAPIMeta
func NewJSONAPI(prefix string) (*JSONAPI, error) {
	return &JSONAPI{
		Path:      prefix,
		container: new(sync.Map),
	}, nil
}

// GetAPIMeta
//
//	@receiver apiMetaGetter
//	@param name
//	@param env
//	@return *api.APIMeta
//	@return error
func (meta *JSONAPI) Get(name string, env string) *APIMeta {

	key := GetCacheKeyByAPI(name, env)
	fmt.Println(key)
	if value, ok := meta.container.Load(key); ok {
		return value.(*APIMeta)
	}

	key = GetCacheKeyByAPI(name, "default")
	fmt.Println(key)
	if value, ok := meta.container.Load(key); ok {
		return value.(*APIMeta)
	}

	return nil
}

func (meta *JSONAPI) LoadAll() error {
	files := meta.loadFiles()
	for name, item := range files {
		object, err := parseFile(item)
		if err != nil {
			return fmt.Errorf("parse file `%s` error: %s", item, err.Error())
		}

		for _, api := range object.APIList {
			for env, service := range object.ServiceMap {
				apiMeta := &APIMeta{
					SuccessCode: object.SuccessCode,
					CodeKey:     object.CodeKey,
					DataKey:     object.DataKey,
					MessageKey:  object.MessageKey,
					Timeout:     service.Timeout,
					Path:        api.Path,
					Method:      api.Method,
					ContentType: api.ContentType,
					Name:        api.Name,
				}
				apiMeta.Host = service.Host
				apiMeta.SignConfig = service.SignConfig
				cacheKey := FormCacheKey(name, api.Name, env)
				meta.container.Store(cacheKey, apiMeta)
			}

		}

	}
	return nil
}
func (meta *JSONAPI) loadFiles() map[string]string {
	files := util.ReadSubFiles(meta.Path)
	retData := map[string]string{}
	for _, file := range files {
		if !strings.HasSuffix(file, jsonExtension) {
			continue
		}
		_, fileName := filepath.Split(file)
		serviceName := strings.TrimSuffix(fileName, jsonExtension)
		retData[serviceName] = file
	}
	return retData
}

func parseFile(filePath string) (*JsonObject, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	retData := &JsonObject{}
	if err := json.Unmarshal(bytes, retData); err != nil {
		return nil, err
	}
	return retData, nil
}
