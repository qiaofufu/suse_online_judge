package models

import (
	"gorm.io/gorm"
	"onlineJudge/common"
)

type UserIntegral struct {
	gorm.Model
	Value uint	`gorm:"default:0"`
	ConsumptionValue uint `gorm:"default:0"`
	UserID uint
}

func (u UserIntegral) AddIntegral(userID uint, value uint)  {
	common.DB.Model(&UserIntegral{}).Where("user_id = ?", userID).First(&u)
	u.Value += value
	common.DB.Model(&UserIntegral{}).Save(&u)
}