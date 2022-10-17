package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success  ...
//  @param ctx
//  @param data
func Success(ctx *gin.Context, data interface{}) {
	SetResponse(ctx, CodeSuccess, data, MessageSuccess)
}

// Failure ...
//  @param ctx
//  @param code
//  @param message
func Failure(ctx *gin.Context, code int64, message string) {
	SetResponse(ctx, code, nil, message)
}

// SetResponse ...
//  @param ctx
//  @param code
//  @param data
//  @param message
func SetResponse(ctx *gin.Context, code int64, data interface{}, message string) {
	if response, exists := getJSONResponse(ctx); exists {
		response.Code = code
		response.Message = message
		response.Data = data
	} else {
		ctx.Set(keyResponse, &JSONResponse{
			Code:    code,
			Message: message,
			Data:    data,
		})
	}
}

// getJSONResponse
//  @param ctx
//  @return *JSONResponse
//  @return bool
func getJSONResponse(ctx *gin.Context) (*JSONResponse, bool) {
	if body, exists := ctx.Get(keyResponse); exists {
		response, flag := body.(*JSONResponse)
		return response, flag
	}
	return nil, false
}

// ResponseJSON ...
//  @param ctx N
func ResponseJSON(ctx *gin.Context) {
	if ctx.IsAborted() {
		return
	}
	if response, found := getJSONResponse(ctx); found {
		ctx.Abort()
		ctx.PureJSON(http.StatusOK, response)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, &JSONResponse{
		Code:    CodeDefaultError,
		Message: MessageEmptyResponse,
		Data:    nil,
	})
}
