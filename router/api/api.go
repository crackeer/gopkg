package api

// APIMeta API
type APIMeta struct {
	Host        string            `json:"host"`
	Path        string            `json:"path"`
	ContentType string            `json:"content_type"`
	Method      string            `json:"method"`
	SuccessCode string            `json:"success_code"`
	MessageKey  string            `json:"message_key"`
	CodeKey     string            `json:"code_key"`
	DataKey     string            `json:"data_key"`
	Timeout     int64             `json:"timeout"`
	SignConfig  *SignConfig       `json:"sign_config"`
	Header      map[string]string `json:"header"`
	CacheTime   int64             `json:"cache_time"`
}

// APIResponse
type APIResponse struct {
	Code           string
	Message        string
	Data           string
	HttpStatusCode int64
	OriginBody     []byte
}

// APIRequest
type APIRequest struct {
	*APIMeta
	Logger
}

// APIMetaFactory ...
type APIMetaFactory interface {
	GetAPIMeta(string, string) (*APIMeta, error)
}

type Logger interface {
	Error(map[string]interface{}, string)
	Info(map[string]interface{}, string)
}
