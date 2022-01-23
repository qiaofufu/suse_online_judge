package controller

import (
	"github.com/gin-gonic/gin"
	"onlineJudge/common"
	"onlineJudge/models"
	"onlineJudge/response"
)

// GetsSlideshowList 获取轮播图列表
func GetsSlideshowList(ctx *gin.Context) {
	// 数据库操作
	var slideshowList []models.Slideshow
	result := common.DB.Model(&models.Slideshow{}).Order("id desc").Find(&slideshowList)
	if result.Error != nil {
		response.Fail(ctx, nil, "获取数据错误")
	} else {
		response.Success(ctx, slideshowList, "获取成功")
	}
}