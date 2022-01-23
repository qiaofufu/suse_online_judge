package middlewares

import (

	"github.com/gin-gonic/gin"
	"net/http"
	"onlineJudge/common"
)

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == ""{
			ctx.JSON(http.StatusUnauthorized, gin.H{"code" : 401, "msg" : "请先登录"})
			ctx.Abort()
			return
		}
		if !common.VerifyToken(tokenString) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code" : 401, "msg" : "请先登录"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}


