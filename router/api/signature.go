package api

// SignConfig
type SignConfig struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
}
type Signature interface {
	Sign(*APIMeta, map[string]interface{}, map[string]string) (map[string]interface{}, map[string]string)
}
