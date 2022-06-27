package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

const YamlExt = ".yaml"

// LoadYamlConfig
//  @param path
//  @param dest
//  @return error
func LoadYamlConfig(path string, dest interface{}) error {
	if !strings.HasSuffix(path, YamlExt) {
		path = fmt.Sprintf("%s%s", path, YamlExt)
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config `%s` error: %s", path, err.Error())
	}

	return yaml.Unmarshal(bytes, dest)

}

// LoadYamlAPIConfig
//  @param path
//  @return *YamlAPI
//  @return error
func LoadYamlAPIConfig(path string) (*YamlAPI, error) {
	retData := &YamlAPI{}
	err := LoadYamlConfig(path, retData)
	return retData, err
}

// LoadYamlPageConfig
//  @param path
//  @return *YamlPage
//  @return error
func LoadYamlGroupPageConfig(path string) (*YamlGroupPage, error) {
	retData := &YamlGroupPage{}
	err := LoadYamlConfig(path, retData)
	return retData, err
}
