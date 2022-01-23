package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"onlineJudge/common"
	"onlineJudge/models"
	"onlineJudge/response"
	"onlineJudge/util"
)

// Register 用户注册
func Register(ctx *gin.Context) {
	// 数据绑定
	type requiredDTO struct {
		UserName string `json:"user_name"`
		UserEmail string `json:"user_email"`
		UserPassword string `json:"user_password"`
		UserEmailCode string `json:"user_email_code"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil , "数据绑定失败")
		return
	}
	// 检验
	if !util.CheckPasswordMatch(dto.UserPassword) {
		response.Fail(ctx, nil, "密码格式错误")
		return
	}
	if !util.CheckEmailMatch(dto.UserEmail) {
		response.Fail(ctx, nil, "邮箱格式错误")
		return
	}
	if isEmailExit(dto.UserEmail) {
		response.Fail(ctx, nil, "邮箱已注册")
		return
	}
	if !util.VerifyEmailCode(dto.UserEmailCode, dto.UserEmail) {
		response.Fail(ctx, nil, "验证码不正确请重新获取")
		return
	}
	// 处理
	user := models.User{
		UserName: dto.UserName,
		UserEmail: dto.UserEmail,
		UserPassword: dto.UserPassword,
	}
	result := common.DB.Create(&user)
	if result.Error != nil {
		response.Fail(ctx, nil, "注册失败")
	} else {
		userIntegral := models.UserIntegral{UserID: user.ID}
		result := common.DB.Create(&userIntegral)
		if result.Error != nil {
			response.Fail(ctx, nil, "注册失败")
		} else {
			response.Success(ctx, nil, "注册成功")
		}
	}
}

// Login 用户登录
func Login(ctx *gin.Context) {
	// 数据绑定
	type requiredDTO struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 验证
	if !util.CheckEmailMatch(dto.Email) {
		response.Fail(ctx, nil, "邮箱格式错误")
		return
	}
	// 处理
	var user models.User
	result := common.DB.Model(&models.User{}).Where("user_email = ? and user_password = ?", dto.Email, dto.Password).First(&user)
	tokenStr, err := common.CreateToken(map[string]interface{}{"id":user.ID, "":user.UserEmail})
	if err != nil {
		response.Fail(ctx, nil, "token口令创建失败")
		return
	}
	if result.Error == gorm.ErrRecordNotFound {
		response.Fail(ctx, nil, "用户邮箱或密码错误")
	} else if result.Error != nil {
		response.Fail(ctx, nil, "内部错误")
	} else {
		type responseDTO struct {
			JWT string `json:"jwt"`
		}
		var ret responseDTO
		ret.JWT = tokenStr
		response.Success(ctx, ret, "登录成功")
	}
}

// GetUserInfo 获取用户个人信息
func GetUserInfo(ctx *gin.Context) {
	type problemList struct {
		ProblemID uint `json:"problem_id"`
		Status int `json:"status"`
	}
	type responseDTO struct {
		AcceptNumber int64 `json:"accept_number"`
		ProblemStatus []problemList `json:"problem_status"`
		User models.User `json:"user"`
		SubmissionNumber int64 `json:"submission_number"`
	}
	var res responseDTO
	user, err := getUserByJWT(ctx)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}
	user.UserPassword = ""
	res.User = user

	if err := common.DB.Model(&models.Submission{}).Where("user_id = ?", user.ID).Count(&res.SubmissionNumber).Error; err != nil {
		response.Fail(ctx, nil, "获取失败1")
		return
	}
	var p []models.ProblemStatusRecord
	if err := common.DB.Model(&models.ProblemStatusRecord{}).Where("user_id = ?", user.ID).Find(&p).Error; err != nil {
		response.Fail(ctx, nil, "获取失败1")
		return
	}
	for _, v := range p {
		elem := problemList{ProblemID: v.ProblemID, Status: v.Status}
		res.ProblemStatus = append(res.ProblemStatus, elem)
		if v.Status == 1 {
			res.AcceptNumber++
		}
	}
	response.Success(ctx, res, "获取成功")
}

// UpdateUserPassword 更新用户密码
func UpdateUserPassword(ctx *gin.Context) {
	// 数据绑定
	type requiredDTO struct {
		Password string `json:"password"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	user, err := getUserByJWT(ctx)
	if err != nil {
		response.Fail(ctx, nil, "获取用户失败")
		return
	}
	if !util.CheckPasswordMatch(dto.Password) {
		response.Fail(ctx, nil ,"密码格式不正确")
		return
	}
	// 处理
	user.UserPassword = dto.Password
	result := common.DB.Model(&models.User{}).Save(user)
	if result.Error != nil {
		response.Fail(ctx, nil, "更新失败")
	} else {
		response.Success(ctx, nil, "更新成功")
	}
}

// GetUserName 获取用户名
func GetUserName(ctx *gin.Context) {
	type requiredDTO struct {
		UserID uint `json:"user_id"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
	}
	var user models.User
	user, err := user.GetByID(dto.UserID)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}
	type responseDTO struct {
		UserName string `json:"user_name"`
	}
	var res responseDTO
	res.UserName = user.UserName
	response.Success(ctx,res,"获取成功")
}

func Rank(ctx *gin.Context) {
	rank, _ := models.User{}.Rank()
	response.Success(ctx, rank, "")
}

// getUserByJWT 从jwt获取用户
func getUserByJWT(ctx *gin.Context) (models.User, error) {
	tokenString := ctx.GetHeader("Authorization")
	token, err := common.ParseToken(tokenString)
	if err != nil {
		return models.User{}, errors.New("token口令错误")
	}
	id, ok := token.Get("id")
	if !ok {
		return models.User{}, errors.New("获取id失败")
	}
	var user models.User
	result := common.DB.Model(&models.User{}).Where("id = ?", id).First(&user)
	if result.Error != nil {
		return models.User{}, errors.New("权限认证失败-1")
	} else {
		return user, nil
	}
}