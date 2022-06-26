package ginhelper

// JSONRespone
type JSONRespone struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const keyResponse = "__RESPONSE__"

const (

	// CodeSuccess
	CodeSuccess int64 = 0

	CodeDefaultError int64 = -1
)

const (

	// MessageSuccess
	MessageSuccess = "success"

	// MessageEmptyResponse
	MessageEmptyResponse = "no router match or no data response"
)
