package api

import (
	"fmt"
	"sync"

	"github.com/crackeer/gopkg/mapbuilder"
	"github.com/crackeer/gopkg/util"
)

const (

	// InputKey ...
	InputKey = "_input_"

	// HeaderKey
	HeaderKey = "_header_"
)

// RequestClient ...
type RequestClient struct {
	factory APIFactory
	env     string
	logger  Logger

	meshAPISeperator string
	meshAPIPrefix    string
}

// NewRequestClient
//
//	@param getter
//	@return *RequestClient
func NewRequestClient(getter APIFactory) *RequestClient {
	return &RequestClient{
		factory:          getter,
		meshAPISeperator: mapbuilder.DefaultSeperator,
		meshAPIPrefix:    mapbuilder.DefaultPrefix,
	}
}

func (client *RequestClient) UseEnv(env string) {
	client.env = env
}

func (client *RequestClient) SetMeshAPIConfig(prefix, seq string) {
	client.meshAPISeperator = seq
	client.meshAPIPrefix = prefix
}

// Request
//
//	@receiver client
//	@param name
//	@param query
//	@param header
//	@return *APIResponse
//	@return error
func (client *RequestClient) Request(apiID string, query map[string]interface{}, header map[string]string) (*APIResponse, error) {
	apiMeta := client.factory.Get(apiID, client.env)
	if apiMeta == nil {
		return nil, fmt.Errorf("api `%s` not found", apiID)
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
func (client *RequestClient) requestList(list []*RequestItem) (map[string]*APIResponse, map[string]string, error) {

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
			apiResponse, err := client.Request(tmp.API, tmp.Params, tmp.Header)
			locker.Lock()
			defer locker.Unlock()
			if err != nil {
				errMap[tmp.As] = err.Error()
				if tmp.ErrorExit {
					retError = fmt.Errorf("request api `%s` response error:%s", tmp.As, err.Error())
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
func (client *RequestClient) Mesh(list [][]*RequestItem, query map[string]interface{}, header map[string]string) (map[string]*APIResponse, map[string]string, error) {

	var (
		retError error
		errMap   map[string]string       = make(map[string]string)
		retData  map[string]*APIResponse = make(map[string]*APIResponse)
	)

	input := map[string]interface{}{
		InputKey:  query,
		HeaderKey: header,
	}

	for _, items := range list {
		newItems, err := transfer(items, input, client.meshAPIPrefix, client.meshAPISeperator)
		if err != nil {
			retError = err
			break
		}
		mapAPIResponse, mapError, err := client.requestList(newItems)
		if err != nil {
			retError = err
			break
		}

		for key, response := range mapAPIResponse {
			retData[key] = response
			input[key] = response.Data
		}

		for key, message := range mapError {
			errMap[key] = message
		}

	}
	return retData, errMap, retError
}

func transfer(list []*RequestItem, input map[string]interface{}, prefix, seq string) ([]*RequestItem, error) {
	src := util.ToString(input)
	for i, item := range list {

		if len(item.Params) > 0 {
			if value := mapbuilder.Build([]byte(src), item.Params, prefix, seq); value != nil {
				list[i].Params = value.(map[string]interface{})
			}
		}

		if len(item.Header) > 0 {
			header := util.Convert2Map(item.Header)
			if newHeader := mapbuilder.Build([]byte(src), header, prefix, seq); newHeader != nil {
				list[i].Header = util.Map2MapString(newHeader.(map[string]interface{}))
			}

		} else {
			header := util.LoadMap(input).GetMap(HeaderKey)
			list[i].Header = util.Map2MapString(header)
		}
	}
	return list, nil
}
