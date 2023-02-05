package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/crackeer/gopkg/util"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

// APIRequest
type APIRequest struct {
	*APIMeta
	LogEntry *LogEntry
}

// NewAPIRequest
//
//	@param apiMeta
//	@param logger
//	@return *APIRequest
func NewAPIRequest(apiMeta *APIMeta) *APIRequest {
	return &APIRequest{
		APIMeta: apiMeta,
	}
}

// AddHeader
//
//	@receiver apiRequest
//	@param header
func (apiRequest *APIRequest) AddHeader(header map[string]string) {
	if apiRequest.APIMeta.Header == nil {
		apiRequest.APIMeta.Header = map[string]string{}
	}
	for key, value := range header {
		apiRequest.APIMeta.Header[key] = value
	}

}

func (apiRequest *APIRequest) Do(parameter map[string]interface{}, header map[string]string) (*APIResponse, error) {
	apiRequest.AddHeader(header)
	fullURL, contentType, requestBody := apiRequest.getPrimary(parameter)
	apiRequest.AddHeader(map[string]string{
		"Content-Type": contentType,
	})

	var (
		request     *http.Request
		response    *http.Response
		retError    error
		err         error
		apiResponse = &APIResponse{
			Name: apiRequest.APIMeta.Name,
		}
	)

	for {

		request, err = http.NewRequest(apiRequest.Method, fullURL, requestBody)
		if err != nil {
			retError = fmt.Errorf("new request error: %s", err.Error())
			break
		}

		// Build header
		for key, value := range apiRequest.Header {
			request.Header.Set(key, value)
		}
		apiRequest.LogEntry = NewLogEntryFromRequest(request)
		apiRequest.LogEntry.SetStart()
		client := &http.Client{Timeout: time.Millisecond * time.Duration(apiRequest.Timeout)}
		response, err = client.Do(request)
		apiRequest.LogEntry.SetEnd()
		if err != nil {
			retError = fmt.Errorf("request error: %s", err.Error())
			break
		}

		var byteBody []byte

		byteBody, err = ioutil.ReadAll(response.Body)

		if err != nil {
			retError = fmt.Errorf("read request body error: %s", err.Error())
			break
		}
		apiRequest.LogEntry.AddRespone(byteBody, int64(response.StatusCode))

		// Build response
		apiResponse.HttpStatusCode = int64(response.StatusCode)
		apiResponse.OriginBody = byteBody
		if response.StatusCode != http.StatusOK {
			retError = fmt.Errorf("http_error %d, body=%s", response.StatusCode, string(apiResponse.OriginBody))
			break
		}

		if len(apiRequest.CodeKey) > 0 {
			apiResponse.Code = gjson.GetBytes(byteBody, apiRequest.CodeKey).String()
		}

		if len(apiRequest.MessageKey) > 0 {
			apiResponse.Message = gjson.GetBytes(byteBody, apiRequest.MessageKey).String()
		}

		if len(apiRequest.DataKey) < 1 {
			apiResponse.Data = util.TryJSONParse(string(byteBody))
		} else {
			apiResponse.Data = gjson.GetBytes(byteBody, apiRequest.DataKey).Value()
		}

		if len(apiRequest.SuccessCode) > 0 && len(apiResponse.Code) > 0 {
			if apiResponse.Code != apiRequest.SuccessCode {
				retError = fmt.Errorf("response error: %s", apiResponse.Message)
				break
			}
		}
		break
	}
	if retError != nil {
		apiRequest.LogEntry.AddError(retError.Error())
	}

	return apiResponse, retError
}

func (apiRequest *APIRequest) getPrimary(parameter map[string]interface{}) (string, string, io.Reader) {
	fullURL := fmt.Sprintf("%s/%s", apiRequest.Host, apiRequest.Path)
	contentType := apiRequest.getContentType()

	var body io.Reader

	urlParameter := &url.Values{}
	for key, value := range parameter {
		strVal := util.ToString(value)
		urlParameter.Set(key, strVal)
	}
	if apiRequest.Method == http.MethodGet {
		if value := urlParameter.Encode(); len(value) > 0 {
			fullURL = fmt.Sprintf("%s?%s", fullURL, urlParameter.Encode())
		}
	}

	if contentType == gin.MIMEJSON {
		byteData, _ := util.Marshal(parameter)
		body = strings.NewReader(string(byteData))
	} else if contentType == gin.MIMEPOSTForm {
		body = strings.NewReader(urlParameter.Encode())
	}

	return fullURL, contentType, body
}

func (apiRequest *APIRequest) getContentType() string {
	if value, ok := apiRequest.Header["Content-Type"]; ok {
		return value
	}

	if len(apiRequest.ContentType) > 0 {
		return apiRequest.ContentType
	}

	return ""
}

// GetLog GetLogMap
//
//	@receiver apiRequest
//	@return map
func (apiRequest *APIRequest) GetLog() map[string]interface{} {
	if apiRequest.LogEntry != nil {
		return apiRequest.LogEntry.Map()
	}
	return nil
}
