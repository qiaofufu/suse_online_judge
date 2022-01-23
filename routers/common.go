package routers

import (
	"github.com/gin-gonic/gin"
	"onlineJudge/controller"
	"onlineJudge/middlewares"
)

func MountCommonRouter(group *gin.RouterGroup) {
	userRouter := group.Group("/user")
	{
		userRouter.POST("/register", controller.Register)
		userRouter.POST("/register/code", controller.SendCode)
		userRouter.POST("/login", controller.Login)
		userRouter.POST("/get_user_name", controller.GetUserName)
		userRouter.GET("/rank", controller.Rank)
	}
	userRouterAuth := group.Group("/user")
	{
		userRouterAuth.Use(middlewares.AuthRequired())
		userRouterAuth.GET("/get/info", controller.GetUserInfo)
		userRouterAuth.POST("/update/password", controller.UpdateUserPassword)
		userRouterAuth.GET("/get/integral", controller.GetUserIntegral)
	}
	announcementRouter := group.Group("/announcement")
	{
		announcementRouter.GET("/gets/list", controller.GetsAnnouncementList)
		announcementRouter.POST("/get", controller.GetAnnouncement)
	}
	slideshowRouter := group.Group("/slideshow")
	{
		slideshowRouter.GET("/gets/list", controller.GetsSlideshowList)
	}
	contestRouter := group.Group("/contest")
	{
		contestRouter.POST("/list", controller.GetsContestList)
	}
	contestRouterAuth := group.Group("/contest")
	{
		contestRouterAuth.POST("/get", controller.GetContest)
		contestRouterAuth.POST("/problem", controller.GetContestProblemList)
		contestRouterAuth.POST("/submission", controller.GetsContestSubmission)
		contestRouterAuth.POST("/rank", controller.GetContestRank)

	}
	goodsRouter := group.Group("/goods")
	{
		goodsRouter.GET("/gets/typelist", controller.GetsGoodsTypeList)
		goodsRouter.POST("/gets/goodslist", controller.GetsGoodsList)
	}
	goodsRouterAuth := group.Group("/goods")
	{
		goodsRouterAuth.Use(middlewares.AuthRequired())
		goodsRouterAuth.POST("/buy", controller.BuyGoods)
		goodsRouterAuth.POST("/order", controller.GetsUserOrderList)
	}
	problemRouter := group.Group("/problem")
	{
		problemRouter.POST("/list", controller.GetsProblemList)
		problemRouter.POST("/content", controller.GetProblemContent)
	}
	problemRouterAuth := group.Group("/problem")
	{
		problemRouterAuth.Use(middlewares.AuthRequired())
		problemRouterAuth.POST("/contest", controller.GetContestProblemContent)
		problemRouterAuth.POST("/submission", controller.SubmissionCode)
		problemRouterAuth.POST("/exists", controller.SubmissionExists)
	}
	submissionRouter := group.Group("/submission")
	{
		submissionRouter.POST("/list", controller.SubmissionList)
	}
	submissionRouterAuth := group.Group("/submission")
	{
		submissionRouterAuth.Use(middlewares.AuthRequired())

	}
	commonRouter := group.Group("/common")
	{
		commonRouter.POST("/get/admin/name", controller.GetAdminName)
	}
}
