package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(ctx *gin.Context, httpStatus int, code int, data interface{}, msg string) {
	ctx.JSONP(httpStatus, gin.H{
		"code" : code,
		"data" : data,
		"msg" : msg,
	})
}

func Fail(ctx *gin.Context, data interface{}, msg string) {
	Response(ctx, http.StatusOK, -1, data, msg)
}

func Success(ctx *gin.Context, data interface{}, msg string) {
	Response(ctx, http.StatusOK, 0, data, msg)
}

func CheckFail(ctx *gin.Context, data interface{}, msg string) {
	Response(ctx, http.StatusOK, -2, data, msg)
}
