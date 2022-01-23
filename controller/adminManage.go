package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"onlineJudge/common"
	"onlineJudge/models"
	"onlineJudge/response"
)

/*
*	管理员管理功能
*/

// AddAdmin 添加管理员
func AddAdmin(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_ADMINMANAGER) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		AdminName 		string 	`json:"admin_name"`
		AdminAccount 	string 	`json:"admin_account"`
		Roles			string 	`json:"roles"`
		Password		string	`json:"password"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作 并 响应
	common.DB.Transaction(func(tx *gorm.DB) error {
		admin := models.Admin{AdminName: dto.AdminName, AdminAccount: dto.AdminAccount, AdminRoles: dto.Roles, AdminPassword: dto.Password}
		var total int64
		if err := tx.Model(&models.Admin{}).Where("admin_account = ?", dto.AdminAccount).Count(&total).Error; err != nil {
			response.Fail(ctx, nil, "内部错误")
		}
		if total == 0 {
			if err := tx.Model(&models.Admin{}).Create(&admin).Error; err != nil {
				response.Fail(ctx, nil, "添加失败")
				return err
			} else {
				response.Success(ctx, nil, "添加成功")
			}
		} else {
			response.Fail(ctx, nil, "账号已存在")
		}
		return nil
	})
}

// UpdateAdmin 更新管理员角色
func UpdateAdmin(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_ADMINMANAGER) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		AlterID			uint 	`json:"alter_id"`
		AdminName 		string 	`json:"admin_name"`
		AdminAccount 	string 	`json:"admin_account"`
		Roles			string 	`json:"roles"`
		Password		string	`json:"password"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	result := common.DB.Model(&models.Admin{}).Where("id = ?", dto.AlterID).Updates(models.Admin{AdminName: dto.AdminName, AdminAccount: dto.AdminAccount, AdminPassword: dto.Password, AdminRoles: dto.Roles})
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "更新管理员角色失败")
	} else {
		response.Success(ctx, nil, "更新管理员角色成功")
	}
}

// DeleteAdmin 删除管理员
func DeleteAdmin(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_ADMINMANAGER) {
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
	result := common.DB.Where("id = ?", dto.DeleteID).Delete(&models.Admin{})
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "删除失败")
	} else {
		response.Success(ctx, nil, "删除成功")
	}
}

// GetsAdminsList 获取所有管理员列表
func GetsAdminsList(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_ADMINMANAGER) {
		return
	}
	// 数据库操作
	var adminList []models.Admin
	result := common.DB.Model(&models.Admin{}).Find(&adminList)

	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "获取信息错误")
	} else {
		type responseDTO struct {
			ID uint `json:"id"`
			AdminName string `json:"admin_name"`
			AdminAccount string `json:"admin_account"`
			AdminRoles	string `json:"admin_roles"`
			AdminPassword string `json:"admin_password"`
		}
		var resDTO []responseDTO
		for _, value := range adminList {
			resDTO = append(resDTO, responseDTO{AdminName: value.AdminName, ID: value.ID, AdminRoles: value.AdminRoles, AdminAccount: value.AdminAccount, AdminPassword: value.AdminPassword})
		}
		response.Success(ctx, resDTO, "获取信息成功")
	}
}

// AddRole 添加角色
func AddRole(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_GOODSMANAGER) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		Content string `json:"content"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	common.DB.Transaction(func(tx *gorm.DB) error {
		var total int64
		role := models.Roles{RolesContent: dto.Content}
		if err := tx.Model(&models.Roles{}).Where("roles_content = ?", dto.Content).Count(&total).Error; err != nil {
			response.Fail(ctx, nil, "内部错误")
			return err
		}
		if total == 0 {
			result := tx.Model(&models.Roles{}).Create(&role)
			if result.Error != nil {
				response.Fail(ctx, nil, "创建失败")
			} else {
				response.Success(ctx, nil, "创建成功")
			}
		} else {
			response.Fail(ctx, nil, "角色以存在")
		}
		return nil
	})
}

// UpdateRole 更新角色
func UpdateRole(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_GOODSMANAGER) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		AlterID		uint		`json:"alter_id"`
		Content		string		`json:"content"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	result := common.DB.Model(&models.Roles{}).Where("id = ?", dto.AlterID).Update("roles_content", dto.Content)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "更新失败")
	} else {
		response.Success(ctx, nil, "更新成功")
	}
}

// DeleteRole 删除角色
func DeleteRole(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_GOODSMANAGER) {
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
	result := common.DB.Where("id = ?", dto.DeleteID).Delete(&models.Roles{})
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "删除失败")
	} else {
		response.Success(ctx, nil, "删除成功")
	}
}

// GetsRolesList 获取角色列表
func GetsRolesList(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_GOODSMANAGER) {
		return
	}
	// 数据库操作
	var role []models.Roles
	result := common.DB.Model(&models.Roles{}).Find(&role)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "获取信息错误")
	} else {
		type responseDTO struct {
			ID		uint 	`json:"id"`
			Content	string 	`json:"content"`
		}
		var resDTO []responseDTO
		for _, value := range role {
			elem := responseDTO{ID: value.ID, Content: value.RolesContent}
			resDTO = append(resDTO, elem)
		}
		response.Success(ctx, resDTO, "获取信息成功")
	}
}




