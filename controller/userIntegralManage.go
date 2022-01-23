package controller

import (
	"github.com/gin-gonic/gin"
	"onlineJudge/common"
	"onlineJudge/models"
	"onlineJudge/response"
)

/*
*	用户积分管理
*/

// SetUserIntegral 设置用户积分
func SetUserIntegral(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_USERMANAGER) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		UserID uint `json:"user_id"`
		Value uint `json:"value"`
		ConsumptionValue uint `json:"consumption_value"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	result := common.DB.Model(&models.UserIntegral{}).Where("user_id = ?", dto.UserID).Updates(models.UserIntegral{ConsumptionValue: dto.ConsumptionValue, Value: dto.Value})
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "设置错误")
	} else {
		response.Success(ctx, nil, "设置成功")
	}
}

// GetsUserIntegralList 获取用户积分列表
func GetsUserIntegralList(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_USERMANAGER) {
		return
	}
	// 数据库操作
	var userIntegralList []models.UserIntegral
	result := common.DB.Model(&models.UserIntegral{}).Find(&userIntegralList)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "数据获取失败")
	} else {
		response.Success(ctx, userIntegralList, "获取用户积分列表成功")
	}
}

