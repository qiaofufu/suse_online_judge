package controller

import (
	"github.com/gin-gonic/gin"
	"onlineJudge/common"
	"onlineJudge/models"
	"onlineJudge/response"
)

// GetUserIntegral 获取当前用户积分
func GetUserIntegral(ctx *gin.Context) {
	user, err := getUserByJWT(ctx)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}
	// 操作
	var userIntegral models.UserIntegral
	result := common.DB.Model(&models.UserIntegral{}).Where("user_id = ?", user.ID).First(&userIntegral)
	if result.Error != nil {
		response.Fail(ctx, nil, "获取积分失败")
	} else {
		type responseDTO struct {
			Value uint `json:"value"`
			ConsumptionValue uint `json:"consumption_value"`
		}
		ret := responseDTO{Value: userIntegral.Value, ConsumptionValue: userIntegral.ConsumptionValue}
		response.Success(ctx, ret, "获取积分成功")
	}
}
