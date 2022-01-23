package routers

import (
	"github.com/gin-gonic/gin"
	"onlineJudge/controller"
	"onlineJudge/middlewares"
)

func MountAdminRouter(group *gin.RouterGroup) {
	// 路由分组
	router := group.Group("/admin")
	routerAuth := group.Group("/admin")
	routerAuth.Use(middlewares.AuthRequired())

	// 分发
	{
		// 无需权限管理员
		router.POST("/login", controller.AdminLogin)
	}
	{
		// 管理员操作
		routerAuth.POST("/info/update",controller.AdminInfoUpdate)
		routerAuth.GET("/get/info", controller.GetAdminInfo)
		// 管理员管理功能
		routerAuth.POST("/manage/admin/add", controller.AddAdmin)
		routerAuth.POST("/manage/admin/update", controller.UpdateAdmin)
		routerAuth.DELETE("/manage/admin/delete", controller.DeleteAdmin)
		routerAuth.GET("/manage/admin/gets/adminlist", controller.GetsAdminsList)
		routerAuth.POST("manage/admin/role/add", controller.AddRole)
		routerAuth.POST("/manage/admin/role/update", controller.UpdateRole)
		routerAuth.DELETE("/manage/admin/role/delete", controller.DeleteRole)
		routerAuth.GET("/manage/admin/role/rolelist", controller.GetsRolesList)
		// 用户管理功能
		routerAuth.POST("/manage/user/add", controller.AddUser)
		routerAuth.POST("/manage/user/update", controller.UpdateUser)
		routerAuth.DELETE("/manage/user/delete", controller.DeleteUser)
		routerAuth.GET("/manage/user/gets/userlist", controller.GetsUsersList)
		// 公告管理
		routerAuth.POST("/manage/announcement/add", controller.AddAnnouncement)
		routerAuth.DELETE("/manage/announcement/delete", controller.DeleteAnnouncement)
		routerAuth.POST("/manage/announcement/update", controller.UpdateAnnouncement)
		routerAuth.POST("/manage/announcement/picture/update", controller.UpdateAnnouncementPicture)
		routerAuth.GET("/manage/announcement/gets/announcementlist", controller.AdminGetsAnnouncementList)
		// 积分管理
		routerAuth.POST("/manage/user/integral/set", controller.SetUserIntegral)
		routerAuth.GET("/manage/user/integral/integrallist", controller.GetsUserIntegralList)
		// 货物管理
		routerAuth.POST("/manage/goods/type/add", controller.AddGoodsType)
		routerAuth.POST("/manage/goods/type/update", controller.UpdateGoodsType)
		routerAuth.DELETE("/manage/goods/type/delete", controller.DeleteGoodsType)
		routerAuth.GET("/manage/goods/type/gets/goodstypelist", controller.AdminGetsGoodsTypeList)
		routerAuth.POST("/manage/goods/add", controller.AddGoods)
		routerAuth.POST("/manage/goods/update", controller.UpdateGoods)
		routerAuth.POST("/manage/goods/picture/update", controller.UpdateGoodsPicture)
		routerAuth.DELETE("/manage/goods/delete", controller.DeleteGoods)
		routerAuth.GET("/manage/goods/gets/goodslist", controller.AdminGetsGoodsList)
		routerAuth.POST("/manage/order/gets/orderlist", controller.GetsOrderList)
		routerAuth.POST("/manage/order/finish", controller.FinishOrder)
		// 竞赛管理
		routerAuth.POST("/manage/contest/add", controller.AddContest)
		routerAuth.GET("/manage/contest/gets/contestlist", controller.AdminGetsContestList)
		routerAuth.POST("/manage/contest/update", controller.UpdateContest)
		routerAuth.DELETE("/manage/contest/delete", controller.DeleteContest)
		routerAuth.POST("/manage/problem/add", controller.ADDContestProblem)
		routerAuth.DELETE("/manage/problem/delete", controller.DeleteContestProblem)
		routerAuth.POST("/manage/problem/update", controller.UpdateContestProblem)
		routerAuth.POST("/manage/problem/gets/problemlist", controller.GetsContestProblemList)
		routerAuth.POST("/manage/example/add", controller.AddExample)
		routerAuth.DELETE("/manage/example/delete", controller.DeleteExample)
		routerAuth.POST("/manage/example/update", controller.UpdateExample)
		routerAuth.POST("/manage/example/gets/examplelist", controller.GetsExampleList)
		routerAuth.POST("/manage/contest/answer/release", controller.ReleaseAnswer)
		// 轮播图管理
		routerAuth.POST("/manage/slideshow/add", controller.AddSlideshow)
		routerAuth.POST("/manage/slideshow/update", controller.UpdateSlideshow)
		routerAuth.DELETE("/manage/slideshow/delete", controller.DeleteSlideshow)
		routerAuth.GET("/manage/slideshow/gets/slideshowlist", controller.AdminGetsSlideshowList)
	}
}
