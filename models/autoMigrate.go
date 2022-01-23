package models

import "onlineJudge/common"

// AutoMigrate 数据迁移
func AutoMigrate() {
	common.DB.AutoMigrate(&Admin{})
	common.DB.AutoMigrate(&Announcement{})
	common.DB.AutoMigrate(&User{})
	common.DB.AutoMigrate(&UserIntegral{})
	common.DB.AutoMigrate(&Goods{})
	common.DB.AutoMigrate(&GoodsType{})
	common.DB.AutoMigrate(&Roles{})
	common.DB.AutoMigrate(&Order{})
	common.DB.AutoMigrate(&Contest{})
	common.DB.AutoMigrate(&Problem{})
	common.DB.AutoMigrate(&Example{})
	common.DB.AutoMigrate(&Slideshow{})
	common.DB.AutoMigrate(&Submission{})
	common.DB.AutoMigrate(&ProblemStatusRecord{})
	common.DB.AutoMigrate(&ContestRecord{})
	//common.DB.AutoMigrate(&ContestProblemRecord{})
	common.DB.AutoMigrate(&JudgeMessage{})

}