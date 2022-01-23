package controller

import (
	"github.com/gin-gonic/gin"
	"onlineJudge/common"
	"onlineJudge/models"
	"onlineJudge/response"
	"strconv"
)

// AddSlideshow 添加轮播图
func AddSlideshow(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_SLIDESHOWMANAGE) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		Title string `json:"title"`
	}
	var dto requiredDTO
	dto.Title = ctx.PostForm("title")
	picture, err := savePicture(ctx, "picture", dto.Title, "slideshow")
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}
	// 数据库操作
	elem := models.Slideshow{Title: dto.Title, Url: picture}
	result := common.DB.Model(&models.Slideshow{}).Create(&elem)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "添加失败")
	} else {
		response.Success(ctx, nil, "添加成功")
	}
}

// UpdateSlideshow 更新轮播图
func UpdateSlideshow(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_SLIDESHOWMANAGE) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		AlterID uint `json:"alter_id"`
		Title string `json:"title"`
	}
	var dto requiredDTO
	id, _ := strconv.Atoi(ctx.PostForm("alter_id"))
	dto.AlterID = uint(id)
	dto.Title = ctx.PostForm("title")
	picture, err := savePicture(ctx, "picture", dto.Title, "slideshow")
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}
	// 数据库操作
	elem := models.Slideshow{Title: dto.Title, Url: picture}
	result := common.DB.Where("id = ?", dto.AlterID).Updates(elem)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "更新失败")
	} else {
		response.Success(ctx, nil, "更新成功")
	}
}

// DeleteSlideshow 删除轮播图
func DeleteSlideshow(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_SLIDESHOWMANAGE) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		DeleteID uint `json:"delete_id"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	result := common.DB.Where("id = ?", dto.DeleteID).Delete(&models.Slideshow{})
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "删除失败")
	} else {
		response.Success(ctx, nil, "删除成功")
	}
}

// AdminGetsSlideshowList 获取轮播图列表
func AdminGetsSlideshowList(ctx *gin.Context) {
	// 数据库操作
	var slideshowList []models.Slideshow
	result := common.DB.Model(&models.Slideshow{}).Order("id desc").Find(&slideshowList)
	if result.Error != nil {
		response.Fail(ctx, nil, "获取数据错误")
	} else {
		response.Success(ctx, slideshowList, "获取成功")
	}
}