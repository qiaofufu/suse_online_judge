package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"onlineJudge/common"
	"onlineJudge/models"
	"onlineJudge/response"
)

// GetsGoodsTypeList 获取货物类型列表
func GetsGoodsTypeList(ctx *gin.Context) {
	// 数据库操作
	var goodType []models.GoodsType
	result := common.DB.Model(&models.GoodsType{}).Find(&goodType)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "获取信息错误")
	} else {
		type responseDTO struct {
			ID		uint 	`json:"id"`
			Content	string 	`json:"content"`
		}
		var resDTO []responseDTO
		for _, value := range goodType {
			elem := responseDTO{ID: value.ID, Content: value.Content}
			resDTO = append(resDTO, elem)
		}
		response.Success(ctx, resDTO, "获取信息成功")
	}
}

// GetsGoodsList 获取商品列表
func GetsGoodsList(ctx *gin.Context) {
	// 数据绑定
	type requiredDTO struct {
		GoodsType string `json:"goods_type"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	var goods []models.Goods
	result := common.DB.Model(&models.Goods{}).Where("goods_type = ?", dto.GoodsType).Find(&goods)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "获取信息错误")
	} else {
		fmt.Println()
		type responseDTO struct {
			ID		uint 	`json:"id"`
			GoodsName string `json:"goods_name"`
			GoodsValue uint `json:"goods_value"`
			GoodsCount uint `json:"goods_count"`
			GoodsType string `json:"goods_type"`
			GoodsPhotoUrl string `json:"goods_photo_url"`
			GoodsDescription string `json:"goods_description"`
		}
		var resDTO []responseDTO
		for _, value := range goods {
			elem := responseDTO{ID: value.ID, GoodsName: value.GoodsName, GoodsDescription: value.GoodsDescription, GoodsType: value.GoodsType, GoodsValue: value.GoodsValue, GoodsCount: value.GoodsCount, GoodsPhotoUrl: value.GoodsPhotoUrl}
			resDTO = append(resDTO, elem)
		}
		response.Success(ctx, resDTO, "获取信息成功")
	}
}

// BuyGoods 购买商品
func BuyGoods(ctx *gin.Context) {
	// 数据绑定
	type requiredDTO struct {
		GoodsID uint `json:"goods_id"`
		Number uint `json:"number"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	user, err := getUserByJWT(ctx)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}
	// 处理
	common.DB.Transaction(func(tx *gorm.DB) error {
		var userIntegral models.UserIntegral
		if err := tx.Model(&models.UserIntegral{}).Where("user_id = ?", user.ID).First(&userIntegral).Error; err != nil {
			response.Fail(ctx, nil, "购买失败1")
			return err
		}
		var goods models.Goods
		if err := tx.Model(&models.Goods{}).Where("id = ?", dto.GoodsID).First(&goods).Error; err != nil {
			response.Fail(ctx, nil, "购买失败2")
			return err
		}
		if goods.GoodsCount < dto.Number {
			response.Fail(ctx, nil, "库存不足")
			return errors.New("库存不足")
		}
		if userIntegral.Value - userIntegral.ConsumptionValue < goods.GoodsValue * dto.Number {
			response.Fail(ctx, nil, "积分不足，无法购买")
			return errors.New("积分不足")
		}
		userIntegral.ConsumptionValue += goods.GoodsValue * dto.Number
		if err := tx.Save(&userIntegral).Error; err != nil {
			response.Fail(ctx, nil, "购买失败3")
			return err
		}
		goods.GoodsCount -= dto.Number
		if err := tx.Save(&goods).Error; err != nil {
			response.Fail(ctx, nil, "购买失败4")
			return err
		}
		order := models.Order{GoodsID: dto.GoodsID, UserID: user.ID, Value: goods.GoodsValue * dto.Number, BookCount: dto.Number}
		if err := tx.Create(&order).Error; err != nil {
			response.Fail(ctx, nil, "购买失败5")
			return err
		}
		response.Success(ctx, nil, "购买成功，请前往算法竞赛部领取奖品")
		return nil
	})
}

// GetsUserOrderList 获取用户订单列表
func GetsUserOrderList(ctx *gin.Context) {
	type requiredDTO struct {
		Status int `json:"status"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	user, err := getUserByJWT(ctx)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}

	var orderList []models.Order
	result := common.DB.Model(&models.Order{}).Where("user_id = ? and status = ?", user.ID, dto.Status).Find(&orderList)
	if result.Error != nil {
		response.Fail(ctx, nil, "查找失败")
	}
	response.Success(ctx, orderList, "查找成功")
}