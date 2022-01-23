package models

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	BookCount uint
	GetTime time.Time `gorm:"default:null"`
	Status int	`gorm:"default:0"`
	Value uint
	UserID uint
	GoodsID uint
}
