package api

type Signature interface {
	Sign(*APIMeta, map[string]interface{}, map[string]string) (map[string]interface{}, map[string]string)
}
