package gin

import (
	"github.com/gin-gonic/gin"
)

// DoResponseJSON
//  @return gin.HandlerFunc
func DoResponseJSON() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		ResponseJSON(ctx)
	}
}
