package models

import "gorm.io/gorm"

type Goods struct {
	gorm.Model
	GoodsName string
	GoodsDescription string
	GoodsPhotoUrl string
	GoodsValue uint
	GoodsCount uint
	GoodsType string
}

type GoodsType struct {
	gorm.Model
	Content string
}
