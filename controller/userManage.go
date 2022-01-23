package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"onlineJudge/common"
	"onlineJudge/models"
	"onlineJudge/response"
)

/*
*	用户管理功能
*/

// AddUser 添加用户
func AddUser(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_USERMANAGER) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		UserName 		string `json:"user_name"`
		UserEmail 		string `json:"user_email"`
		UserPassword 	string `json:"user_password"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	common.DB.Transaction(func(tx *gorm.DB) error {
		// 检测用户邮箱是否存在
		if flag, err := isUserExist("user_email", dto.UserEmail); flag && err == nil {
			response.Fail(ctx, nil, "用户邮箱已存在")
			return errors.New("邮箱以存在")
		} else if err != nil {
			response.Fail(ctx, nil, err.Error())
			return err
		}
		// 创建用户
		user := models.User{UserName: dto.UserName, UserPassword: dto.UserPassword, UserEmail: dto.UserEmail}
		result := tx.Create(&user)
		if result.Error != nil {
			response.Fail(ctx, nil, "创建用户失败")
			return result.Error
		} else {
			userIntegral := models.UserIntegral{UserID: user.ID}
			result := tx.Create(&userIntegral)
			if result.Error != nil {
				response.Fail(ctx, nil, "创建用户失败")
			} else {
				response.Success(ctx, nil, "创建用户成功")
			}
		}

		return nil
	})
}

// UpdateUser 更新用户信息
func UpdateUser(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_USERMANAGER) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		AlterID			uint 	`json:"alter_id"`
		UserName 		string 	`json:"user_name"`
		Password		string	`json:"password"`
		Email			string 	`json:"email"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	result := common.DB.Model(&models.User{}).Where("id = ?", dto.AlterID).Updates(models.User{UserName: dto.UserName, UserPassword: dto.Password, UserEmail: dto.Email})
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "更新用户失败")
	} else {
		response.Success(ctx, nil, "更新用户成功")
	}
}

// DeleteUser 删除用户
func DeleteUser(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_USERMANAGER) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		DeleteID 	uint	`json:"delete_id"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	result := common.DB.Where("id = ?", dto.DeleteID).Delete(&models.User{})
	if result.Error != nil {
		response.Fail(ctx, nil, "删除失败")
		return
	}
	result = common.DB.Where("user_id = ?", dto.DeleteID).Delete(&models.UserIntegral{})
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "删除失败")
	} else {
		response.Success(ctx, nil, "删除成功")
	}
}

// GetsUsersList 获取所有用户列表
func GetsUsersList(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_USERMANAGER) {
		return
	}
	// 数据库操作
	var userList []models.User
	result := common.DB.Model(&models.User{}).Find(&userList)

	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "获取信息错误")
	} else {
		type responseDTO struct {
			Id 			uint 	`json:"id"`
			UserName 	string 	`json:"user_name"`
			UserEmail 	string 	`json:"user_email"`
			DeletedAt 	gorm.DeletedAt `json:"deleted_at"`
		}
		var resDTO []responseDTO
		for _, value := range userList {
			elem := responseDTO{Id: value.ID, UserName: value.UserName, UserEmail: value.UserEmail, DeletedAt: value.DeletedAt}
			resDTO = append(resDTO, elem)
		}
		response.Success(ctx, resDTO, "获取信息成功")
	}
}

// isUserExist 检测User表中字段值是否存在
func isUserExist(column string, value string) (bool, error) {
	var total int64
	sqlStr := column + "= ?"
	if err := common.DB.Model(&models.User{}).Where(sqlStr,value).Count(&total).Error; err != nil {
		return false, errors.New("内部错误")
	}
	if total == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

