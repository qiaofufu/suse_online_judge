package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"onlineJudge/common"
	"onlineJudge/models"
	"onlineJudge/response"
	"onlineJudge/util"
	"strconv"
)

// AddGoodsType 添加商品类型
func AddGoodsType(ctx *gin.Context) {
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
		goodsType := models.GoodsType{Content: dto.Content}
		if err := tx.Model(&models.GoodsType{}).Where("content = ?", dto.Content).Count(&total).Error; err != nil {
			response.Fail(ctx, nil, "内部错误")
			return err
		}
		if total == 0 {
			result := tx.Model(&models.GoodsType{}).Create(&goodsType)
			if result.Error != nil {
				response.Fail(ctx, nil, "创建失败")
			} else {
				response.Success(ctx, nil, "创建成功")
			}
		} else {
			response.Fail(ctx, nil, "商品类别以存在")
		}
		return nil
	})
}

// UpdateGoodsType 更新商品类型
func UpdateGoodsType(ctx *gin.Context) {
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
	result := common.DB.Model(&models.GoodsType{}).Where("id = ?", dto.AlterID).Update("content", dto.Content)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "更新失败")
	} else {
		response.Success(ctx, nil, "更新成功")
	}
}

// DeleteGoodsType 删除商品类型
func DeleteGoodsType(ctx *gin.Context) {
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
	result := common.DB.Where("id = ?", dto.DeleteID).Delete(&models.GoodsType{})
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "删除失败")
	} else {
		response.Success(ctx, nil, "删除成功")
	}
}

// AdminGetsGoodsTypeList 获取货物类型列表
func AdminGetsGoodsTypeList(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_GOODSMANAGER) {
		return
	}
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

// AddGoods 添加商品
func AddGoods(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_GOODSMANAGER) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		GoodsName string `json:"goods_name"`
		GoodsCount uint `json:"goods_count"`
		GoodsValue uint `json:"goods_value"`
		GoodsType string `json:"goods_type"`
		GoodsDescription string `json:"goods_description"`
	}
	var dto requiredDTO
	dto.GoodsName = ctx.PostForm("goods_name")
	v, _ := strconv.Atoi(ctx.PostForm("goods_count"))
	dto.GoodsCount = uint(v)
	v, _ = strconv.Atoi(ctx.PostForm("goods_value"))
	dto.GoodsValue = uint(v)
	dto.GoodsDescription = ctx.PostForm("goods_description")
	dto.GoodsType = ctx.PostForm("goods_type")
	photo, err := savePicture(ctx,"photo",dto.GoodsType + "-" + dto.GoodsName, "goods")
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}
	// 数据库操作
	goods := models.Goods{GoodsName: dto.GoodsName, GoodsCount: dto.GoodsCount, GoodsType: dto.GoodsType, GoodsValue: dto.GoodsValue, GoodsDescription: dto.GoodsDescription, GoodsPhotoUrl: photo}
	result := common.DB.Model(&models.Goods{}).Create(&goods)
	if result.Error != nil {
		response.Fail(ctx, nil, "添加货物失败")
	} else {
		response.Success(ctx, nil, "添加成功")
	}
}

// UpdateGoods 更新商品
func UpdateGoods(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_GOODSMANAGER) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		AlterID uint `json:"alter_id"`
		GoodsName string `json:"goods_name"`
		GoodsDescription string `json:"goods_description"`
		GoodsValue uint `json:"goods_value"`
		GoodsCount uint `json:"goods_count"`
		GoodsType string `json:"goods_type"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	result := common.DB.Model(&models.Goods{}).Where("id = ?", dto.AlterID).Updates(models.Goods{GoodsName: dto.GoodsName, GoodsDescription: dto.GoodsDescription, GoodsType: dto.GoodsType, GoodsValue: dto.GoodsValue, GoodsCount: dto.GoodsCount})
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "更新失败")
	} else {
		response.Success(ctx, nil, "更新成功")
	}
}

// UpdateGoodsPicture 更新商品图片
func UpdateGoodsPicture(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_GOODSMANAGER) {
		return
	}
	// 数据绑定
	goodsID := ctx.PostForm("alter_id")
	var goods models.Goods
	result := common.DB.Model(&models.Goods{}).Where("id = ?", goodsID).First(&goods)
	if result.Error != nil {
		response.Fail(ctx, nil, "商品不存在")
		return
	}
	pictureUrl, err := savePicture(ctx, "photo", goods.GoodsType + "-" + goods.GoodsName, "goods")
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}
	goods.GoodsPhotoUrl = pictureUrl
	result = common.DB.Save(&goods)
	if result.Error != nil {
		response.Fail(ctx, nil, "更新失败")
	} else {
		response.Success(ctx, nil, "更新成功")
	}
}

// DeleteGoods 删除商品
func DeleteGoods(ctx *gin.Context) {
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
	result := common.DB.Where("id = ?", dto.DeleteID).Delete(&models.Goods{})
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "删除失败")
	} else {
		response.Success(ctx, nil, "删除成功")
	}
}

// AdminGetsGoodsList 获取商品列表
func AdminGetsGoodsList(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_GOODSMANAGER) {
		return
	}
	// 数据库操作
	var goods []models.Goods
	result := common.DB.Model(&models.Goods{}).Find(&goods)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "获取信息错误")
	} else {
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

// GetsOrderList 获取订单列表
func GetsOrderList(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_GOODSMANAGER) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		UserID uint `json:"user_id"`
		GoodsID uint `json:"goods_id"`
		Status int `json:"status"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	var orderList []models.Order
	result := common.DB.Model(&models.Order{})
	if dto.UserID != 0 {
		result = result.Where("user_id = ?",dto.UserID)
	}
	if dto.GoodsID != 0 {
		result = result.Where("goods_id = ?", dto.GoodsID)
	}
	if dto.Status != 0 {
		result = result.Where("status = ?", dto.Status)
	}
	result = result.Order("id desc").Find(&orderList)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "查询失败")
	} else {
		response.Success(ctx, orderList, "查询成功")
	}
}

// FinishOrder 完成订单
func FinishOrder(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_GOODSMANAGER) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		OrderID uint `json:"order_id"`
		GetTime string `json:"get_time"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	getTime, err := util.GetTime(dto.GetTime)
	if err != nil {
		response.Fail(ctx, nil, "时间转换失败")
	}
	// 数据库操作
	result := common.DB.Model(&models.Order{}).Where("id = ?", dto.OrderID).Updates(models.Order{GetTime: getTime, Status: 1})
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "更新订单状态失败")
	} else {
		response.Success(ctx, nil, "完成订单成功")
	}
}
