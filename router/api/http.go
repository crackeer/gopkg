package api

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/crackeer/gopkg/util"
	"github.com/gin-gonic/gin/binding"
	"github.com/tidwall/gjson"
)

// NewAPIRequest
//
//	@param apiMeta
//	@param logger
//	@return *APIRequest
func NewAPIRequest(apiMeta *APIMeta, logger Logger) *APIRequest {
	return &APIRequest{
		APIMeta: apiMeta,
		Logger:  logger,
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

		client := &http.Client{Timeout: time.Millisecond * time.Duration(apiRequest.Timeout)}
		response, err = client.Do(request)

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
		fullURL = fmt.Sprintf("%s?%s", fullURL, urlParameter.Encode())
	}

	if contentType == "application/json" {
		byteData, _ := util.Marshal(parameter)
		body = strings.NewReader(string(byteData))
	} else if contentType == "application/json" {
		body = strings.NewReader(urlParameter.Encode())
	} else if contentType == binding.MIMEMultipartPOSTForm {
		body, contentType = packMultipart(parameter)
	}

	return fullURL, contentType, body
}

func packMultipart(data map[string]interface{}) (*bytes.Buffer, string) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for k, v := range data {
		if list, ok := v.([]string); ok {
			for _, item := range list {
				writer.WriteField(k, item)
			}
			continue
		}

		if fileList, ok := v.([]*multipart.FileHeader); ok {
			for _, f := range fileList {
				formFile, _ := writer.CreateFormFile(k, f.Filename)
				if tmp, err := f.Open(); err == nil {
					io.Copy(formFile, tmp)
				}
			}
			continue
		}
	}
	err := writer.Close()
	if err != nil {
		return nil, ""
	}
	return payload, writer.FormDataContentType()
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
