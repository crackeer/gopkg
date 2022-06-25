package api

import (
	"fmt"
	"sync"

	"github.com/crackeer/gopkg/mapbuilder"
	"github.com/crackeer/gopkg/util"
)

// RequestClient
type RequestClient struct {
	factory APIMetaGetter
	logger  Logger
}

// RequestItem
type RequestItem struct {
	APIID  string                 `json:"api_id"`
	Params map[string]interface{} `json:"params"`
	Header map[string]string      `json:"header"`
	As     string                 `json:"as"`
	Key    bool                   `json:"key"`
}

// NewRequestClient
//  @param getter
//  @return *RequestClient
func NewRequestClient(getter APIMetaGetter) *RequestClient {
	return &RequestClient{
		factory: getter,
	}
}

// Request
//  @receiver client
//  @param name
//  @param query
//  @param header
//  @return *APIResponse
//  @return error
func (client *RequestClient) Request(apiID string, query map[string]interface{}, header map[string]string) (*APIResponse, error) {
	apiMeta, err := client.factory.GetAPIMeta(apiID)
	if err != nil {
		return nil, fmt.Errorf("get api meta error: %s", err.Error())
	}

	apiRequest := NewAPIRequest(apiMeta, client.logger)
	return apiRequest.Do(query, header)
}

// RequestList
//  @receiver client
//  @param list
//  @return map[string]*APIResponse
//  @return error
func (client *RequestClient) RequestList(list []*RequestItem) (map[string]*APIResponse, map[string]string, error) {

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
			apiResponse, err := client.Request(tmp.APIID, tmp.Params, tmp.Header)
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
//  @receiver client
//  @param list
//  @return map[string]*APIResponse
//  @return error
func (client *RequestClient) Mesh(list [][]*RequestItem, query map[string]interface{}, header map[string]string) (map[string]*APIResponse, map[string]string, error) {

	var (
		retError error
		errMap   map[string]string       = make(map[string]string)
		retData  map[string]*APIResponse = make(map[string]*APIResponse)
	)

	input := map[string]interface{}{
		"root":   query,
		"header": header,
	}

	for _, items := range list {
		newItems, err := transfer(items, input)
		if err != nil {
			retError = err
			break
		}
		mapAPIResponse, mapError, err := client.RequestList(newItems)
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
		if list[i].Params, err = mapbuilder.Build(src, item.Params); err != nil && item.Key {
			return list, fmt.Errorf("build `%s` params error: %s", item.As, err.Error())
		}

		header := util.Convert2Map(item.Header)

		newHeader, err := mapbuilder.Build(src, header)
		if err != nil && item.Key {
			return list, fmt.Errorf("build `%s` header error: %s", item.As, err.Error())
		}
		list[i].Header = util.Map2MapString(newHeader)
	}
	return list, nil
}
