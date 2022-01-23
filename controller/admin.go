package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"onlineJudge/common"
	"onlineJudge/models"
	"onlineJudge/response"
	"strings"
)

/*
*	管理员基本功能
*/

// AdminLogin 管理员登录
func AdminLogin(ctx *gin.Context) {
	// 数据绑定
	type requiredDTO struct {
		AdminAccount string `json:"admin_account"`
		AdminPassword string `json:"admin_password"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	var admin models.Admin
	result := common.DB.Model(&models.Admin{}).Where("admin_account = ? and admin_password = ?", dto.AdminAccount, dto.AdminPassword).First(&admin)
	// 响应
	if result.Error == gorm.ErrRecordNotFound {
		response.Fail(ctx, nil, "账号或密码错误")
	} else if result.Error != nil {
		response.Fail(ctx, nil, "内部错误")
	} else {
		tokenString, err := common.CreateToken(map[string]interface{}{"admin_id":admin.ID})
		if err != nil {
			response.Fail(ctx, nil, "token获取失败")
		} else {
			type responseDTO struct {
				ID uint `json:"id"`
				AdminName string `json:"admin_name"`
				AdminAccount string `json:"admin_account"`
				AdminRoles	string `json:"admin_roles"`
				Token string `json:"token"`
			}
			response.Success(ctx, responseDTO{ID: admin.ID, AdminName: admin.AdminName, AdminAccount: admin.AdminAccount, AdminRoles: admin.AdminRoles, Token: tokenString}, "获取信息成功")
		}
	}
}

// AdminInfoUpdate 管理员信息更新
func AdminInfoUpdate(ctx *gin.Context) {
	// 数据绑定
	type requiredDTO struct {
		AdminName string `json:"admin_name"`
		AdminPassword string `json:"admin_password"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	admin, err := getAdminJwt(ctx)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}
	admin.AdminName = dto.AdminName
	admin.AdminPassword = dto.AdminPassword
	result := common.DB.Save(&admin)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "更新失败")
	} else {
		response.Success(ctx, admin, "更新成功")
	}
}

// GetAdminInfo 获取管理员个人信息
func GetAdminInfo(ctx *gin.Context) {
	// 数据绑定
	admin, err := getAdminJwt(ctx)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}
	// 响应
	response.Success(ctx, admin, "获取成功")
}

// roleCheck 角色检查
func roleCheck(ctx *gin.Context, role string) bool {
	admin, err := getAdminJwt(ctx)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return false
	}
	if strings.Index(admin.AdminRoles, role) == -1 && strings.Index(admin.AdminRoles, "admin") == -1 {
		response.Fail(ctx, nil, "权限等级不够，非法访问")
		return false
	}
	return true
}

// getAdminJwt 根据请求头中的 Authorization 来获取admin
func getAdminJwt(ctx *gin.Context) (models.Admin, error) {
	tokenString := ctx.GetHeader("Authorization")
	token, err := common.ParseToken(tokenString)
	if err != nil {
		return models.Admin{}, errors.New("token认证失败，请重新登录")
	}
	adminID, ok := token.Get("admin_id")
	if !ok {
		return models.Admin{}, errors.New("获取信息错误")
	}
	var admin models.Admin
	result := common.DB.Model(&models.Admin{}).Where("id = ?", adminID).First(&admin)
	if result.Error != nil {
		return models.Admin{}, errors.New("未找到用户，请检查用户密码")
	}
	return admin, nil
}

// GetAdminName 获取管理员名字
func GetAdminName(ctx *gin.Context) {
	type requiredDTO struct {
		ID uint `json:"id"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	type responseDTO struct {
		AdminName string `json:"admin_name"`
	}
	var res responseDTO
	var admin models.Admin
	result := common.DB.Model(&models.Admin{}).Where("id = ?", dto.ID).First(&admin)
	res.AdminName = admin.AdminName
	if result.Error != nil {
		response.Fail(ctx, nil, "查询失败")
	} else {
		response.Success(ctx, res, "查询成功")
	}
}