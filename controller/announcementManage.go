package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"onlineJudge/common"
	"onlineJudge/models"
	"onlineJudge/response"
	"strconv"
	"time"
)

/*
*	公告管理功能
*/

// AddAnnouncement 创建公告
func AddAnnouncement(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_ANNOUNCEMENTMANAGER) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		Title string `form:"title"`
		Content string `form:"content"`
	}
	var dto requiredDTO
	dto.Title = ctx.PostForm("title")
	dto.Content = ctx.PostForm("content")
	admin, err := getAdminJwt(ctx)
	if err != nil {
		response.Fail(ctx, nil, "管理员角色获取失败")
		return
	}
	picture, err := savePicture(ctx,"picture", "announcement-" + dto.Title +"-" + strconv.FormatUint(uint64(admin.ID), 10), "picture")
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}
	// 数据库操作
	announcement := models.Announcement{Title: dto.Title, Content: dto.Content, PictureUrl: picture, AdminID: admin.ID}
	result := common.DB.Model(&models.Announcement{}).Create(&announcement)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "创建公告失败")
	} else {
		response.Success(ctx, nil, "创建公告成功")
	}
}

// DeleteAnnouncement 删除公告
func DeleteAnnouncement(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_ANNOUNCEMENTMANAGER) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		DeleteId uint `json:"delete_id"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	result := common.DB.Where("id = ?", dto.DeleteId).Delete(&models.Announcement{})
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "删除失败")
	} else {
		response.Success(ctx, nil, "删除成功")
	}
}

// UpdateAnnouncement 更新公告信息
func UpdateAnnouncement(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_ANNOUNCEMENTMANAGER) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		AlterId uint `json:"alter_id"`
		Title string `json:"title"`
		Content string `json:"content"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	fmt.Println(dto)
	// 数据库操作
	result := common.DB.Model(&models.Announcement{}).Where("id = ?", dto.AlterId).Updates(models.Announcement{Title: dto.Title, Content: dto.Content})
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "数据更新失败")
	} else {
		response.Success(ctx, nil, "数据更新成功")
	}
}

// UpdateAnnouncementPicture 更新公告图片
func UpdateAnnouncementPicture (ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_ANNOUNCEMENTMANAGER) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		Title string `json:"title"`
		AlterID uint `json:"alter_id"`
	}
	var dto requiredDTO
	id, err := strconv.Atoi(ctx.PostForm("alter_id"))

	if err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}

	dto.AlterID = uint(id)
	dto.Title = ctx.PostForm("title")
	picture, err := savePicture(ctx,"picture", "announcement-" + dto.Title +"-" + strconv.FormatUint(uint64(dto.AlterID), 10), "picture")
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}
	// 数据库操作
	result := common.DB.Model(&models.Announcement{}).Where("id = ?", dto.AlterID).Update("picture_url", picture)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "更新失败")
	} else {
		response.Success(ctx, nil, "更新成功")
	}
}

// AdminGetsAnnouncementList 获取公告列表-后台
func AdminGetsAnnouncementList(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_ANNOUNCEMENTMANAGER) {
		return
	}
	// 数据库操作
	var announcementList []models.Announcement
	result := common.DB.Model(&models.Announcement{}).Find(&announcementList)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "获取信息错误")
	} else {
		type responseDTO struct {
			Id uint `json:"id"`
			Content string `json:"content"`
			Title string `json:"title"`
			CreateAt time.Time `json:"create_at"`
			Publisher string `json:"publisher"`
			PictureUrl string `json:"picture_url"`
		}
		var resDTO []responseDTO
		for _, value := range announcementList {
			var admin models.Admin
			if err := common.DB.Model(&models.Admin{}).Unscoped().Where("id = ?", value.AdminID).First(&admin).Error; err != nil {
				response.Fail(ctx, nil, "获取信息错误2")
				return
			}
			elem := responseDTO{Id: value.ID, Content: value.Content, Title: value.Title, CreateAt: value.CreatedAt, PictureUrl: value.PictureUrl, Publisher: admin.AdminName}
			resDTO = append(resDTO, elem)
		}
		response.Success(ctx, resDTO, "获取信息成功")
	}
}

