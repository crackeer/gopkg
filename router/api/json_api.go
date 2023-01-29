package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/crackeer/gopkg/config"
)

const (
	jsonExtension = ".json"
)

// JSONAPI
type JsonObject struct {
	ServiceMap     map[string]JsonServiceItem `json:"service_map"`
	SuccessCodeKey string                     `json:"success_code_key"`
	MessageKey     string                     `json:"message_key"`
	SuccessCode    string                     `json:"success_code"`
	DataKey        string                     `json:"data_key"`
	Timeout        int64                      `json:"timeout"`
	APIList        []JsonAPIItem              `json:"api_list"`
}

// JsonServiceItem
type JsonServiceItem struct {
	Host       string      `json:"host"`
	SignConfig *SignConfig `json:"sign_config"`
}

// JsonAPIItem
type JsonAPIItem struct {
	Name        string `json:"name"`
	URI         string `json:"uri"`
	Method      string `json:"method"`
	ContentType string `json:"content_type"`
}

// JSONAPIMeta
type JSONAPIMeta struct {
	Path      string
	container *sync.Map
}

// NewYamlAPIMeta
//
//	@param prefix
//	@return *YamlAPIMeta
func NewJSONAPIMeta(prefix string) *JSONAPIMeta {
	return &JSONAPIMeta{
		Path:      prefix,
		container: new(sync.Map),
	}
}

// GetAPIMeta
//
//	@receiver apiMetaGetter
//	@param name
//	@param env
//	@return *api.APIMeta
//	@return error
func (apiMetaGetter *JSONAPIMeta) GetAPIMeta(name string, env string) (*APIMeta, error) {

	key := name + "@" + env
	if value, ok := apiMetaGetter.container.Load(key); ok {
		if apiMeta, ok := value.(*APIMeta); ok {
			return apiMeta, nil
		}
	}

	apiMeta, err := apiMetaGetter.loadAPIMeta(name, env)

	if err != nil {
		return nil, err
	}

	apiMetaGetter.container.Store(key, apiMeta)

	return apiMeta, nil
}

func (apiMetaGetter *JSONAPIMeta) LoadAllAPI() error {
	return nil
}

func (apiMetaGetter *JSONAPIMeta) readFile(name string) (*JsonObject, error) {
	path := filepath.Join(apiMetaGetter.Path, name)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	retData := &JsonObject{}
	if err := json.Unmarshal(bytes, retData); err != nil {
		return nil, err
	}
	return retData, nil
}

func (apiMetaGetter *JSONAPIMeta) loadAPIMeta(name string, env string) (*APIMeta, error) {

	parts := strings.Split(name, "/")
	if len(parts) < 2 {
		return nil, errors.New("name error")
	}
	apiName := parts[1]
	path := filepath.Join(apiMetaGetter.Path, parts[0]+".json")
	apiConfig, err := config.LoadYamlAPIConfig(path)
	if err != nil {
		return nil, fmt.Errorf("load api `%s` config error: %s", name, err.Error())
	}

	if len(apiConfig.BaseURI) < 1 {
		return nil, errors.New("get api error: base_uri list nil")
	}

	if len(apiConfig.API) < 1 {
		return nil, errors.New("get api error: api list nil")
	}

	item, exists := apiConfig.API[apiName]
	if !exists {
		return nil, fmt.Errorf("api `%s` not exists", apiName)
	}

	retData := &APIMeta{
		SuccessCode: apiConfig.SuccessCode,
		CodeKey:     apiConfig.CodeKey,
		DataKey:     apiConfig.DataKey,
		MessageKey:  apiConfig.MessageKey,
		Timeout:     apiConfig.Timeout,
		Path:        item.Path,
		Method:      item.Method,
		ContentType: item.ContentType,
		Header:      map[string]string{},
	}

	if baseURI, exists := apiConfig.BaseURI[config.EnvDefault]; exists {
		retData.BaseURI = baseURI
	}

	if baseURI, exists := apiConfig.BaseURI[env]; exists {
		retData.BaseURI = baseURI
	}

	if len(retData.BaseURI) < 1 {
		return nil, errors.New("api base_uri nil")
	}

	for key, value := range apiConfig.Header {
		retData.Header[key] = value
	}

	for key, value := range item.Header {
		retData.Header[key] = value
	}

	return retData, nil
}
