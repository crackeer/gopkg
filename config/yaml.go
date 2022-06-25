package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadYamlConfig
//  @param path
//  @param dest
//  @return error
func LoadYamlConfig(path string, dest interface{}) error {
	if !strings.HasSuffix(path, ".yaml") {
		path = fmt.Sprintf("%s.yaml", path)
	}
	bytes, err := ioutil.ReadFile(path)
	fmt.Println(string(bytes))
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
