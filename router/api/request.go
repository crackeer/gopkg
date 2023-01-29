package api

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/crackeer/gopkg/mapbuilder"
	"github.com/crackeer/gopkg/util"
)

const rootKey = "_root"
const headerKey = "_header"

// RequestClient ...
type RequestClient struct {
	factory APIMetaGetter
	logger  Logger
}

// RequestItem ...
type RequestItem struct {
	API    string                 `yaml:"api"`
	Params map[string]interface{} `yaml:"params"`
	Header map[string]string      `yaml:"header"`
	As     string                 `yaml:"as"`
	Key    bool                   `yaml:"key"`
}

// ParseMeshConfig
//
//	@param raw
//	@return [][]*RequestItem
//	@return error
func ParseMeshConfig(raw string) ([][]*RequestItem, error) {
	retData := [][]*RequestItem{}
	if err := json.Unmarshal([]byte(raw), &retData); err != nil {
		return nil, err
	}
	return retData, nil
}

// NewRequestClient
//
//	@param getter
//	@return *RequestClient
func NewRequestClient(getter APIMetaGetter) *RequestClient {
	return &RequestClient{
		factory: getter,
	}
}

// Request
//
//	@receiver client
//	@param name
//	@param query
//	@param header
//	@return *APIResponse
//	@return error
func (client *RequestClient) Request(apiID string, query map[string]interface{}, header map[string]string, env string) (*APIResponse, error) {
	apiMeta, err := client.factory.GetAPIMeta(apiID, env)
	if err != nil {
		return nil, fmt.Errorf("get api meta error: %s", err.Error())
	}

	apiRequest := NewAPIRequest(apiMeta, client.logger)
	return apiRequest.Do(query, header)
}

// RequestList
//
//	@receiver client
//	@param list
//	@return map[string]*APIResponse
//	@return error
func (client *RequestClient) requestList(list []*RequestItem, env string) (map[string]*APIResponse, map[string]string, error) {

	wg := &sync.WaitGroup{}
	wg.Add(len(list))
	locker := &sync.RWMutex{}
	var (
		retError error
		errMap   map[string]string       = make(map[string]string)
		retData  map[string]*APIResponse = make(map[string]*APIResponse)
	)

	for _, item := range list {
		go func(tmp *RequestItem) {
			apiResponse, err := client.Request(tmp.API, tmp.Params, tmp.Header, env)
			locker.Lock()
			defer locker.Unlock()
			if err != nil {
				errMap[tmp.As] = err.Error()
				if tmp.Key {
					retError = err
				}
			} else {
				retData[tmp.As] = apiResponse
			}
			wg.Done()

		}(item)
	}
	wg.Wait()
	return retData, errMap, retError
}

// RequestList
//
//	@receiver client
//	@param list
//	@return map[string]*APIResponse
//	@return error
func (client *RequestClient) Mesh(list [][]*RequestItem, query map[string]interface{}, header map[string]string, env string) (map[string]*APIResponse, map[string]string, error) {

	var (
		retError error
		errMap   map[string]string       = make(map[string]string)
		retData  map[string]*APIResponse = make(map[string]*APIResponse)
	)

	input := map[string]interface{}{
		rootKey:   query,
		headerKey: header,
	}

	for _, items := range list {
		newItems, err := transfer(items, input)
		if err != nil {
			retError = err
			break
		}
		mapAPIResponse, mapError, err := client.requestList(newItems, env)
		if err != nil {
			retError = err
			break
		}

		for key, response := range mapAPIResponse {
			retData[key] = response
			var realData interface{}
			if err := util.Unmarshal(string(response.OriginBody), &realData); err == nil {
				input[key] = realData
			} else {
				input[key] = string(response.OriginBody)
			}
		}

		for key, message := range mapError {
			errMap[key] = message
		}

	}
	return retData, errMap, retError
}

func transfer(list []*RequestItem, input map[string]interface{}) ([]*RequestItem, error) {
	src := util.ToString(input)

	var err error
	for i, item := range list {
		if len(item.Params) > 0 {
			if list[i].Params, err = mapbuilder.Build(src, item.Params); err != nil && item.Key {
				return list, fmt.Errorf("build `%s` params error: %s", item.As, err.Error())
			}
		} else {
			list[i].Params = util.LoadMap(input).GetMap(rootKey)
		}

		if len(item.Params) > 0 {
			header := util.Convert2Map(item.Header)

			newHeader, err := mapbuilder.Build(src, header)
			if err != nil && item.Key {
				return list, fmt.Errorf("build `%s` header error: %s", item.As, err.Error())
			}
			list[i].Header = util.Map2MapString(newHeader)
		} else {
			header := util.LoadMap(input).GetMap(headerKey)
			list[i].Header = util.Map2MapString(header)
		}
	}
	return list, nil
}
