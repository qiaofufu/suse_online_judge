package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"onlineJudge/common"
	"onlineJudge/models"
	"onlineJudge/response"
	"onlineJudge/util"
	"time"
)

// isEmailExit 邮箱是否存在
func isEmailExit(email string) bool {
	var user models.User
	result := common.DB.Model(&models.User{}).Where("user_email = ?", email).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return false
	} else {
		return true
	}

}

// SendCode 发送邮件
func SendCode(ctx *gin.Context) {
	type Email struct {
		Email string `json:"email"`
	}
	email := Email{}
	err := ctx.ShouldBindJSON(&email)
	if err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	if !util.CheckEmailMatch(email.Email) {
		response.CheckFail(ctx, nil, "邮箱格式错误")
		return
	}
	code, _ := common.RedisClient.Get(util.CodeKey(email.Email)).Result()
	if code != "" {
		duration, _ := common.RedisClient.TTL(util.CodeKey(email.Email)).Result()
		if duration >= time.Minute {
			response.Fail(ctx, nil, "操作过于频繁")
			return
		}
	}
	randomCode := util.RandomCode(6)
	common.RedisClient.Set(util.CodeKey(email.Email), randomCode, time.Minute * 10)
	send := util.SendEmail(email.Email, randomCode)
	if send {
		response.Success(ctx, nil, "发送成功")
	} else {
		common.RedisClient.Del(util.CodeKey(email.Email))
		response.Fail(ctx, nil, "发送失败")
	}
}

