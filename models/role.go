package models

import "gorm.io/gorm"

type Roles struct {
	gorm.Model
	RolesContent string `json:"roles_content"`
}
