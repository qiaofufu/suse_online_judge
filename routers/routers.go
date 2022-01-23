package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func MountRouter(engine *gin.Engine) {
	router := engine.Group("/onlinejudge")
	router.StaticFS("/file", http.Dir("./file"))
	MountAdminRouter(router)
	MountCommonRouter(router)
}
