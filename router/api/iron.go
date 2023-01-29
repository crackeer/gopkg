package api

import "encoding/json"

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
