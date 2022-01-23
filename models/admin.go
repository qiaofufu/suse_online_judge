package models

import (
	"gorm.io/gorm"
	"onlineJudge/common"
)

type Admin struct {
	gorm.Model
	AdminName string  `json:"admin_name"`
	AdminAccount string	`json:"admin_account"`
	AdminRoles string	`json:"admin_roles"`
	AdminPassword string `json:"admin_password"`
}

func (a Admin) getName(aid uint) string {
	common.DB.Model(&Admin{}).Where("id = ?", aid).First(&a)
	return a.AdminName
}