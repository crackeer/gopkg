package render

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/crackeer/gopkg/util"
)

// RenderHTML
//  @param framePath
//  @param contentPath
//  @param injectData
//  @return string
//  @return error
func RenderHTML(framePath string, contentPath string, opt *Option) (string, error) {

	var (
		content string
	)
	if opt == nil {
		opt = DefaultOption()
	}
	if len(contentPath) > 0 {
		bytes, err := ioutil.ReadFile(contentPath)
		if err != nil {
			return "", errors.New("content file not exists")
		}
		content = string(bytes)
	}

	var injectData string
	if opt.InjectData != nil {
		injectData, _ = util.MarshalEscapeHtml(opt.InjectData)
	}

	if len(framePath) < 1 {
		return strings.Replace(content, opt.PlaceholderJSData, injectData, -1), nil
	}
	bytes, err := ioutil.ReadFile(framePath)
	if err != nil {
		return "", errors.New("frame file not exists")
	}
	fulltext := strings.Replace(string(bytes), opt.PlaceholderContent, content, -1)

	return strings.Replace(fulltext, opt.PlaceholderJSData, injectData, -1), nil

}
