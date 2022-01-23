package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"onlineJudge/common"
	"onlineJudge/models"
	"onlineJudge/response"
	"time"
)

// GetsAnnouncementList 获取公告列表
func GetsAnnouncementList(ctx *gin.Context) {
	// 数据库操作
	var announcementList []models.Announcement
	result := common.DB.Model(&models.Announcement{}).Find(&announcementList)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "获取信息错误")
	} else {
		type responseDTO struct {
			Id uint `json:"id"`
			Title string `json:"title"`
			CreateAt time.Time `json:"create_at"`
		}
		var resDTO []responseDTO
		for _, value := range announcementList {
			var admin models.Admin
			if err := common.DB.Model(&models.Admin{}).Unscoped().Where("id = ?", value.AdminID).First(&admin).Error; err != nil {
				response.Fail(ctx, nil, "获取信息错误2")
				return
			}
			elem := responseDTO{Id: value.ID, Title: value.Title, CreateAt: value.CreatedAt}
			resDTO = append(resDTO, elem)
		}
		response.Success(ctx, resDTO, "获取信息成功")
	}
}

// GetAnnouncement 获取公告
func GetAnnouncement(ctx *gin.Context) {
	// 数据绑定
	type requiredDTO struct {
		ID uint `json:"id"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 处理
	var announcement models.Announcement
	result := common.DB.Model(&models.Announcement{}).Where("id = ?", dto.ID).First(&announcement)
	if result.Error == gorm.ErrRecordNotFound {
		response.Fail(ctx, nil, "未找到公告")
	} else if result.Error != nil {
		response.Fail(ctx, nil, "内部错误")
	} else {
		response.Success(ctx, announcement, "获取成功")
	}
}