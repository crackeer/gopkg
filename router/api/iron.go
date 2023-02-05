package api

import "encoding/json"

// RequestItem ...
type RequestItem struct {
	API       string                 `json:"api"`
	Params    map[string]interface{} `json:"params"`
	Header    map[string]string      `json:"header"`
	As        string                 `json:"as"`
	ErrorExit bool                   `json:"error_exit"`
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
