package api

// APIMeta API
type APIMeta struct {
	BaseURI     string
	Path        string
	ContentType string
	Method      string
	SuccessCode string
	MessageKey  string
	CodeKey     string
	DataKey     string
	Timeout     int64
	SignAlg     string
	Header      map[string]string
}

// APIResponse
type APIResponse struct {
	Code           string
	Message        string
	Data           interface{}
	HttpStatusCode int64
	OriginBody     []byte
}

// APIRequest
type APIRequest struct {
	*APIMeta
	Logger
}

// APIMetaGetter
type APIMetaGetter interface {
	GetAPIMeta(string, string) (*APIMeta, error)
}

type Logger interface {
	Error(map[string]interface{}, string)
	Info(map[string]interface{}, string)
}
