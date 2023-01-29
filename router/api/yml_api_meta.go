package api

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/crackeer/gopkg/config"
)

// YamlAPIMetaGetter
type YamlAPIMetaGetter struct {
	ConfigPrefix string
	container    *sync.Map
}

// NewYamlAPIMetaGetter
//
//	@param prefix
//	@return *YamlAPIMetaGetter
func NewYamlAPIMetaGetter(prefix string) *YamlAPIMetaGetter {
	return &YamlAPIMetaGetter{
		ConfigPrefix: prefix,
		container:    new(sync.Map),
	}
}

// GetAPIMeta
//
//	@receiver apiMetaGetter
//	@param name
//	@param env
//	@return *api.APIMeta
//	@return error
func (apiMetaGetter *YamlAPIMetaGetter) GetAPIMeta(name string, env string) (*APIMeta, error) {

	key := name + "@" + env
	fmt.Println(key)
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

func (apiMetaGetter *YamlAPIMetaGetter) loadAPIMeta(name string, env string) (*APIMeta, error) {

	parts := strings.Split(name, "/")
	fmt.Println(parts)
	if len(parts) < 2 {
		return nil, errors.New("name error")
	}
	apiName := parts[1]
	path := apiMetaGetter.ConfigPrefix + "/" + parts[0]
	apiConfig, err := config.LoadYamlAPIConfig(path)
	fmt.Println(apiConfig)
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
		SignAlg:     apiConfig.SignAlg,
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
